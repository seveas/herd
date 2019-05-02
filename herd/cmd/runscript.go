package cmd

import (
	"fmt"

	"github.com/seveas/herd"
	"github.com/spf13/cobra"
)

var runScriptCmd = &cobra.Command{
	Use:                   "run-script script [glob [filters]]",
	Short:                 "Run a script",
	RunE:                  runScript,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(runScriptCmd)
}

func runScript(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("No script provided")
	}

	commands, err := filterCommands(args[1:])
	if err != nil {
		return err
	}

	scriptCommands, err := herd.ParseScript(args[0])
	if err != nil {
		// This should not show the usage message
		herd.UI.Errorf("Unable to parse script %s: %s", args[0], err)
		return nil
	}
	for _, command := range scriptCommands {
		commands = append(commands, command)
	}

	runCommands(commands, true)
	return nil
}
