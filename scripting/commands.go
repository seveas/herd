package scripting

import (
	"fmt"
	"strings"
	"time"

	"github.com/seveas/herd"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
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

type showVariablesCommand struct{}

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
	hosts := e.Registry.Search(c.glob, c.attributes, c.sampled, c.count)
	e.Hosts.AddHosts(hosts)
}

func (c addHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.glob, c.attributes)
}

type removeHostsCommand struct {
	glob       string
	attributes herd.MatchAttributes
}

func (c removeHostsCommand) execute(e *ScriptEngine) {
	e.Hosts.Remove(c.glob, c.attributes)
}

func (c removeHostsCommand) String() string {
	return fmt.Sprintf("remove hosts %s %v", c.glob, c.attributes)
}

type listHostsCommand struct {
	opts herd.HostListOptions
}

func (c listHostsCommand) execute(e *ScriptEngine) {
	e.Ui.PrintHostList(c.opts)
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
	oc := e.Ui.OutputChannel()
	pc := e.Ui.ProgressChannel(time.Now().Add(e.Runner.GetTimeout()))
	hi, err := e.Runner.Run(c.command, pc, oc)
	if err != nil {
		logrus.Errorf("Unable to execute %s: %s", c.command, err)
	}
	if oc != nil {
		close(oc)
	}
	close(pc)
	e.Ui.Sync()
	if hi != nil {
		e.History = append(e.History, hi)
		e.Ui.PrintHistoryItem(hi)
	}
}

func (c runCommand) String() string {
	return "run " + c.command
}
