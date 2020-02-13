package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/herd"
	"github.com/seveas/herd/scripting"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var currentUser *userData

type userData struct {
	user            *user.User
	cacheDir        string
	configDir       string
	systemConfigDir string
	dataDir         string
	historyDir      string
}

func (u *userData) makeDirectories() {
	// We ignore errors, as we should function fine without these
	os.MkdirAll(u.configDir, 0700)
	os.MkdirAll(u.systemConfigDir, 0755)
	os.MkdirAll(u.dataDir, 0700)
	os.MkdirAll(u.cacheDir, 0700)
	os.MkdirAll(u.historyDir, 0700)
}

var rootCmd = &cobra.Command{
	Use: "herd",
	Long: `Replace your ssh for loops with a tool that
- can find hosts for you
- can handle thousands of hosts in parallel
- does not fork a command for every host
- stores all its history, including output, for you to reuse
- can run interactively!
`,
	Example: `  herd run '*' os=Debian -- dpkg -l bash
  herd interactive *vpn-gateway*`,
	Args:    cobra.NoArgs,
	Version: herd.Version(),
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error
	currentUser, err = getCurrentUser()
	if err != nil {
		bail("%s", err)
	}
	currentUser.makeDirectories()
	rootCmd.SetHelpTemplate(fmt.Sprintf(`%s

Configuration: %s, %s
Datadir: %s
History: %s
Cache: %s
`,
		rootCmd.HelpTemplate(),
		filepath.Join(currentUser.configDir, "config.yaml"),
		filepath.Join(currentUser.systemConfigDir, "config.yaml"),
		currentUser.dataDir,
		currentUser.historyDir,
		currentUser.cacheDir))
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Duration("splay", 0, "Wait a random duration up to this argument before and between each host")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 60*time.Second, "Global timeout for commands")
	rootCmd.PersistentFlags().Duration("host-timeout", 10*time.Second, "Per-host timeout for commands")
	rootCmd.PersistentFlags().Duration("connect-timeout", 3*time.Second, "Per-host ssh connect timeout")
	rootCmd.PersistentFlags().IntP("parallel", "p", 0, "Maximum number of hosts to run on in parallel")
	rootCmd.PersistentFlags().StringP("output", "o", "all", "When to print command output (all at once, per host or per line)")
	rootCmd.PersistentFlags().Bool("no-pager", false, "Disable the use of the pager")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable the use of the colors in the output")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "Log level")
	rootCmd.PersistentFlags().StringSliceP("sort", "s", []string{"name"}, "Sort hosts by these fields before running commands")
	rootCmd.PersistentFlags().Bool("timestamp", false, "In tail mode, prefix each line with the current time")
	viper.BindPFlag("Splay", rootCmd.PersistentFlags().Lookup("splay"))
	viper.BindPFlag("Timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("HostTimeout", rootCmd.PersistentFlags().Lookup("host-timeout"))
	viper.BindPFlag("ConnectTimeout", rootCmd.PersistentFlags().Lookup("connect-timeout"))
	viper.BindPFlag("Parallel", rootCmd.PersistentFlags().Lookup("parallel"))
	viper.BindPFlag("Output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("Sort", rootCmd.PersistentFlags().Lookup("sort"))
	viper.BindPFlag("NoPager", rootCmd.PersistentFlags().Lookup("no-pager"))
	viper.BindPFlag("NoColor", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("Timestamp", rootCmd.PersistentFlags().Lookup("timestamp"))
}

func initConfig() {
	viper.AddConfigPath(currentUser.configDir)
	viper.AddConfigPath("/etc/herd")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("herd")

	// Read the configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			bail("Can't read configuration: %s", err)
		}
	}
	viper.AutomaticEnv()

	// Check configuration variables
	level, err := logrus.ParseLevel(viper.GetString("LogLevel"))
	if err != nil {
		bail("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", viper.GetString("LogLevel"))
	}
	logrus.SetLevel(level)

	if viper.GetBool("NoColor") {
		ansi.DisableColors(true)
	}
	outputModes := map[string]herd.OutputMode{
		"all":      herd.OutputAll,
		"inline":   herd.OutputInline,
		"per-host": herd.OutputPerhost,
		"tail":     herd.OutputTail,
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
	ui := herd.NewSimpleUI()
	ui.SetOutputMode(viper.Get("Output").(herd.OutputMode))
	ui.SetOutputTimestamp(viper.GetBool("Timestamp"))
	ui.SetPagerEnabled(!viper.GetBool("NoPager"))
	ui.BindLogrus()

	registry := herd.NewRegistry(currentUser.dataDir, currentUser.cacheDir)
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
	var mc chan herd.CacheMessage
	if logrus.IsLevelEnabled(logrus.InfoLevel) {
		mc = ui.CacheUpdateChannel()
	}
	if err := registry.LoadHosts(mc); err != nil {
		// Do not log this error, registry.LoadHosts() does its own error logging
		ui.End()
		return nil, err
	}
	ui.Sync()
	runner := herd.NewRunner(registry)
	runner.SetSplay(viper.GetDuration("Splay"))
	runner.SetParallel(viper.GetInt("Parallel"))
	runner.SetTimeout(viper.GetDuration("Timeout"))
	runner.SetHostTimeout(viper.GetDuration("HostTimeout"))
	runner.SetConnectTimeout(viper.GetDuration("ConnectTimeout"))
	return scripting.NewScriptEngine(ui, runner), nil
}
