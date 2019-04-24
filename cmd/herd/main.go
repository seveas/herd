package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/pflag"

	"github.com/seveas/herd"
)

func main() {
	c := herd.NewAppConfig()
	herd.UI = herd.NewSimpleUI(&c.UI)

	pflag.BoolVarP(&c.List, "list", "l", c.List, "List matching hosts (one per line) instead of executing commands")
	pflag.BoolVarP(&c.ListOneline, "list-oneline", "L", c.List, "List matching hosts (all on one line) instead of executing commands")
	pflag.DurationVar(&c.Runner.Timeout, "timeout", c.Runner.Timeout, "Global timeout for commands")
	pflag.DurationVar(&c.Runner.HostTimeout, "host-timeout", c.Runner.HostTimeout, "Per-host timeout for commands")
	pflag.DurationVar(&c.Runner.ConnectTimeout, "connect-timeout", c.Runner.ConnectTimeout, "SSH connection timeout for commands")
	pflag.IntVarP(&c.Runner.Parallel, "parallel", "p", c.Runner.Parallel, "Maximum number of hosts to run on in parallel")
	pflag.StringVarP(&c.ScriptFile, "script", "s", c.ScriptFile, "Script file to execute")
	pflag.BoolVarP(&c.Interactive, "interactive", "i", c.Interactive, "Interactive mode")
	pflag.CommandLine.SetOutput(os.Stderr)
	pflag.Parse()

	args := pflag.Args()
	commandStart := pflag.CommandLine.ArgsLenAtDash()
	haveCommand := commandStart != -1
	hostSpecs := args
	command := []string{}
	if haveCommand {
		hostSpecs = args[:commandStart]
		command = args[commandStart:]
	}

	// We can have only one mode: list, command-line command, script, or interactive
	modes := 0
	if c.List {
		modes++
	}
	if c.ListOneline {
		modes++
	}
	if haveCommand {
		modes++
	}
	if c.ScriptFile != "" {
		modes++
	}
	if c.Interactive {
		modes++
	}
	if modes != 1 {
		usage()
		os.Exit(1)
	}

	// If we have a command, or a list we must have hostspecs
	if (c.List || haveCommand) && len(hostSpecs) == 0 {
		usage()
		os.Exit(1)
	}

	// First we add hosts from the command line, in all modes
	commands := make([]herd.Command, 0)
	add := true
hostspecLoop:
	for len(hostSpecs) > 0 {
		glob := hostSpecs[0]
		attrs := make(herd.HostAttributes)
		for i, arg := range hostSpecs[1:] {
			if arg == "+" || arg == "-" {
				hostSpecs = hostSpecs[i+2:]
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
			if len(parts) != 2 {
				usage()
				os.Exit(1)
			}
			attrs[parts[0]] = parts[1]
		}
		// We've fallen through, so no more hostspecs
		if add {
			commands = append(commands, herd.AddHostsCommand{Glob: glob, Attributes: attrs})
		} else {
			commands = append(commands, herd.RemoveHostsCommand{Glob: glob, Attributes: attrs})
		}
		break
	}

	// Add a command specified on the command line, if we have one
	if haveCommand {
		commands = append(commands, herd.RunCommand{Command: strings.Join(command, " ")})
	}

	// When we have a script, parse it
	if c.ScriptFile != "" {
		var err error
		scriptCommands, err := herd.ParseScript(c.ScriptFile, &c)
		if err != nil {
			herd.UI.Errorf("Unable to parse script %s: %s", c.ScriptFile, err)
			os.Exit(1)
		}
		for _, command := range scriptCommands {
			commands = append(commands, command)
		}
	}

	// Display list of hosts if requested
	if c.List {
		commands = append(commands, herd.ListHostsCommand{OneLine: false})
	}
	if c.ListOneline {
		commands = append(commands, herd.ListHostsCommand{OneLine: true})
	}

	// Execute commands
	providers := herd.LoadProviders(c)
	runner := herd.NewRunner(providers, &c.Runner)

	for _, command := range commands {
		herd.UI.Debugf("%s", command)
		command.Execute(runner)
	}

	// Enter interactive mode if requested
	if c.Interactive {
		il := &InteractiveLoop{Config: &c, Runner: runner}
		il.Run()
	}

	runner.End()

	// Save history, if there is any
	if len(runner.History) > 0 {
		if err := os.MkdirAll(c.HistoryDir, 0700); err != nil {
			herd.UI.Warnf("Unable to create history path %s: %s", c.HistoryDir, err)
		} else {
			fn := path.Join(c.HistoryDir, runner.History[0].StartTime.Format("2006-01-02T15:04:05.json"))
			runner.History.Save(fn)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: herd [opts] hostglob [attr=value...] [+ hostglob [attr=value]...] -- command\n")
	fmt.Fprintf(os.Stderr, "       herd [opts] --list[-oneline] hostglob [attr=value...] [+ hostglob [attr=value]...]\n")
	fmt.Fprintf(os.Stderr, "       herd [opts] --script scriptfile [hostglob [attr=value...] [+ hostglob [attr=value]...]]\n\n")
	fmt.Fprintf(os.Stderr, "       herd [opts] --interactive [hostglob [attr=value...] [+ hostglob [attr=value]...]]\n\n")
	pflag.CommandLine.PrintDefaults()
}
