package scripting

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/katyusha"
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
		e.Ui.SetOutputMode(c.value.(katyusha.OutputMode))
	case "Timestamp":
		e.Ui.SetOutputTimestamp(c.value.(bool))
	case "NoPager":
		e.Ui.SetPagerEnabled(!c.value.(bool))
	case "NoColor":
		ansi.DisableColors(c.value.(bool))
	case "Timeout":
		e.Runner.SetTimeout(c.value.(time.Duration))
	case "HostTimeout":
		e.Runner.SetHostTimeout(c.value.(time.Duration))
	case "ConnectTimeout":
		e.Runner.SetConnectTimeout(c.value.(time.Duration))
	case "Parallel":
		e.Runner.SetParallel(c.value.(int))
	}
	return nil
}

func (c setCommand) String() string {
	return fmt.Sprintf("set %s %v", c.variable, c.value)
}

type addHostsCommand struct {
	glob       string
	attributes katyusha.MatchAttributes
}

func (c addHostsCommand) execute(e *ScriptEngine) error {
	e.Runner.AddHosts(c.glob, c.attributes)
	return nil
}

func (c addHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.glob, c.attributes)
}

type removeHostsCommand struct {
	glob       string
	attributes katyusha.MatchAttributes
}

func (c removeHostsCommand) execute(e *ScriptEngine) error {
	e.Runner.RemoveHosts(c.glob, c.attributes)
	return nil
}

func (c removeHostsCommand) String() string {
	return fmt.Sprintf("remove hosts %s %v", c.glob, c.attributes)
}

type listHostsCommand struct {
	opts katyusha.HostListOptions
}

func (c listHostsCommand) execute(e *ScriptEngine) error {
	e.Ui.PrintHostList(e.Runner.GetHosts(), c.opts)
	return nil
}

func (c listHostsCommand) String() string {
	ret := "list hosts"
	if c.opts.OneLine {
		ret += " --oneline"
	}
	if c.opts.AllAttributes {
		ret += " --all-attributes"
	} else if len(c.opts.Attributes) != 0 {
		ret += " --attributes=" + strings.Join(c.opts.Attributes, ",")
	}
	return ret
}

type runCommand struct {
	command string
}

func (c runCommand) execute(e *ScriptEngine) error {
	oc := e.Ui.OutputChannel(e.Runner)
	pc := e.Ui.ProgressChannel(e.Runner)
	hi := e.Runner.Run(c.command, pc, oc)
	e.Ui.Sync()
	e.History = append(e.History, hi)
	if !strings.HasPrefix(c.command, "katyusha:") {
		e.Ui.PrintHistoryItem(hi)
	}
	return nil
}

func (c runCommand) String() string {
	return "run " + c.command
}
