package cmd

import (
	"fmt"

	"github.com/seveas/katyusha"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:                   "list [--oneline] glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Query your datasources and hosts matching globs and filters",
	Example:               "  katyusha list *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runList,
}

func init() {
	listCmd.Flags().Bool("oneline", false, "List all hosts on one line, separated by commas")
	viper.BindPFlag("OneLine", listCmd.Flags().Lookup("oneline"))
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
	commands = append(commands, katyusha.ListHostsCommand{OneLine: viper.GetBool("OneLine")})
	runCommands(commands, true)
	return nil
}
