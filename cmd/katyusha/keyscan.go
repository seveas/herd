package main

import (
	"fmt"

	"github.com/seveas/katyusha"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var keyScanCmd = &cobra.Command{
	Use:                   "keyscan glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Scan ssh keys and output them in known_hosts format, similar to ssh-keyscan",
	Example:               "  katyusha keyscan *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runKeyScan,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !rootCmd.PersistentFlags().Lookup("loglevel").Changed {
			logrus.SetLevel(logrus.WarnLevel)
		}
	},
}

func init() {
	rootCmd.AddCommand(keyScanCmd)
}

func runKeyScan(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but keyscan mode doesn't support that")
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
	engine.Execute()
	if len(args) == 0 {
		engine.Runner.AddHosts("*", []katyusha.MatchAttribute{})
	}
	engine.Runner.Run("katyusha:connect", nil, nil)
	engine.Runner.RemoveHosts("*", []katyusha.MatchAttribute{{Name: "sshKey", Value: nil}})
	template := `{{ $host := . }}{{ range $key := .PublicKeys -}}
{{ $host.Name }}{{ if $host.Address }},{{ $host.Address }}{{ end }} {{ sshkey $key }}
{{ end -}}
`
	engine.Ui.PrintHostList(engine.Runner.GetHosts(), katyusha.HostListOptions{Template: template})
	return nil
}
