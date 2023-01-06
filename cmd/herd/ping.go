package main

import (
	"fmt"
	"time"

	"github.com/seveas/herd"
	"github.com/seveas/herd/ssh"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pingCmd = &cobra.Command{
	Use:                   "ping [options] glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Connect to hosts and check if they are alive",
	Example:               "  herd ping *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runPing,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !rootCmd.PersistentFlags().Lookup("loglevel").Changed {
			logrus.SetLevel(logrus.WarnLevel)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}

func runPing(cmd *cobra.Command, args []string) error {
	if !rootCmd.PersistentFlags().Lookup("output").Changed {
		viper.Set("Output", herd.OutputInline)
	}
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but ping mode doesn't support that")
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	executor, err := ssh.NewPingExecutor(viper.GetDuration("SshAgentTimeout"), *currentUser.user)
	if err != nil {
		return err
	}
	engine, err := setupScriptEngine(executor)
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
		hosts := engine.Registry.Search("*", []herd.MatchAttribute{}, []string{}, 0)
		engine.Hosts.AddHosts(hosts)
	}
	oc := engine.Ui.OutputChannel()
	pc := engine.Ui.ProgressChannel(time.Now().Add(engine.Runner.GetTimeout()))
	hi, err := engine.Runner.Run("", pc, oc)
	if err != nil {
		logrus.Error(err.Error())
		return nil
	}
	engine.Ui.PrintHistoryItem(hi)
	return nil
}
