package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:                   "list [--oneline] glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Query your datasources and list hosts matching globs and filters",
	Example:               "  herd list *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runList,
}

func init() {
	listCmd.Flags().Bool("oneline", false, "List all hosts on one line, separated by commas")
	listCmd.Flags().StringSlice("attributes", []string{}, "Show not onlt the names, but also the specified attributes")
	listCmd.Flags().Bool("all-attributes", false, "List hosts with all their attributes")
	listCmd.Flags().Bool("csv", false, "Output in csv format")
	viper.BindPFlag("OneLine", listCmd.Flags().Lookup("oneline"))
	viper.BindPFlag("AllAttributes", listCmd.Flags().Lookup("all-attributes"))
	viper.BindPFlag("Attributes", listCmd.Flags().Lookup("attributes"))
	viper.BindPFlag("Csv", listCmd.Flags().Lookup("csv"))
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but list mode doesn't support that")
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
	engine.AddListHostsCommand(viper.GetBool("OneLine"), viper.GetBool("csv"), viper.GetBool("AllAttributes"), viper.GetStringSlice("Attributes"))
	engine.Execute()
	return nil
}
