package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var keyScanCmd = &cobra.Command{
	Use:                   "keyscan glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Scan ssh keys and output them in known_hosts format, similar to ssh-keyscan",
	Example:               "  katyusha keyscan *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runKeyScan,
}

func init() {
	rootCmd.AddCommand(keyScanCmd)
}

func runKeyScan(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but keyscan mode doesn't support that")
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	engine, err := setupScriptEngine()
	if err != nil {
		return err
	}
	defer engine.End()
	if err = engine.ParseCommandLine(args, splitAt); err != nil {
		logrus.Error(err.Error())
		return err
	}
	if len(engine.ActiveHosts()) == 0 {
		engine.ParseCodeLine("add hosts *\n")
	}
	engine.AddKeyScanCommand()
	engine.Execute()
	return nil
}
