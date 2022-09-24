package main

import (
	"fmt"
	"path/filepath"
	"time"

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

	executor, err := ssh.NewExecutor(viper.GetDuration("SshAgentTimeout"), *currentUser.user)
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
	fn := filepath.Join(currentUser.historyDir, time.Now().Format("2006-01-02_150405.json"))
	engine.Execute()
	return engine.History.Save(fn)
}
