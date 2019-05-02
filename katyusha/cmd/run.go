package cmd

import (
	"fmt"
	"strings"

	"github.com/seveas/katyusha"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:                   "run glob [filters] -- command [args...]",
	Short:                 "Run a single command",
	RunE:                  runCommand,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runCommand(cmd *cobra.Command, args []string) error {
	filters, rest := splitArgs(cmd, args)
	if len(rest) == 0 {
		return fmt.Errorf("A command is mandatory")
	}
	commands, err := filterCommands(filters)
	if err != nil {
		return err
	}
	commands = append(commands, katyusha.RunCommand{Command: strings.Join(rest, " ")})
	runCommands(commands, true)
	return nil
}
