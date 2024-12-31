package main

import (
	"fmt"

	"github.com/seveas/herd/ssh"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	splitAt := cmd.ArgsLenAtDash()
	var filters []string
	if splitAt != -1 {
		filters = args[:splitAt]
		args = args[splitAt:]
	}
	if len(args) != 1 {
		return fmt.Errorf("A single script must be provided")
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	executor, err := ssh.NewExecutor(viper.GetInt("SshAgentCount"), viper.GetDuration("SshAgentTimeout"), *currentUser.user, false)
	if err != nil {
		bail(err.Error())
	}
	engine, err := setupScriptEngine(executor)
	if err != nil {
		return err
	}
	defer engine.End()
	if err = engine.ParseCommandLine(filters, -1); err != nil {
		logrus.Error(err.Error())
		return err
	}
	if err = engine.ParseScriptFile(args[0]); err != nil {
		logrus.Errorf("Unable to parse script %s: %s", args[0], err)
		return err
	}
	engine.Execute()
	fn := historyFile(currentUser.historyDir)
	return engine.History.Save(fn)
}
