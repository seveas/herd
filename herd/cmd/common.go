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
	// First we add hosts from the command line, in all modes
	commands := make([]herd.Command, 0)
	add := true
hostspecLoop:
	for len(filters) > 0 {
		glob := filters[0]
		attrs := make(herd.HostAttributes)
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
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
				return nil, fmt.Errorf("incorrect filter: %s", arg)
			}
			if parts[1][0] == '~' {
				re, err := regexp.Compile(parts[1][1:])
				if err != nil {
					return nil, fmt.Errorf("Invalid regexp /%s/: %s", parts[1][1:], err)
				} else {
					attrs[parts[0]] = re
				}
			} else {
				attrs[parts[0]] = parts[1]
			}
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
