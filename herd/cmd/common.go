package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/seveas/herd"
	"github.com/spf13/cobra"
)

func splitArgs(cmd *cobra.Command, args []string) ([]string, []string) {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt == -1 {
		return args, []string{}
	}
	return args[:splitAt], args[splitAt:]
}

func filterCommands(filters []string) ([]herd.Command, error) {
	comparison := regexp.MustCompile("^(.*?)(=~|==?|!=|!~)(.*)$")
	// First we add hosts from the command line, in all modes
	commands := make([]herd.Command, 0)
	add := true
hostspecLoop:
	for len(filters) > 0 {
		glob := filters[0]
		attrs := make(herd.MatchAttributes, 0)
		for i, arg := range filters[1:] {
			if arg == "+" || arg == "-" {
				filters = filters[i+2:]
				if add {
					commands = append(commands, herd.AddHostsCommand{Glob: glob, Attributes: attrs})
				} else {
					commands = append(commands, herd.RemoveHostsCommand{Glob: glob, Attributes: attrs})
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
				return nil, fmt.Errorf("incorrect filter: %s", arg)
			}
			key, comp, val := parts[1], parts[2], parts[3]
			attr := herd.MatchAttribute{Name: key, Value: val, FuzzyTyping: true}
			if strings.HasPrefix(comp, "!") {
				attr.Negate = true
			}
			if strings.HasSuffix(comp, "~") {
				re, err := regexp.Compile(val)
				if err != nil {
					return nil, fmt.Errorf("Invalid regexp /%s/: %s", val, err)
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
			commands = append(commands, herd.AddHostsCommand{Glob: glob, Attributes: attrs})
		} else {
			commands = append(commands, herd.RemoveHostsCommand{Glob: glob, Attributes: attrs})
		}
		break
	}
	return commands, nil
}

func runCommands(commands []herd.Command, doEnd bool) *herd.Runner {
	providers := herd.LoadProviders()
	runner := herd.NewRunner(providers)

	for _, command := range commands {
		herd.UI.Debugf("%s", command)
		command.Execute(runner)
	}
	if doEnd {
		runner.End()
	}
	return runner
}
