package cmd

import (
	"fmt"

	"github.com/seveas/herd/scripting"
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
	filters, rest := splitArgs(cmd, args)
	if len(rest) != 0 {
		return fmt.Errorf("Command provided, but list mode doesn't support that")
	}
	commands, err := filterCommands(filters)
	if err != nil {
		return err
	}
	commands = append(commands, scripting.ListHostsCommand{
		OneLine:       viper.GetBool("OneLine"),
		AllAttributes: viper.GetBool("AllAttributes"),
		Attributes:    viper.GetStringSlice("Attributes"),
		Csv:           viper.GetBool("Csv"),
	})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	_, err = runCommands(commands, true)
	return err
}
