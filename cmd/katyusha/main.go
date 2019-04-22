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
	katyusha.UI = katyusha.NewSimpleUI()

	pflag.BoolVarP(&c.List, "list", "l", c.List, "List matching hosts (one per line) instead of executing commands")
	pflag.BoolVarP(&c.ListOneline, "list-oneline", "L", c.List, "List matching hosts (all on one line) instead of executing commands")
	pflag.DurationVar(&c.Runner.Timeout, "timeout", c.Runner.Timeout, "Global timeout for commands")
	pflag.DurationVar(&c.Runner.HostTimeout, "host-timeout", c.Runner.HostTimeout, "Per-host timeout for commands")
	pflag.DurationVar(&c.Runner.ConnectTimeout, "connect-timeout", c.Runner.ConnectTimeout, "SSH connection timeout for commands")
	pflag.IntVarP(&c.Runner.Parallel, "parallel", "p", c.Runner.Parallel, "Maximum number of hosts to run on in parallel")
	pflag.StringVarP(&c.ScriptFile, "script", "s", c.ScriptFile, "Script file to execute")
	pflag.CommandLine.SetOutput(os.Stderr)
	pflag.Parse()
	if c.ListOneline {
		c.List = true
	}

	args := pflag.Args()
	commandStart := pflag.CommandLine.ArgsLenAtDash()
	if !c.List && c.ScriptFile == "" && (commandStart == -1 || commandStart == len(args) || commandStart == 0) {
		usage()
		os.Exit(1)
	}
	if c.ScriptFile != "" && (c.List || len(args) != 0) {
		usage()
		os.Exit(1)
	}
	if c.List && commandStart == -1 {
		commandStart = len(args)
	}

	commands := make([]katyusha.Command, 0)

	if c.ScriptFile != "" {
		var err error
		commands, err = katyusha.ParseScript(c.ScriptFile, &c)
		if err != nil {
			katyusha.UI.Errorf("Unable to parse script %s: %s", c.ScriptFile, err)
			os.Exit(1)
		}
	} else {
		hostSpecs := args[:commandStart]
		command := args[commandStart:]

	hostspecLoop:
		for true {
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
		if len(command) > 0 {
			commands = append(commands, katyusha.RunCommand{Command: strings.Join(command, " "), Formatter: c.Formatter})
		}
	}

	providers := katyusha.LoadProviders(c)
	runner := katyusha.NewRunner(providers, &c.Runner)

	for _, command := range commands {
		katyusha.UI.Debugf("%s", command)
		command.Execute(runner)
	}
	runner.End()

	if c.List {
		if c.ListOneline {
			for i, host := range runner.Hosts {
				if i == 0 {
					os.Stdout.WriteString(host.Name)
				} else {
					fmt.Printf(",%s", host.Name)
				}
			}
			os.Stdout.WriteString("\n")
		} else {
			for _, host := range runner.Hosts {
				fmt.Println(host.Name)
			}
		}
		return
	}

	if err := os.MkdirAll(c.HistoryDir, 0700); err != nil {
		katyusha.UI.Warnf("Unable to create history path %s: %s", c.HistoryDir, err)
	} else {
		fn := path.Join(c.HistoryDir, runner.History[0].StartTime.Format("2006-01-02T15:04:05.json"))
		runner.History.Save(fn)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: katyusha [opts] hostglob [attr=value...] [+ hostglob [attr=value]...] -- command\n")
	fmt.Fprintf(os.Stderr, "       katyusha [opts] --list[-oneline] hostglob [attr=value...] [+ hostglob [attr=value]...]\n")
	fmt.Fprintf(os.Stderr, "       katyusha [opts] --script scriptfile\n\n")
	pflag.CommandLine.PrintDefaults()
}
