package scripting

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/herd"
)

type command interface {
	execute(e *ScriptEngine) error
}

type setCommand struct {
	variable string
	value    interface{}
}

func (c setCommand) execute(e *ScriptEngine) error {
	switch c.variable {
	case "Output":
		e.ui.SetOutputMode(c.value.(herd.OutputMode))
	case "Timestamp":
		e.ui.SetOutputTimestamp(c.value.(bool))
	case "NoPager":
		e.ui.SetPagerEnabled(!c.value.(bool))
	case "NoColor":
		ansi.DisableColors(c.value.(bool))
	case "Timeout":
		e.runner.SetTimeout(c.value.(time.Duration))
	case "HostTimeout":
		e.runner.SetHostTimeout(c.value.(time.Duration))
	case "ConnectTimeout":
		e.runner.SetConnectTimeout(c.value.(time.Duration))
	case "Parallel":
		e.runner.SetParallel(c.value.(int))
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
	e.ui.PrintHostList(e.runner.GetHosts(), c.oneLine, c.csv, c.allAttributes, c.attributes)
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
	e.history = append(e.history, hi)
	e.ui.PrintHistoryItem(hi)
	return nil
}

func (c runCommand) String() string {
	return "run " + c.command
}
