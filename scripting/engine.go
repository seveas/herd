package scripting

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/seveas/katyusha"
)

type ScriptEngine struct {
	commands []command
	ui       katyusha.UI
	runner   *katyusha.Runner
	position int
	history  katyusha.History
}

func NewScriptEngine(ui katyusha.UI, runner *katyusha.Runner) *ScriptEngine {
	return &ScriptEngine{
		commands: []command{},
		ui:       ui,
		runner:   runner,
		position: 0,
		history:  make(katyusha.History, 0),
	}
}

func (e *ScriptEngine) ActiveHosts() katyusha.Hosts {
	return e.runner.GetHosts()
}

func (e *ScriptEngine) ParseCommandLine(args []string, splitAt int) error {
	filters := args
	if splitAt != -1 {
		filters = filters[:splitAt]
	}
	comparison := regexp.MustCompile("^(.*?)(=~|==?|!=|!~)(.*)$")
	// First we add hosts from the command line, in all modes
	commands := make([]command, 0)
	add := true
hostspecLoop:
	for len(filters) > 0 {
		glob := filters[0]
		// Do we have a glob or not?
		if comparison.MatchString(glob) {
			glob = "*"
		} else {
			filters = filters[1:]
		}
		attrs := make(katyusha.MatchAttributes, 0)
		for i, arg := range filters[:] {
			if arg == "+" || arg == "-" {
				filters = filters[i+1:]
				if add {
					commands = append(commands, addHostsCommand{glob: glob, attributes: attrs})
				} else {
					commands = append(commands, removeHostsCommand{glob: glob, attributes: attrs})
				}
				if arg == "+" {
					add = true
				} else {
					add = false
				}
				continue hostspecLoop
			}
			parts := comparison.FindStringSubmatch(arg)
			if len(parts) == 0 {
				return fmt.Errorf("incorrect filter: %s", arg)
			}
			key, comp, val := parts[1], parts[2], parts[3]
			attr := katyusha.MatchAttribute{Name: key, Value: val, FuzzyTyping: true}
			if strings.HasPrefix(comp, "!") {
				attr.Negate = true
			}
			if strings.HasSuffix(comp, "~") {
				re, err := regexp.Compile(val)
				if err != nil {
					return fmt.Errorf("Invalid regexp /%s/: %s", val, err)
				} else {
					attr.Value = re
					attr.Regex = true
					attr.FuzzyTyping = false
				}
			}
			attrs = append(attrs, attr)
		}
		// We've fallen through, so no more hostspecs
		if add {
			commands = append(commands, addHostsCommand{glob: glob, attributes: attrs})
		} else {
			commands = append(commands, removeHostsCommand{glob: glob, attributes: attrs})
		}
		break
	}
	e.commands = append(e.commands, commands...)
	if splitAt != -1 {
		e.commands = append(e.commands, runCommand{command: strings.Join(args[splitAt:], " ")})
	}
	return nil
}

func (e *ScriptEngine) ParseScriptFile(fn string) error {
	code, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	commands, err := ParseCode(string(code))
	if err != nil {
		return err
	}
	e.commands = append(e.commands, commands...)
	return nil
}

func (e *ScriptEngine) ParseCodeLine(code string) error {
	commands, err := ParseCode(code)
	if err != nil {
		return err
	}
	e.commands = append(e.commands, commands...)
	return nil
}

func (e *ScriptEngine) AddListHostsCommand(oneLine, csv, allAttributes bool, attributes []string) {
	e.commands = append(e.commands, listHostsCommand{oneLine: oneLine, csv: csv, allAttributes: allAttributes, attributes: attributes})
}

func (e *ScriptEngine) Execute() {
	if len(e.commands) < e.position {
		return
	}
	for _, command := range e.commands[e.position:] {
		logrus.Debugf("%s", command)
		command.execute(e)
		e.position++
	}
}

func (e *ScriptEngine) SaveHistory(fn string) error {
	return e.history.Save(fn)
}

func (e *ScriptEngine) End() {
	e.runner.End()
	e.ui.Wait()
}
