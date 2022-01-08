package main

import (
	"fmt"

	"github.com/seveas/herd"
	"github.com/seveas/herd/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var keyScanCmd = &cobra.Command{
	Use:                   "keyscan [options] glob [filters] [<+|-> glob [filters]...]",
	Short:                 "Scan ssh keys and output them in known_hosts format, similar to ssh-keyscan",
	Example:               "  herd keyscan *.site2.example.com os=Debian",
	DisableFlagsInUseLine: true,
	RunE:                  runKeyScan,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !rootCmd.PersistentFlags().Lookup("loglevel").Changed {
			logrus.SetLevel(logrus.WarnLevel)
		}
	},
}

func init() {
	keyScanCmd.Flags().StringSlice("type", []string{"ssh-rsa", "ecdsa-sha2-nistp256", "ssh-ed25519"}, "Which key algorithm(s) to scan for")
	viper.BindPFlag("KeyType", keyScanCmd.Flags().Lookup("type"))
	rootCmd.AddCommand(keyScanCmd)
}

func runKeyScan(cmd *cobra.Command, args []string) error {
	splitAt := cmd.ArgsLenAtDash()
	if splitAt != -1 {
		return fmt.Errorf("Command provided, but keyscan mode doesn't support that")
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	knownTypes := map[string]string{
		"dsa":                 "ssh-dss",
		"dss":                 "ssh-dss",
		"ssh-dss":             "ssh-dss",
		"rsa":                 "ssh-rsa",
		"ssh-rsa":             "ssh-rsa",
		"ecdsa":               "ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521",
		"ecdsa-sha2-nistp256": "ecdsa-sha2-nistp256",
		"ecdsa-sha2-nistp384": "ecdsa-sha2-nistp384",
		"ecdsa-sha2-nistp521": "ecdsa-sha2-nistp521",
		"ed25519":             "ssh-ed25519",
		"ssh-ed25519":         "ssh-ed25519",
	}
	keyTypes := make([]string, 0)
	for _, keyType := range viper.GetStringSlice("KeyType") {
		if expandedKeyType, ok := knownTypes[keyType]; ok {
			keyTypes = append(keyTypes, expandedKeyType)
		} else {
			return fmt.Errorf("Unknown public key type: %s", keyType)
		}
	}

	engine, err := setupScriptEngine(ssh.NewKeyScanExecutor(keyTypes))
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
		hosts := engine.Registry.GetHosts("*", []herd.MatchAttribute{}, []string{}, 0)
		engine.Runner.AddHosts(hosts)
	}
	engine.Runner.Run("herd:keyscan", nil, nil)
	template := `{{ $host := . }}{{ range $key := .PublicKeys -}}
{{ $host.Name }}{{ if $host.Address }},{{ $host.Address }}{{ end }} {{ sshkey $key }}
{{ end -}}
`
	engine.Ui.PrintHostList(engine.Runner.GetHosts(), herd.HostListOptions{Template: template})
	return nil
}
