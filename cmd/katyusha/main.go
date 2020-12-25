package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/scripting"
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
		if f := viper.GetString("Profile"); f != "" {
			pfd, err := os.Create(f + ".cpuprofile")
			if err != nil {
				bail("Could not create CPU profile file: ", err)
			}
			if err = pprof.StartCPUProfile(pfd); err != nil {
				bail("Could not start CPU profile: ", err)
			}
			tfd, err := os.Create(f + ".trace")
			if err != nil {
				bail("Could not create trace file: ", err)
			}
			if err = trace.Start(tfd); err != nil {
				bail("Could not start trace: ", err)
			}
			cmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
				pprof.StopCPUProfile()
				pfd.Close()
				trace.Stop()
				tfd.Close()
			}
		}
	},
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
Providers: %s
`,
		rootCmd.HelpTemplate(),
		filepath.Join(currentUser.configDir, "config.yaml"),
		filepath.Join(currentUser.systemConfigDir, "config.yaml"),
		currentUser.dataDir,
		currentUser.historyDir,
		currentUser.cacheDir,
		strings.Join(katyusha.Providers(), ",")))
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Duration("splay", 0, "Wait a random duration up to this argument before and between each host")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 60*time.Second, "Global timeout for commands")
	rootCmd.PersistentFlags().Duration("host-timeout", 10*time.Second, "Per-host timeout for commands")
	rootCmd.PersistentFlags().Duration("connect-timeout", 3*time.Second, "Per-host ssh connect timeout")
	rootCmd.PersistentFlags().Duration("ssh-agent-timeout", 50*time.Millisecond, "SSH agent timeout when checking functionality")
	rootCmd.PersistentFlags().IntP("parallel", "p", 0, "Maximum number of hosts to run on in parallel")
	rootCmd.PersistentFlags().StringP("output", "o", "all", "When to print command output (all at once, per host or per line)")
	rootCmd.PersistentFlags().Bool("no-pager", false, "Disable the use of the pager")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable the use of the colors in the output")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "Log level")
	rootCmd.PersistentFlags().StringSliceP("sort", "s", []string{"name"}, "Sort hosts by these fields before running commands")
	rootCmd.PersistentFlags().Bool("timestamp", false, "In tail mode, prefix each line with the current time")
	rootCmd.PersistentFlags().String("profile", "", "Write profiling and tracing data to files starting with this name")
	rootCmd.PersistentFlags().Bool("refresh", false, "Force caches to be refreshed")
	viper.BindPFlag("Splay", rootCmd.PersistentFlags().Lookup("splay"))
	viper.BindPFlag("Timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("HostTimeout", rootCmd.PersistentFlags().Lookup("host-timeout"))
	viper.BindPFlag("ConnectTimeout", rootCmd.PersistentFlags().Lookup("connect-timeout"))
	viper.BindPFlag("SshAgentTimeout", rootCmd.PersistentFlags().Lookup("ssh-agent-timeout"))
	viper.BindPFlag("Parallel", rootCmd.PersistentFlags().Lookup("parallel"))
	viper.BindPFlag("Output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("Sort", rootCmd.PersistentFlags().Lookup("sort"))
	viper.BindPFlag("NoPager", rootCmd.PersistentFlags().Lookup("no-pager"))
	viper.BindPFlag("NoColor", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("Timestamp", rootCmd.PersistentFlags().Lookup("timestamp"))
	viper.BindPFlag("Profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("Refresh", rootCmd.PersistentFlags().Lookup("refresh"))
}

func initConfig() {
	viper.AddConfigPath(currentUser.configDir)
	viper.AddConfigPath("/etc/katyusha")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("katyusha")

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
	ui := katyusha.NewSimpleUI()
	ui.SetOutputMode(viper.Get("Output").(katyusha.OutputMode))
	ui.SetOutputTimestamp(viper.GetBool("Timestamp"))
	ui.SetPagerEnabled(!viper.GetBool("NoPager"))
	ui.BindLogrus()

	registry := katyusha.NewRegistry(currentUser.dataDir, currentUser.cacheDir)
	registry.SetSortFields(viper.GetStringSlice("Sort"))
	conf := viper.Sub("Providers")
	if conf != nil {
		if err := registry.LoadProviders(conf); err != nil {
			logrus.Error(err.Error())
			ui.End()
			return nil, err
		}
	}
	registry.LoadMagicProviders()
	if viper.GetBool("Refresh") {
		registry.InvalidateCache()
	}
	if err := registry.LoadHosts(ui.LoadingMessage); err != nil {
		// Do not log this error, registry.LoadHosts() does its own error logging
		ui.End()
		return nil, err
	}
	ui.Sync()
	runner := katyusha.NewRunner(registry)
	runner.SetSplay(viper.GetDuration("Splay"))
	runner.SetParallel(viper.GetInt("Parallel"))
	runner.SetTimeout(viper.GetDuration("Timeout"))
	runner.SetHostTimeout(viper.GetDuration("HostTimeout"))
	runner.SetConnectTimeout(viper.GetDuration("ConnectTimeout"))
	runner.SetSshAgentTimeout(viper.GetDuration("SshAgentTimeout"))
	return scripting.NewScriptEngine(ui, runner), nil
}
