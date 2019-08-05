package cmd

import (
	"fmt"
	"os"
	"path"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/seveas/katyusha"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "katyusha",
	Long: `Replace your ssh for loops with a tool that
- can find hosts for you
- can handle thousands of hosts in parallel
- does not fork a command for every host
- stores all its history, including output, for you to reuse
- can run interactively!
`,
	Example: `  katyusha run '*' os=Debian -- dpkg -l bash
  katyusha interactive *vpn-gateway*`,
	Args:    cobra.NoArgs,
	Version: katyusha.Version(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		katyusha.UI = katyusha.NewSimpleUI()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().DurationP("timeout", "t", 60*time.Second, "Global timeout for commands")
	rootCmd.PersistentFlags().Duration("host-timeout", 10*time.Second, "Per-host timeout for commands")
	rootCmd.PersistentFlags().Duration("connect-timeout", 3*time.Second, "Per-host ssh connect timeout")
	rootCmd.PersistentFlags().IntP("parallel", "p", 0, "Maximum number of hosts to run on in parallel")
	rootCmd.PersistentFlags().StringP("output", "o", "all", "When to print command output (all at once, per host or per line)")
	viper.BindPFlag("Timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("HostTimeout", rootCmd.PersistentFlags().Lookup("host-timeout"))
	viper.BindPFlag("ConnectTimeout", rootCmd.PersistentFlags().Lookup("connect-timeout"))
	viper.BindPFlag("Parallel", rootCmd.PersistentFlags().Lookup("parallel"))
	viper.BindPFlag("Output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		bail("%s", err)
	}

	viper.SetDefault("HistoryDir", path.Join(home, ".katyusha", "history"))
	viper.SetDefault("LogLevel", katyusha.INFO)
	viper.SetDefault("Formatter", "pretty")

	viper.AddConfigPath(path.Join(home, ".katyusha"))
	viper.AddConfigPath("/etc/katyusha")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Can't read configuration:", err)
			os.Exit(1)
		}
	}

	viper.SetEnvPrefix("katyusha")
	viper.AutomaticEnv()

	// Check configuration variables
	if _, ok := katyusha.Formatters[viper.GetString("Formatter")]; !ok {
		fmt.Fprintln(os.Stderr, "Unknown formatter:", viper.GetString("Formatter"))
		os.Exit(1)
	}
	logLevels := map[string]int{"DEBUG": katyusha.DEBUG, "INFO": katyusha.INFO, "NORMAL": katyusha.NORMAL, "WARNING": katyusha.WARNING, "ERROR": katyusha.ERROR}
	if level, ok := logLevels[viper.GetString("LogLevel")]; ok {
		viper.Set("LogLevel", level)
	}
}
