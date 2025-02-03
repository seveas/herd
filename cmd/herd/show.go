package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showCmd = &cobra.Command{
	Use:                   "show hostname",
	Short:                 "Show all attributes of one or more hosts",
	Example:               "  herd show host1.site1.example.com",
	DisableFlagsInUseLine: true,
	RunE:                  runShow,
}

func init() {
	f := showCmd.Flags()
	bindFlagsAndEnv(f)
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	viper.Set("Template", "{{.|yaml}}")
	return runList(cmd, args)
}
