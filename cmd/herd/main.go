package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/seveas/herd"
)

func main() {
	c := herd.NewAppConfig()
	herd.UI = herd.NewSimpleUI()

	pflag.BoolVarP(&c.List, "list", "l", c.List, "List matching hosts (one per line) instead of executing commands")
	pflag.BoolVarP(&c.ListOneline, "list-oneline", "L", c.List, "List matching hosts (all on one line) instead of executing commands")
	pflag.CommandLine.SetOutput(os.Stderr)
	pflag.Parse()
	if c.ListOneline {
		c.List = true
	}

	args := pflag.Args()
	commandStart := pflag.CommandLine.ArgsLenAtDash()
	if !c.List && (commandStart == -1 || commandStart == len(args) || commandStart == 0) {
		usage()
		os.Exit(1)
	}
	if c.List && commandStart == -1 {
		commandStart = len(args)
	}

	hostSpecs := args[:commandStart]
	command := args[commandStart:]

	providers := herd.LoadProviders(c)

	hosts := make(herd.Hosts, 0)

hostspecLoop:
	for true {
		glob := hostSpecs[0]
		attrs := make(herd.HostAttributes)
		for i, arg := range hostSpecs[1:] {
			if arg == "+" {
				hostSpecs = hostSpecs[i+2:]
				hosts = append(hosts, providers.GetHosts(glob, attrs)...)
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
		hosts = append(hosts, providers.GetHosts(glob, attrs)...)
		break
	}

	hosts = hosts.SortAndUniq()

	if c.List {
		if c.ListOneline {
			for i, host := range hosts {
				if i == 0 {
					os.Stdout.WriteString(host.Name)
				} else {
					fmt.Printf(",%s", host.Name)
				}
			}
			os.Stdout.WriteString("\n")
		} else {
			for _, host := range hosts {
				fmt.Println(host.Name)
			}
		}
		return
	}

	runner := herd.NewRunner(hosts)
	hi := runner.Run(strings.Join(command, " "))

	c.Formatter.Format(hi, os.Stdout)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: herd [args] hostlgob [attr=value...] [+ hostglob [attr=value]...] -- command\n\n")
	pflag.CommandLine.PrintDefaults()
}
