package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/scripting"
	"github.com/sirupsen/logrus"
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
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().DurationP("timeout", "t", 60*time.Second, "Global timeout for commands")
	rootCmd.PersistentFlags().Duration("host-timeout", 10*time.Second, "Per-host timeout for commands")
	rootCmd.PersistentFlags().Duration("connect-timeout", 3*time.Second, "Per-host ssh connect timeout")
	rootCmd.PersistentFlags().IntP("parallel", "p", 0, "Maximum number of hosts to run on in parallel")
	rootCmd.PersistentFlags().StringP("output", "o", "all", "When to print command output (all at once, per host or per line)")
	rootCmd.PersistentFlags().Bool("no-pager", false, "Disable the use of the pager")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "Log level")
	rootCmd.PersistentFlags().StringSliceP("sort", "s", []string{"name"}, "Sort hosts by these fields before running commands")
	viper.BindPFlag("Timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("HostTimeout", rootCmd.PersistentFlags().Lookup("host-timeout"))
	viper.BindPFlag("ConnectTimeout", rootCmd.PersistentFlags().Lookup("connect-timeout"))
	viper.BindPFlag("Parallel", rootCmd.PersistentFlags().Lookup("parallel"))
	viper.BindPFlag("Output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("Sort", rootCmd.PersistentFlags().Lookup("sort"))
	viper.BindPFlag("NoPager", rootCmd.PersistentFlags().Lookup("no-pager"))
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		bail("%s", err)
	}

	// We only need to set defaults for things that don't have a flag bound to them
	root := filepath.Join(home, ".katyusha")
	viper.Set("RootDir", root)
	viper.SetDefault("Formatter", "pretty")

	viper.AddConfigPath(root)
	viper.AddConfigPath("/etc/katyusha")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			bail("Can't read configuration: %s", err)
		}
	}

	viper.SetEnvPrefix("katyusha")
	viper.AutomaticEnv()

	// Check configuration variables
	if _, ok := katyusha.Formatters[viper.GetString("Formatter")]; !ok {
		bail("Unknown formatter: %s. Known formatters: pretty", viper.GetString("Formatter"))
	}

	level, err := logrus.ParseLevel(viper.GetString("LogLevel"))
	if err != nil {
		bail("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", viper.GetString("LogLevel"))
	}
	logrus.SetLevel(level)

	outputModes := map[string]katyusha.OutputMode{
		"all":      katyusha.OutputAll,
		"inline":   katyusha.OutputInline,
		"per-host": katyusha.OutputPerhost,
		"tail":     katyusha.OutputTail,
	}
	om, ok := outputModes[viper.GetString("Output")]
	if !ok {
		bail("Unknown output mode: %s. Known modes: all, inline, per-host, tail", viper.GetString("Output"))
	}
	viper.Set("Output", om)
}

func bail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func setupScriptEngine() (*scripting.ScriptEngine, error) {
	formatter := katyusha.Formatters[viper.GetString("Formatter")]
	ui := katyusha.NewSimpleUI(formatter)
	ui.SetOutputMode(viper.Get("Output").(katyusha.OutputMode))
	ui.SetPagerEnabled(!viper.GetBool("NoPager"))
	logrus.SetFormatter(formatter)
	logrus.SetOutput(ui)

	registry := katyusha.NewRegistry(viper.GetString("RootDir"))
	registry.SetSortFields(viper.GetStringSlice("Sort"))
	registry.LoadMagicProviders()
	conf := viper.Sub("Providers")
	if conf != nil {
		if err := registry.LoadProviders(conf); err != nil {
			logrus.Error(err.Error())
			ui.End()
			return nil, err
		}
	}
	if err := registry.LoadHosts(ui.CacheUpdateChannel()); err != nil {
		// Do not log this error, registry.LoadHosts() does its own error logging
		ui.End()
		return nil, err
	}
	runner := katyusha.NewRunner(registry)
	runner.SetParallel(viper.GetInt("Parallel"))
	runner.SetTimeout(viper.GetDuration("Timeout"))
	runner.SetHostTimeout(viper.GetDuration("HostTimeout"))
	runner.SetConnectTimeout(viper.GetDuration("ConnectTimeout"))
	return scripting.NewScriptEngine(ui, runner), nil
}
