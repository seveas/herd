package scripting

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/scripting/parser"
)

type ConfigSetter struct {
	tokenType int
}

var typeNames = map[int]string{
	parser.KatyushaParserNUMBER:   "number",
	parser.KatyushaParserSTRING:   "string",
	parser.KatyushaParserDURATION: "duration",
}

var tokenConverters = map[int]func(string) (interface{}, error){
	parser.KatyushaParserNUMBER:   func(s string) (interface{}, error) { return strconv.ParseInt(s, 0, 64) },
	parser.KatyushaParserSTRING:   func(s string) (interface{}, error) { return strconv.Unquote(s) },
	parser.KatyushaParserDURATION: func(s string) (interface{}, error) { return time.ParseDuration(s) },
	parser.KatyushaParserREGEXP: func(s string) (interface{}, error) {
		return regexp.Compile(strings.Replace(s[1:len(s)-1], "\\/", "/", -1))
	},
	parser.KatyushaParserIDENTIFIER: func(s string) (interface{}, error) {
		if s == "nil" {
			return nil, nil
		}
		if s == "true" {
			return true, nil
		}
		if s == "false" {
			return false, nil
		}
		return nil, fmt.Errorf("Unknown variable: %s", s)
	},
}

var tokenTypes = map[string]int{
	"Timeout":        parser.KatyushaParserDURATION,
	"HostTimeout":    parser.KatyushaParserDURATION,
	"ConnectTimeout": parser.KatyushaParserDURATION,
	"Parallel":       parser.KatyushaParserNUMBER,
	"Output":         parser.KatyushaParserSTRING,
	"LogLevel":       parser.KatyushaParserSTRING,
}

var logLevels = map[string]int{
	"DEBUG":   katyusha.DEBUG,
	"INFO":    katyusha.INFO,
	"NORMAL":  katyusha.NORMAL,
	"WARNING": katyusha.WARNING,
	"ERROR":   katyusha.ERROR,
}
var validators = map[string]func(interface{}) (interface{}, error){
	"Output": func(i interface{}) (interface{}, error) {
		s := i.(string)
		if s == "all" || s == "host" || s == "line" || s == "pager" {
			return s, nil
		}
		return nil, fmt.Errorf("Unknown output mode: %s. Known modes: all, host, line, pager", s)
	},
	"LogLevel": func(i interface{}) (interface{}, error) {
		s := i.(string)
		if level, ok := logLevels[s]; ok {
			return level, nil
		} else {
			return nil, fmt.Errorf("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", s)
		}
	},
}

type katyushaListener struct {
	*parser.BaseKatyushaListener
	Commands      []Command
	ConfigSetters map[string]ConfigSetter
	ErrorListener *KatyushaErrorListener
}

func (l *katyushaListener) ExitSet(c *parser.SetContext) {
	if l.ErrorListener.HasErrors {
		return
	}
	varName := c.GetVarname().GetText()
	tokenType, ok := tokenTypes[varName]
	if !ok {
		err := fmt.Sprintf("unknown setting: %s", varName)
		c.GetParser().NotifyErrorListeners(err, c.GetVarname(), nil)
		return
	}

	valueCtx := c.GetVarvalue()
	valueToken := valueCtx.GetStart()
	valueType := valueToken.GetTokenType()

	if valueType != tokenType {
		p := valueCtx.GetParser()
		err := fmt.Sprintf("%s value should be a %s, not a %s", varName, typeNames[tokenType], typeNames[valueType])
		p.NotifyErrorListeners(err, valueToken, nil)
		return
	}

	val, err := tokenConverters[valueType](valueToken.GetText())
	if err != nil {
		p := valueCtx.GetParser()
		p.NotifyErrorListeners(err.Error(), valueToken, nil)
		return
	}
	validator, ok := validators[varName]
	if ok {
		val, err = validator(val)
		if err != nil {
			p := valueCtx.GetParser()
			p.NotifyErrorListeners(err.Error(), valueToken, nil)
			return
		}
	}

	command := SetCommand{
		Variable: varName,
		Value:    val,
	}

	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ExitAdd(c *parser.AddContext) {
	if l.ErrorListener.HasErrors {
		return
	}
	glob := "*"
	if g := c.GetGlob(); g != nil {
		glob = g.GetText()
	}
	attrs := l.ParseFilters(c.AllFilter())
	command := AddHostsCommand{Glob: glob, Attributes: attrs}
	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ExitRemove(c *parser.RemoveContext) {
	if l.ErrorListener.HasErrors {
		return
	}
	glob := "*"
	if g := c.GetGlob(); g != nil {
		glob = g.GetText()
	}
	attrs := l.ParseFilters(c.AllFilter())
	command := RemoveHostsCommand{Glob: glob, Attributes: attrs}
	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ParseFilters(filters []parser.IFilterContext) katyusha.MatchAttributes {
	attrs := make(katyusha.MatchAttributes, 0)
	for _, filter := range filters {
		// If there already are lexer/parser errors, don't bother anymore
		for _, child := range filter.GetChildren() {
			if _, ok := child.(*antlr.ErrorNodeImpl); ok {
				return attrs
			}
		}
		key := filter.GetChild(0).(antlr.ParseTree)
		attr := katyusha.MatchAttribute{Name: key.GetText()}
		comp := filter.GetChild(1).(antlr.ParseTree).GetText()
		if strings.HasPrefix(comp, "!") {
			attr.Negate = true
		}
		if strings.HasSuffix(comp, "~") {
			valueToken := filter.GetChild(2).(*antlr.TerminalNodeImpl).GetSymbol()
			value, err := tokenConverters[valueToken.GetTokenType()](valueToken.GetText())
			if err != nil {
				filter.GetParser().NotifyErrorListeners(err.Error(), valueToken, nil)
				continue
			}
			attr.Regex = true
			attr.FuzzyTyping = false
			attr.Value = value
		} else {
			valueCtx := filter.GetChild(2).(*parser.ValueContext)
			valueToken := valueCtx.GetStart()
			if _, ok := tokenConverters[valueToken.GetTokenType()]; !ok {
				// Unknown value, implying a syntax error
				return attrs
			}
			value, err := tokenConverters[valueToken.GetTokenType()](valueToken.GetText())
			if err != nil {
				valueCtx.GetParser().NotifyErrorListeners(err.Error(), valueToken, nil)
				continue
			}
			attr.Value = value
		}
		attrs = append(attrs, attr)
	}
	return attrs
}

func (l *katyushaListener) ExitList(c *parser.ListContext) {
	if l.ErrorListener.HasErrors {
		return
	}
	oneline := c.GetOneline()
	command := ListHostsCommand{OneLine: false}
	if oneline != nil {
		command.OneLine = true
	}
	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ExitRun(c *parser.RunContext) {
	if l.ErrorListener.HasErrors {
		return
	}
	command := strings.TrimLeft(c.GetText()[3:], " \t")
	if len(command) == 0 {
		err := fmt.Errorf("no command specified")
		c.GetParser().NotifyErrorListeners(err.Error(), c.GetStart(), nil)
		return
	}
	l.Commands = append(l.Commands, RunCommand{Command: command})
}

func ParseScript(fn string) ([]Command, error) {
	code, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return ParseCode(string(code))
}

func ParseCode(code string) ([]Command, error) {
	is := antlr.NewInputStream(code)
	lexer := parser.NewKatyushaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewKatyushaParser(stream)
	el := KatyushaErrorListener{HasErrors: false}
	l := katyushaListener{
		Commands:      make([]Command, 0),
		ErrorListener: &el,
	}
	p.RemoveErrorListeners()
	p.AddErrorListener(&el)
	antlr.ParseTreeWalkerDefault.Walk(&l, p.Prog())
	if el.HasErrors {
		return nil, fmt.Errorf("syntax errors found")
	}
	return l.Commands, nil
}

type KatyushaErrorListener struct {
	*antlr.DefaultErrorListener
	HasErrors bool
}

func (l *KatyushaErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	katyusha.UI.Errorf("line %d:%d %s", line, column, msg)
	l.HasErrors = true
}
