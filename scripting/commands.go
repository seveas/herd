package scripting

import (
	"fmt"
	"strings"

	"github.com/seveas/herd"
	"github.com/spf13/viper"
)

type command interface {
	execute(e *ScriptEngine) error
}

type setCommand struct {
	variable string
	value    interface{}
}

func (c setCommand) execute(e *ScriptEngine) error {
	if c.variable == "Output" {
		e.ui.SetOutputMode(c.value.(herd.OutputMode))
	} else if c.variable == "NoPager" {
		e.ui.SetPagerEnabled(!c.value.(bool))
	} else {
		viper.Set(c.variable, c.value)
	}
	return nil
}

func (c setCommand) String() string {
	return fmt.Sprintf("set %s %v", c.variable, c.value)
}

type addHostsCommand struct {
	glob       string
	attributes herd.MatchAttributes
}

func (c addHostsCommand) execute(e *ScriptEngine) error {
	e.runner.AddHosts(c.glob, c.attributes)
	return nil
}

func (c addHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.glob, c.attributes)
}

type removeHostsCommand struct {
	glob       string
	attributes herd.MatchAttributes
}

func (c removeHostsCommand) execute(e *ScriptEngine) error {
	e.runner.RemoveHosts(c.glob, c.attributes)
	return nil
}

func (c removeHostsCommand) String() string {
	return fmt.Sprintf("remove hosts %s %v", c.glob, c.attributes)
}

type listHostsCommand struct {
	oneLine       bool
	csv           bool
	allAttributes bool
	attributes    []string
}

func (c listHostsCommand) execute(e *ScriptEngine) error {
	e.ui.PrintHostList(e.runner.Hosts, c.oneLine, c.csv, c.allAttributes, c.attributes)
	return nil
}

func (c listHostsCommand) String() string {
	ret := "list hosts"
	if c.oneLine {
		ret += " --oneline"
	}
	if c.allAttributes {
		ret += " --all-attributes"
	}
	if len(c.attributes) != 0 {
		ret += " --attributes=" + strings.Join(c.attributes, ",")
	}
	return ret
}

type runCommand struct {
	command string
}

func (c runCommand) execute(e *ScriptEngine) error {
	oc := e.ui.OutputChannel(e.runner)
	pc := e.ui.ProgressChannel(e.runner)
	hi := e.runner.Run(c.command, pc, oc)
	e.ui.PrintHistoryItem(hi)
	return nil
}

func (c runCommand) String() string {
	return "run " + c.command
}
