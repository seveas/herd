package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:                   "run glob [filters] [<+|-> glob [filters]...] -- command [args...]",
	Short:                 "Run a single command on a set of hosts",
	Example:               "  katyusha run *.site1.example.com os=Debian + *.site2.example.com os=Debian - '*' status=live -- sudo apt-get install bash",
	RunE:                  runCommand,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runCommand(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt == -1 {
		return fmt.Errorf("A command is mandatory")
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	engine, err := setupScriptEngine()
	if err != nil {
		return err
	}
	if err = engine.ParseCommandLine(args, splitAt); err != nil {
		logrus.Error(err.Error())
		return err
	}
	engine.Execute()
	return engine.End()
}
