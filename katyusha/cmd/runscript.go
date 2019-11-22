package cmd

import (
	"fmt"

	"github.com/seveas/katyusha"
	"github.com/spf13/cobra"
)

var runScriptCmd = &cobra.Command{
	Use:   "run-script script [glob [filters] [<+|-> glob [filters]...]]",
	Short: "Run a script on a set of hosts",
	Long: `Katyusha's scripted mode lets you run multiple commands, also allowing you to manipulate
the host list between commands.`,
	Example: `  katyusha run-script myscript

  #!/usr/local/bin/katyusha
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
	scriptCommands, err := katyusha.ParseScript(args[0])
	if err != nil {
		// This should not show the usage message
		katyusha.UI.Errorf("Unable to parse script %s: %s", args[0], err)
		return nil
	}
	for _, command := range scriptCommands {
		commands = append(commands, command)
	}

	_, err = runCommands(commands, true)
	return err
}
