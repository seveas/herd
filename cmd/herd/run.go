package main

import (
	"fmt"

	"github.com/seveas/herd/ssh"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:                   "run glob [filters] [<+|-> glob [filters]...] -- command [args...]",
	Short:                 "Run a single command on a set of hosts",
	Example:               "  herd run *.site1.example.com os=Debian + *.site2.example.com os=Debian - '*' status=live -- sudo apt-get install bash",
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

	executor, err := ssh.NewExecutor(viper.GetInt("SshAgentCount"), viper.GetDuration("SshAgentTimeout"), *currentUser.user, true)
	if err != nil {
		bail(err.Error())
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
	fn := historyFile(currentUser.historyDir)
	return engine.History.Save(fn)
}
