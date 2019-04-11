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

	pflag.BoolVarP(&c.List, "list", "l", c.List, "List matching hosts instead of executing commands")
	pflag.CommandLine.SetOutput(os.Stderr)
	pflag.Parse()

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

	// FIXME proper parsing of hostspecs
	hosts := providers.GetHosts(hostSpecs[0], herd.HostAttributes{})

	if c.List {
		for _, host := range hosts {
			fmt.Println(host.Name)
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
