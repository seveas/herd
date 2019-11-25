package cmd

import (
	"fmt"

	"github.com/seveas/herd"
	"github.com/seveas/herd/scripting"
	"github.com/spf13/cobra"
)

var runScriptCmd = &cobra.Command{
	Use:   "run-script script [glob [filters] [<+|-> glob [filters]...]]",
	Short: "Run a script on a set of hosts",
	Long: `Herd's scripted mode lets you run multiple commands, also allowing you to manipulate
the host list between commands.`,
	Example: `  herd run-script myscript

  #!/usr/local/bin/herd
  add hosts *.site1.example.com
  run id seveas
  remove hosts exitstatus=1
  run userdel seveas`,
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

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	scriptCommands, err := scripting.ParseScript(args[0])
	if err != nil {
		// This should not show the usage message
		herd.UI.Errorf("Unable to parse script %s: %s", args[0], err)
		return err
	}
	for _, command := range scriptCommands {
		commands = append(commands, command)
	}

	_, err = runCommands(commands, true)
	return err
}
