package main

import (
	"fmt"

	"github.com/seveas/herd"

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
	f := listCmd.Flags()
	f.Bool("oneline", false, "List all hosts on one line, separated by commas")
	f.String("separator", ",", "String separating hostnames in --oneline mode")
	f.Bool("all-attributes", false, "List hosts with all their attributes (deprecated, use --attributes=* instead)")
	f.StringSlice("attributes", []string{}, "Show not only the names, but also the specified attributes (supports wildcards)")
	f.Bool("csv", false, "Output in csv format")
	f.Bool("header", true, "Print attribute names in a header line before printing host data")
	f.String("template", "", "Template to use for showing hosts")
	f.StringSlice("count", []string{}, "Show counts for the values of these attributes")
	f.String("group", "", "Group hosts by the values of this attribute")
	// This makes `--count my_attribute` stop working and makes it require `--count=my_attribute` instead.
	// f.Lookup("count").NoOptDefVal = "*"
	bindFlagsAndEnv(f)
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but list mode doesn't support that")
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	engine, err := setupScriptEngine(nil)
	if err != nil {
		return err
	}
	defer engine.End()
	if len(args) == 0 {
		args = append(args, "*")
	}
	if err = engine.ParseCommandLine(args, splitAt); err != nil {
		logrus.Error(err.Error())
		return err
	}
	// Backwards compatibility code
	if viper.GetBool("AllAttributes") {
		viper.Set("Attributes", []string{"*"})
	}
	engine.Execute()
	opts := herd.HostListOptions{
		OneLine:     viper.GetBool("OneLine"),
		Separator:   viper.GetString("Separator"),
		Csv:         viper.GetBool("csv"),
		Attributes:  viper.GetStringSlice("Attributes"),
		Header:      viper.GetBool("Header"),
		Align:       true,
		Template:    viper.GetString("Template"),
		Count:       viper.GetStringSlice("Count"),
		SortByCount: !viper.IsSet("Sort"),
		Group:       viper.GetString("Group"),
	}
	engine.Ui.PrintHostList(opts)
	return nil
}
