package scripting

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/scripting/parser"
	"github.com/sirupsen/logrus"
)

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

type variable struct {
	tokenType int
	validator func(interface{}) (interface{}, error)
}

func mustBeBool(i interface{}) (interface{}, error) {
	if _, ok := i.(bool); ok {
		return i, nil
	}
	return nil, fmt.Errorf("Expected a boolean value, not %v", i)
}

var variables map[string]variable = map[string]variable{
	"Timeout": {
		tokenType: parser.KatyushaParserDURATION,
	},
	"HostTimeout": {
		tokenType: parser.KatyushaParserDURATION,
	},
	"ConnectTimeout": {
		tokenType: parser.KatyushaParserDURATION,
	},
	"Parallel": {
		tokenType: parser.KatyushaParserNUMBER,
	},
	"Output": {
		tokenType: parser.KatyushaParserSTRING,
		validator: func(i interface{}) (interface{}, error) {
			s := i.(string)
			outputModes := map[string]katyusha.OutputMode{
				"all":      katyusha.OutputAll,
				"inline":   katyusha.OutputInline,
				"per-host": katyusha.OutputPerhost,
				"tail":     katyusha.OutputTail,
			}
			if s == "all" || s == "per-host" || s == "inline" || s == "tail" {
				return outputModes[s], nil
			}
			return nil, fmt.Errorf("Unknown output mode: %s. Known modes: all, per-host, inline, tail", s)
		},
	},
	"Timestamp": {
		tokenType: parser.KatyushaParserIDENTIFIER,
		validator: mustBeBool,
	},
	"NoPager": {
		tokenType: parser.KatyushaParserIDENTIFIER,
		validator: mustBeBool,
	},
	"NoColor": {
		tokenType: parser.KatyushaParserIDENTIFIER,
		validator: mustBeBool,
	},
	"LogLevel": {
		tokenType: parser.KatyushaParserSTRING,
		validator: func(i interface{}) (interface{}, error) {
			s := i.(string)
			if level, err := logrus.ParseLevel(s); err == nil {
				return level, nil
			}
			return nil, fmt.Errorf("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", s)
		},
	},
}

type katyushaListener struct {
	*parser.BaseKatyushaListener
	commands      []command
	errorListener *katyushaErrorListener
}

func (l *katyushaListener) ExitSet(c *parser.SetContext) {
	if l.errorListener.hasErrors() {
		return
	}
	varName := c.GetVarname().GetText()
	variable, ok := variables[varName]
	if !ok {
		err := fmt.Sprintf("unknown setting: %s", varName)
		c.GetParser().NotifyErrorListeners(err, c.GetVarname(), nil)
		return
	}

	valueCtx := c.GetVarvalue()
	valueToken := valueCtx.GetStart()
	valueType := valueToken.GetTokenType()

	if valueType != variable.tokenType {
		p := valueCtx.GetParser()
		err := fmt.Sprintf("%s value should be a %s, not a %s", varName, typeNames[variable.tokenType], typeNames[valueType])
		p.NotifyErrorListeners(err, valueToken, nil)
		return
	}

	val, err := tokenConverters[valueType](valueToken.GetText())
	if err != nil {
		p := valueCtx.GetParser()
		p.NotifyErrorListeners(err.Error(), valueToken, nil)
		return
	}
	if variable.validator != nil {
		val, err = variable.validator(val)
		if err != nil {
			p := valueCtx.GetParser()
			p.NotifyErrorListeners(err.Error(), valueToken, nil)
			return
		}
	}

	command := setCommand{
		variable: varName,
		value:    val,
	}

	l.commands = append(l.commands, command)
}

func (l *katyushaListener) ExitAdd(c *parser.AddContext) {
	if l.errorListener.hasErrors() {
		return
	}
	glob := "*"
	if g := c.GetGlob(); g != nil {
		glob = g.GetText()
	}
	attrs := l.parseFilters(c.AllFilter())
	command := addHostsCommand{glob: glob, attributes: attrs}
	l.commands = append(l.commands, command)
}

func (l *katyushaListener) ExitRemove(c *parser.RemoveContext) {
	if l.errorListener.hasErrors() {
		return
	}
	glob := "*"
	if g := c.GetGlob(); g != nil {
		glob = g.GetText()
	}
	attrs := l.parseFilters(c.AllFilter())
	command := removeHostsCommand{glob: glob, attributes: attrs}
	l.commands = append(l.commands, command)
}

func (l *katyushaListener) parseFilters(filters []parser.IFilterContext) katyusha.MatchAttributes {
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
	if l.errorListener.hasErrors() {
		return
	}
	oneline := c.GetOneline() != nil
	command := listHostsCommand{opts: katyusha.HostListOptions{OneLine: oneline, Header: !oneline, Separator: ","}}
	l.commands = append(l.commands, command)
}

func (l *katyushaListener) ExitRun(c *parser.RunContext) {
	if l.errorListener.hasErrors() {
		return
	}
	command := strings.TrimLeft(c.GetText()[3:], " \t")
	if len(command) == 0 {
		err := fmt.Errorf("no command specified")
		c.GetParser().NotifyErrorListeners(err.Error(), c.GetStart(), nil)
		return
	}
	l.commands = append(l.commands, runCommand{command: command})
}

func parseCode(code string) ([]command, error) {
	is := antlr.NewInputStream(code)
	lexer := parser.NewKatyushaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewKatyushaParser(stream)
	el := katyushaErrorListener{
		errors: &katyusha.MultiError{Subject: "Syntax errors found"},
	}
	l := katyushaListener{
		commands:      make([]command, 0),
		errorListener: &el,
	}
	p.RemoveErrorListeners()
	p.AddErrorListener(&el)
	antlr.ParseTreeWalkerDefault.Walk(&l, p.Prog())
	if el.hasErrors() {
		return nil, el.errors
	}
	return l.commands, nil
}

type katyushaErrorListener struct {
	*antlr.DefaultErrorListener
	errors *katyusha.MultiError
}

func (l *katyushaErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	msg = strings.ReplaceAll(msg, "'\n'", "<NEWLINE>")
	l.errors.Add(fmt.Errorf("line %d:%d %s", line, column, msg))
}

func (l *katyushaErrorListener) hasErrors() bool {
	return l.errors.HasErrors()
}
