package scripting

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/herd"
)

type command interface {
	execute(e *ScriptEngine)
}

type setCommand struct {
	variable string
	value    interface{}
}

func (c setCommand) execute(e *ScriptEngine) {
	switch c.variable {
	case "Output":
		e.Ui.SetOutputMode(c.value.(herd.OutputMode))
	case "Timestamp":
		e.Ui.SetOutputTimestamp(c.value.(bool))
	case "NoPager":
		e.Ui.SetPagerEnabled(!c.value.(bool))
	case "NoColor":
		ansi.DisableColors(c.value.(bool))
	case "Splay":
		e.Runner.SetSplay(c.value.(time.Duration))
	case "Timeout":
		e.Runner.SetTimeout(c.value.(time.Duration))
	case "HostTimeout":
		e.Runner.SetHostTimeout(c.value.(time.Duration))
	case "ConnectTimeout":
		e.Runner.SetConnectTimeout(c.value.(time.Duration))
	case "Parallel":
		e.Runner.SetParallel(int(c.value.(int64)))
	}
}

func (c setCommand) String() string {
	return fmt.Sprintf("set %s %v", c.variable, c.value)
}

type showVariablesCommand struct {
	variable string
	value    interface{}
}

func (c showVariablesCommand) execute(e *ScriptEngine) {
	e.Ui.PrintSettings(e.Ui.Settings, e.Registry.Settings, e.Runner.Settings)
}

func (c showVariablesCommand) String() string {
	return "set"
}

type addHostsCommand struct {
	glob       string
	attributes herd.MatchAttributes
	sampled    []string
	count      int
}

func (c addHostsCommand) execute(e *ScriptEngine) {
	hosts := e.Registry.GetHosts(c.glob, c.attributes, c.sampled, c.count)
	if strings.HasPrefix(c.glob, "file:") {
		e.Runner.SetSortFields([]string{})
	}
	e.Runner.AddHosts(hosts)
}

func (c addHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.glob, c.attributes)
}

type removeHostsCommand struct {
	glob       string
	attributes herd.MatchAttributes
}

func (c removeHostsCommand) execute(e *ScriptEngine) {
	e.Runner.RemoveHosts(c.glob, c.attributes)
}

func (c removeHostsCommand) String() string {
	return fmt.Sprintf("remove hosts %s %v", c.glob, c.attributes)
}

type listHostsCommand struct {
	opts herd.HostListOptions
}

func (c listHostsCommand) execute(e *ScriptEngine) {
	e.Ui.PrintHostList(e.Runner.GetHosts(), c.opts)
}

func (c listHostsCommand) String() string {
	return fmt.Sprintf("list hosts {OneLine: %t, Separator: '%s', Csv: %t, Align: %t, Header: %t, AllAttributes: %t, Attributes: [%s]}",
		c.opts.OneLine, c.opts.Separator, c.opts.Csv, c.opts.Align, c.opts.Header,
		c.opts.AllAttributes, strings.Join(c.opts.Attributes, ", "))
}

type runCommand struct {
	command string
}

func (c runCommand) execute(e *ScriptEngine) {
	oc := e.Ui.OutputChannel(e.Runner)
	pc := e.Ui.ProgressChannel(e.Runner)
	hi := e.Runner.Run(c.command, pc, oc)
	if oc != nil {
		close(oc)
	}
	close(pc)
	e.Ui.Sync()
	if hi == nil {
		return
	}
	e.History = append(e.History, hi)
	if !strings.HasPrefix(c.command, "herd:") {
		e.Ui.PrintHistoryItem(hi)
	}
}

func (c runCommand) String() string {
	return "run " + c.command
}
