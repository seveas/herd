package scripting

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/seveas/katyusha"
)

type ScriptEngine struct {
	Ui       katyusha.UI
	Runner   *katyusha.Runner
	History  katyusha.History
	commands []command
	position int
}

func NewScriptEngine(ui katyusha.UI, runner *katyusha.Runner) *ScriptEngine {
	return &ScriptEngine{
		Ui:       ui,
		Runner:   runner,
		History:  make(katyusha.History, 0),
		commands: []command{},
		position: 0,
	}
}

func (e *ScriptEngine) ParseCommandLine(args []string, splitAt int) error {
	filters := args
	if splitAt != -1 {
		filters = filters[:splitAt]
	}
	comparison := regexp.MustCompile("^(.*?)(=~|==?|!=|!~)(.*)$")
	sampling := regexp.MustCompile("((?:(?:[^:]*):)+)([0-9]+)")
	// First we add hosts from the command line, in all modes
	commands := make([]command, 0)
	add := true
hostspecLoop:
	for len(filters) > 0 {
		glob := filters[0]
		// Do we have a glob or not?
		if comparison.MatchString(glob) || sampling.MatchString(glob) {
			glob = "*"
		} else {
			filters = filters[1:]
		}
		attrs := make(katyusha.MatchAttributes, 0)
		sampled := make([]string, 0)
		count := 0
		for i, arg := range filters[:] {
			if arg == "+" || arg == "-" {
				filters = filters[i+1:]
				if add {
					commands = append(commands, addHostsCommand{glob: glob, attributes: attrs, sampled: sampled, count: count})
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
			if sampledAndCount := sampling.FindStringSubmatch(arg); sampledAndCount != nil {
				if len(sampled) != 0 {
					return fmt.Errorf("only one sampling per hostspec allowed")
				}
				add := false
				for _, s := range strings.Split(sampledAndCount[1], ":") {
					if s == "" {
						add = true
					} else if add && len(sampled) > 0 {
						sampled[len(sampled)-1] = fmt.Sprintf("%s:%s", sampled[len(sampled)-1], s)
						add = false
					} else {
						sampled = append(sampled, s)
						add = false
					}
				}
				count64, _ := strconv.ParseInt(sampledAndCount[2], 0, 64)
				count = int(count64)
			} else {
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
		}
		// We've fallen through, so no more hostspecs
		if add {
			commands = append(commands, addHostsCommand{glob: glob, attributes: attrs, sampled: sampled, count: count})
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
	commands, err := parseCode(string(code))
	if err != nil {
		return err
	}
	e.commands = append(e.commands, commands...)
	return nil
}

func (e *ScriptEngine) ParseCodeLine(code string) error {
	commands, err := parseCode(code)
	if err != nil {
		return err
	}
	e.commands = append(e.commands, commands...)
	return nil
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
	e.Ui.Sync()
}

func (e *ScriptEngine) End() {
	e.Runner.End()
	e.Ui.End()
}
