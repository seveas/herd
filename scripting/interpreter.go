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

func convertValue(c parser.IValueContext) (interface{}, error) {
	vc := c.(*parser.ValueContext)
	if a := vc.Array(); a != nil {
		return convertArray(a)
	}
	if h := vc.Hash(); h != nil {
		return convertHash(h)
	}
	if s := vc.Scalar(); s != nil {
		return convertScalar(s)
	}
	return nil, fmt.Errorf("I don't know what to do with this value: %s", c.GetText())
}

func convertScalar(c parser.IScalarContext) (interface{}, error) {
	sc := c.(*parser.ScalarContext)
	if n := sc.NUMBER(); n != nil {
		return strconv.ParseInt(n.GetText(), 0, 64)
	}
	if s := sc.STRING(); s != nil {
		return strconv.Unquote(s.GetText())
	}
	if d := sc.DURATION(); d != nil {
		return time.ParseDuration(d.GetText())
	}
	if i := sc.IDENTIFIER(); i != nil {
		switch i.GetText() {
		case "nil":
			return nil, nil
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("Unknown variable: %s", sc.GetText())
		}
	}
	return nil, fmt.Errorf("I don't know what to do with this scalar: %s", c.GetText())
}

func convertArray(c parser.IArrayContext) ([]interface{}, error) {
	values := c.(*parser.ArrayContext).AllValue()
	ret := make([]interface{}, len(values))
	for i, v := range values {
		gv, err := convertValue(v)
		if err != nil {
			return nil, err
		}
		ret[i] = gv
	}
	return ret, nil
}

func convertHash(c parser.IHashContext) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	hc := c.(*parser.HashContext)
	identifiers := hc.AllIDENTIFIER()
	values := hc.AllValue()
	for i, v := range values {
		gv, err := convertValue(v)
		if err != nil {
			return nil, err
		}
		ret[identifiers[i].GetText()] = gv
	}
	return ret, nil
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
	varValue, err := convertScalar(c.GetVarvalue())

	if err != nil {
		c.GetParser().NotifyErrorListeners(err.Error(), c.GetVarvalue().GetStart(), nil)
		return
	}

	switch varName {
	case "Splay":
		fallthrough
	case "Timeout":
		fallthrough
	case "HostTimeout":
		fallthrough
	case "ConnectTimeout":
		if _, ok := varValue.(time.Duration); !ok {
			err = fmt.Errorf("%s must be a duration", varName)
		}

	case "Parallel":
		if _, ok := varValue.(int64); !ok {
			err = fmt.Errorf("%s must be a number", varName)
		}
	case "Timestamp":
		fallthrough
	case "NoPager":
		fallthrough
	case "NoColor":
		if _, ok := varValue.(bool); !ok {
			err = fmt.Errorf("%s must be a boolean", varName)
		}
	case "Output":
		if s, ok := varValue.(string); ok {
			outputModes := map[string]herd.OutputMode{
				"all":      herd.OutputAll,
				"inline":   herd.OutputInline,
				"per-host": herd.OutputPerhost,
				"tail":     herd.OutputTail,
			}
			if s == "all" || s == "per-host" || s == "inline" || s == "tail" {
				varValue = outputModes[s]
			} else {
				err = fmt.Errorf("Unknown output mode: %s. Known modes: all, per-host, inline, tail", s)
			}
		} else {
			err = fmt.Errorf("%s must be a string", varName)
		}
	case "LogLevel":
		if s, ok := varValue.(string); ok {
			if level, perr := logrus.ParseLevel(s); perr == nil {
				varValue = level
			} else {
				fmt.Println("unknown")
				err = fmt.Errorf("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", s)
			}
		} else {
			err = fmt.Errorf("%s must be a string", varName)
		}
	}

	if err != nil {
		c.GetParser().NotifyErrorListeners(err.Error(), c.GetVarvalue().GetStart(), nil)
		return
	}

	command := setCommand{
		variable: varName,
		value:    varValue,
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

		key := filter.GetKey().GetText()
		attr := herd.MatchAttribute{Name: key}
		comp := filter.GetComp().GetText()
		if strings.HasPrefix(comp, "!") {
			attr.Negate = true
		}
		if strings.HasSuffix(comp, "~") {
			s := filter.GetRx().GetText()
			value, err := regexp.Compile(strings.Replace(s[1:len(s)-1], "\\/", "/", -1))
			if err != nil {
				filter.GetParser().NotifyErrorListeners(err.Error(), filter.GetVal().GetStart(), nil)
				continue
			}
			attr.Regex = true
			attr.FuzzyTyping = false
			attr.Value = value
		} else {
			value, err := convertScalar(filter.GetVal())
			if err != nil {
				filter.GetParser().NotifyErrorListeners(err.Error(), filter.GetVal().GetStart(), nil)
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
	opts := herd.HostListOptions{
		Separator: ",",
		Align:     true,
		Header:    true,
	}
	optval := c.GetOpts()
	if optval != nil {
		optmap, err := convertHash(optval)
		if err != nil {
			c.GetParser().NotifyErrorListeners(err.Error(), c.GetStart(), nil)
			return
		}
		copyBool := func(name string, vp *bool) {
			if v, ok := optmap[name]; ok {
				if b, ok := v.(bool); ok {
					*vp = b
				} else {
					c.GetParser().NotifyErrorListeners(fmt.Sprintf("%s must be a boolean value", name), c.GetStart(), nil)
				}
			}
		}
		copyBool("OneLine", &opts.OneLine)
		copyBool("Csv", &opts.Csv)
		copyBool("Align", &opts.Align)
		copyBool("AllAttributes", &opts.AllAttributes)
		copyBool("Header", &opts.Header)

		if v, ok := optmap["Separator"]; ok {
			if s, ok := v.(string); ok {
				opts.Separator = s
			} else {
				c.GetParser().NotifyErrorListeners("Separator must be a string", c.GetStart(), nil)
			}
		}
		if v, ok := optmap["Attributes"]; ok {
			if ss, ok := v.([]interface{}); ok {
				opts.Attributes = make([]string, 0)
				for _, e := range ss {
					if s, ok := e.(string); ok {
						opts.Attributes = append(opts.Attributes, s)
					} else {
						c.GetParser().NotifyErrorListeners("Attributes must be a list of strings", c.GetStart(), nil)
					}
				}
			} else {
				c.GetParser().NotifyErrorListeners("Attributes must be a list of strings", c.GetStart(), nil)
			}
		}
	}
	command := listHostsCommand{opts: opts}
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

func parseCode(code string) ([]command, error) {
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
