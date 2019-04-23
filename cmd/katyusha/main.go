package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/pflag"

	"github.com/seveas/katyusha"
)

func main() {
	c := katyusha.NewAppConfig()
	katyusha.UI = katyusha.NewSimpleUI(&c.UI)

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
	commands := make([]katyusha.Command, 0)
hostspecLoop:
	for len(hostSpecs) > 0 {
		glob := hostSpecs[0]
		attrs := make(katyusha.HostAttributes)
		for i, arg := range hostSpecs[1:] {
			if arg == "+" {
				hostSpecs = hostSpecs[i+2:]
				commands = append(commands, katyusha.AddHostsCommand{Glob: glob, Attributes: attrs})
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
		commands = append(commands, katyusha.AddHostsCommand{Glob: glob, Attributes: attrs})
		break
	}

	// Add a command specified on the command line, if we have one
	if haveCommand {
		commands = append(commands, katyusha.RunCommand{Command: strings.Join(command, " ")})
	}

	// When we have a script, parse it
	if c.ScriptFile != "" {
		var err error
		scriptCommands, err := katyusha.ParseScript(c.ScriptFile, &c)
		if err != nil {
			katyusha.UI.Errorf("Unable to parse script %s: %s", c.ScriptFile, err)
			os.Exit(1)
		}
		for _, command := range scriptCommands {
			commands = append(commands, command)
		}
	}

	// Execute commands
	providers := katyusha.LoadProviders(c)
	runner := katyusha.NewRunner(providers, &c.Runner)

	for _, command := range commands {
		katyusha.UI.Debugf("%s", command)
		command.Execute(runner)
	}

	// Enter interactive mode if requested
	if c.Interactive {
		il := &InteractiveLoop{Config: &c, Runner: runner}
		il.Run()
	}

	runner.End()

	// Display list of hosts if requested
	if c.List {
		for _, host := range runner.Hosts {
			fmt.Println(host.Name)
		}
	}
	if c.ListOneline {
		for i, host := range runner.Hosts {
			if i == 0 {
				os.Stdout.WriteString(host.Name)
			} else {
				fmt.Printf(",%s", host.Name)
			}
		}
		os.Stdout.WriteString("\n")
	}

	// Save history, if there is any
	if len(runner.History) > 0 {
		if err := os.MkdirAll(c.HistoryDir, 0700); err != nil {
			katyusha.UI.Warnf("Unable to create history path %s: %s", c.HistoryDir, err)
		} else {
			fn := path.Join(c.HistoryDir, runner.History[0].StartTime.Format("2006-01-02T15:04:05.json"))
			runner.History.Save(fn)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: katyusha [opts] hostglob [attr=value...] [+ hostglob [attr=value]...] -- command\n")
	fmt.Fprintf(os.Stderr, "       katyusha [opts] --list[-oneline] hostglob [attr=value...] [+ hostglob [attr=value]...]\n")
	fmt.Fprintf(os.Stderr, "       katyusha [opts] --script scriptfile [hostglob [attr=value...] [+ hostglob [attr=value]...]]\n\n")
	fmt.Fprintf(os.Stderr, "       katyusha [opts] --interactive [hostglob [attr=value...] [+ hostglob [attr=value]...]]\n\n")
	pflag.CommandLine.PrintDefaults()
}
