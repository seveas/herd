package main

import (
	"fmt"

	"github.com/seveas/herd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var keyScanCmd = &cobra.Command{
	Use:                   "keyscan glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Scan ssh keys and output them in known_hosts format, similar to ssh-keyscan",
	Example:               "  herd keyscan *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runKeyScan,
}

func init() {
	rootCmd.AddCommand(keyScanCmd)
	cobra.OnInitialize(func() {
		rootCmd.PersistentFlags().Lookup("loglevel").Value.Set("warn")
		rootCmd.PersistentFlags().Lookup("loglevel").DefValue = "warn"
	})
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
	engine.Execute()
	if len(args) == 0 {
		engine.Runner.AddHosts("*", []herd.MatchAttribute{})
	}
	engine.Runner.Run("herd:connect", nil, nil)
	engine.Runner.RemoveHosts("*", []herd.MatchAttribute{{Name: "sshKey", Value: nil}})
	engine.Ui.PrintHostList(engine.Runner.GetHosts(), herd.HostListOptions{Attributes: []string{"sshKey"}})
	return nil
}
