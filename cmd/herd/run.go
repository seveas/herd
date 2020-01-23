package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	engine, err := setupScriptEngine()
	if err != nil {
		return err
	}
	defer engine.End()
	if err = engine.ParseCommandLine(args, splitAt); err != nil {
		logrus.Error(err.Error())
		return err
	}
	fn := filepath.Join(currentUser.historyDir, time.Now().Format("2006-01-02T15:04:05.json"))
	engine.Execute()
	return engine.SaveHistory(fn)
}
