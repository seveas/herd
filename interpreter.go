package katyusha

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/seveas/katyusha/parser"
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
	parser.KatyushaParserNUMBER:   func(s string) (interface{}, error) { return strconv.Atoi(s) },
	parser.KatyushaParserSTRING:   func(s string) (interface{}, error) { return strconv.Unquote(s) },
	parser.KatyushaParserDURATION: func(s string) (interface{}, error) { return time.ParseDuration(s) },
	parser.KatyushaParserREGEXP: func(s string) (interface{}, error) {
		return regexp.Compile(strings.Replace(s[1:len(s)-1], "\\/", "/", -1))
	},
}

var tokenTypes = map[string]int{
	"Timeout":        parser.KatyushaParserDURATION,
	"HostTimeout":    parser.KatyushaParserDURATION,
	"ConnectTimeout": parser.KatyushaParserDURATION,
	"Parallel":       parser.KatyushaParserNUMBER,
}

type katyushaListener struct {
	*parser.BaseKatyushaListener

	Commands []Command

	ConfigSetters map[string]ConfigSetter
}

func (l *katyushaListener) ExitSet(c *parser.SetContext) {
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

	command := SetCommand{
		Variable: varName,
		Value:    val,
	}

	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ExitAdd(c *parser.AddContext) {
	glob := c.GetGlob().GetText()
	attrs := l.ParseFilters(c.AllFilter())
	command := AddHostsCommand{Glob: glob, Attributes: attrs}
	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ExitRemove(c *parser.RemoveContext) {
	glob := c.GetGlob().GetText()
	attrs := l.ParseFilters(c.AllFilter())
	command := RemoveHostsCommand{Glob: glob, Attributes: attrs}
	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ParseFilters(filters []parser.IFilterContext) map[string]interface{} {
	attrs := HostAttributes{}
	for _, filter := range filters {
		// If there already are lexer/parser errors, don't bother anymore
		for _, child := range filter.GetChildren() {
			if _, ok := child.(*antlr.ErrorNodeImpl); ok {
				return attrs
			}
		}
		key := filter.GetChild(0).(antlr.ParseTree)
		if filter.GetChild(1).(antlr.ParseTree).GetText() == "=" {
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
			attrs[key.GetText()] = value
		} else {
			valueToken := filter.GetChild(2).(*antlr.TerminalNodeImpl).GetSymbol()
			value, err := tokenConverters[valueToken.GetTokenType()](valueToken.GetText())
			if err != nil {
				filter.GetParser().NotifyErrorListeners(err.Error(), valueToken, nil)
				continue
			}
			attrs[key.GetText()] = value
		}
	}
	return attrs
}

func (l *katyushaListener) ExitList(c *parser.ListContext) {
	oneline := c.GetOneline()
	command := ListHostsCommand{OneLine: false}
	if oneline != nil {
		command.OneLine = true
	}
	l.Commands = append(l.Commands, command)
}

func (l *katyushaListener) ExitRun(c *parser.RunContext) {
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
	l := katyushaListener{
		Commands: make([]Command, 0),
	}
	el := KatyushaErrorListener{HasErrors: false}
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
	UI.Errorf("line %d:%d %s", line, column, msg)
	l.HasErrors = true
}
