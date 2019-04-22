package herd

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/seveas/herd/parser"
)

type herdListener struct {
	*parser.BaseHerdListener

	Commands []Command
	Config   *AppConfig
}

func (l *herdListener) ExitSet(c *parser.SetContext) {
	varName := c.GetVarname().GetText()
	valueCtx := c.GetVarvalue()
	valueToken := valueCtx.GetStart()

	command := SetCommand{VariableName: varName}

	// FIXME: proper dispatch?
	if varName == "timeout" {
		command.Variable = &l.Config.Runner.Timeout
		if valueToken.GetTokenType() != parser.HerdParserNUMBER {
			p := valueCtx.GetParser()
			err := fmt.Errorf("timeout value should be a number")
			p.NotifyErrorListeners(err.Error(), valueToken, nil)
			return
		}
		seconds, _ := strconv.Atoi(valueToken.GetText())
		command.Value = time.Duration(seconds) * time.Second
	} else {
		err := fmt.Errorf("unknown variable")
		c.GetParser().NotifyErrorListeners(err.Error(), c.GetVarname(), nil)
		return
	}

	l.Commands = append(l.Commands, command)
}

func (l *herdListener) ExitAdd(c *parser.AddContext) {
	glob := c.GetGlob().GetText()
	attrs := HostAttributes{}
	filters := c.AllFilter()
	for _, filter := range filters {
		key := filter.GetChild(0).(antlr.ParseTree)
		valueCtx := filter.GetChild(2).(*parser.ValueContext)
		valueToken := valueCtx.GetStart()
		if valueToken.GetTokenType() == parser.HerdParserNUMBER {
			attrs[key.GetText()], _ = strconv.Atoi(valueToken.GetText())
		} else {
			val, err := strconv.Unquote(valueToken.GetText())
			if err != nil {
				valueCtx.GetParser().NotifyErrorListeners(err.Error(), valueToken, nil)
				continue
			}
			attrs[key.GetText()] = val
		}
	}
	command := AddHostsCommand{Glob: glob, Attributes: attrs}
	l.Commands = append(l.Commands, command)
}

func (l *herdListener) ExitRun(c *parser.RunContext) {
	command := strings.TrimLeft(c.GetText()[3:], " \t")
	if len(command) == 0 {
		err := fmt.Errorf("no command specified")
		c.GetParser().NotifyErrorListeners(err.Error(), c.GetStart(), nil)
		return
	}
	l.Commands = append(l.Commands, RunCommand{Command: command, Formatter: l.Config.Formatter})
}

func ParseScript(fn string, c *AppConfig) ([]Command, error) {
	code, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	is := antlr.NewInputStream(string(code))
	lexer := parser.NewHerdLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewHerdParser(stream)
	l := herdListener{
		Commands: make([]Command, 0),
		Config:   c,
	}
	el := HerdErrorListener{HasErrors: false}
	p.RemoveErrorListeners()
	p.AddErrorListener(&el)
	antlr.ParseTreeWalkerDefault.Walk(&l, p.Prog())
	if el.HasErrors {
		return nil, fmt.Errorf("syntax errors found")
	}
	return l.Commands, nil
}

type HerdErrorListener struct {
	*antlr.DefaultErrorListener
	HasErrors bool
}

func (l *HerdErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	UI.Errorf("line %d:%d %s", line, column, msg)
	l.HasErrors = true
}
