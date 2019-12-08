package scripting

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/seveas/herd"
	"github.com/seveas/herd/scripting/parser"
	"github.com/sirupsen/logrus"
)

var typeNames = map[int]string{
	parser.HerdParserNUMBER:   "number",
	parser.HerdParserSTRING:   "string",
	parser.HerdParserDURATION: "duration",
}

var tokenConverters = map[int]func(string) (interface{}, error){
	parser.HerdParserNUMBER:   func(s string) (interface{}, error) { return strconv.ParseInt(s, 0, 64) },
	parser.HerdParserSTRING:   func(s string) (interface{}, error) { return strconv.Unquote(s) },
	parser.HerdParserDURATION: func(s string) (interface{}, error) { return time.ParseDuration(s) },
	parser.HerdParserREGEXP: func(s string) (interface{}, error) {
		return regexp.Compile(strings.Replace(s[1:len(s)-1], "\\/", "/", -1))
	},
	parser.HerdParserIDENTIFIER: func(s string) (interface{}, error) {
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

var variables map[string]variable = map[string]variable{
	"Timeout": {
		tokenType: parser.HerdParserDURATION,
	},
	"HostTimeout": {
		tokenType: parser.HerdParserDURATION,
	},
	"ConnectTimeout": {
		tokenType: parser.HerdParserDURATION,
	},
	"Parallel": {
		tokenType: parser.HerdParserDURATION,
	},
	"Output": {
		tokenType: parser.HerdParserSTRING,
		validator: func(i interface{}) (interface{}, error) {
			s := i.(string)
			outputModes := map[string]herd.OutputMode{
				"all":      herd.OutputAll,
				"inline":   herd.OutputInline,
				"per-host": herd.OutputPerhost,
				"tail":     herd.OutputTail,
			}
			if s == "all" || s == "per-host" || s == "inline" || s == "tail" {
				return outputModes[s], nil
			}
			return nil, fmt.Errorf("Unknown output mode: %s. Known modes: all, per-host, inline, tail", s)
		},
	},
	"NoPager": {
		tokenType: parser.HerdParserIDENTIFIER,
		validator: func(i interface{}) (interface{}, error) {
			if _, ok := i.(bool); ok {
				return i, nil
			}
			return nil, fmt.Errorf("Expected a boolean value, not %v", i)
		},
	},
	"LogLevel": {
		tokenType: parser.HerdParserSTRING,
		validator: func(i interface{}) (interface{}, error) {
			s := i.(string)
			if level, err := logrus.ParseLevel(s); err == nil {
				return level, nil
			}
			return nil, fmt.Errorf("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", s)
		},
	},
}

type herdListener struct {
	*parser.BaseHerdListener
	commands      []command
	errorListener *herdErrorListener
}

func (l *herdListener) ExitSet(c *parser.SetContext) {
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

func (l *herdListener) ExitAdd(c *parser.AddContext) {
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

func (l *herdListener) ExitRemove(c *parser.RemoveContext) {
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

func (l *herdListener) parseFilters(filters []parser.IFilterContext) herd.MatchAttributes {
	attrs := make(herd.MatchAttributes, 0)
	for _, filter := range filters {
		// If there already are lexer/parser errors, don't bother anymore
		for _, child := range filter.GetChildren() {
			if _, ok := child.(*antlr.ErrorNodeImpl); ok {
				return attrs
			}
		}
		key := filter.GetChild(0).(antlr.ParseTree)
		attr := herd.MatchAttribute{Name: key.GetText()}
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

func (l *herdListener) ExitList(c *parser.ListContext) {
	if l.errorListener.hasErrors() {
		return
	}
	oneline := c.GetOneline()
	command := listHostsCommand{oneLine: oneline != nil}
	l.commands = append(l.commands, command)
}

func (l *herdListener) ExitRun(c *parser.RunContext) {
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

func ParseCode(code string) ([]command, error) {
	is := antlr.NewInputStream(code)
	lexer := parser.NewHerdLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewHerdParser(stream)
	el := herdErrorListener{
		errors: &herd.MultiError{Subject: "Syntax errors found"},
	}
	l := herdListener{
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

type herdErrorListener struct {
	*antlr.DefaultErrorListener
	errors *herd.MultiError
}

func (l *herdErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	msg = strings.ReplaceAll(msg, "'\n'", "<NEWLINE>")
	l.errors.Add(fmt.Errorf("line %d:%d %s", line, column, msg))
}

func (l *herdErrorListener) hasErrors() bool {
	return l.errors.HasErrors()
}
