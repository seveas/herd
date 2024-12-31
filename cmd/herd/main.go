package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"runtime/trace"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/seveas/herd"
	"github.com/seveas/herd/scripting"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	_ = os.MkdirAll(u.configDir, 0o700)
	_ = os.MkdirAll(u.systemConfigDir, 0o755)
	_ = os.MkdirAll(u.dataDir, 0o700)
	_ = os.MkdirAll(u.cacheDir, 0o700)
	_ = os.MkdirAll(u.historyDir, 0o700)
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
		strings.Join(herd.Providers(), ",")))
	cobra.OnInitialize(initConfig)
	f := rootCmd.PersistentFlags()
	f.Duration("splay", 0, "Wait a random duration up to this argument before and between each host")
	f.DurationP("timeout", "t", 5*time.Minute, "Global timeout for commands")
	f.Duration("load-timeout", 30*time.Second, "Timeout for loading host data from providers")
	f.Duration("host-timeout", time.Minute, "Per-host timeout for commands")
	f.Duration("connect-timeout", 15*time.Second, "Per-host ssh connect timeout")
	f.Duration("ssh-agent-timeout", time.Second, "SSH agent timeout when checking functionality")
	f.Int("ssh-agent-count", 50, "Number of parallel connections to the ssh agent")
	f.IntP("parallel", "p", 0, "Maximum number of hosts to run on in parallel")
	f.StringP("output", "o", "all", "When to print command output (all at once, per host or per line)")
	f.Bool("no-pager", false, "Disable the use of the pager")
	f.Bool("no-color", false, "Disable the use of the colors in the output")
	f.StringP("loglevel", "l", "INFO", "Log level")
	f.StringSliceP("sort", "s", []string{"name"}, "Sort hosts by these fields before running commands")
	f.Bool("timestamp", false, "In tail mode, prefix each line with the current time")
	f.String("profile", "", "Write profiling and tracing data to files starting with this name")
	f.Bool("refresh", false, "Force caches to be refreshed")
	f.Bool("no-refresh", false, "Do not try to refresh cached data")
	f.Bool("strict-loading", false, "Fail if any provider fails to load data")
	f.Bool("no-magic-providers", false, "Do not use magic autodiscovery, only explicitly configured providers")
	bindFlagsAndEnv(f)
}

func bindFlagsAndEnv(s *pflag.FlagSet) {
	rx := regexp.MustCompile("((?:^|-).)")
	toUpper := func(s string) string {
		return strings.ToUpper(strings.Trim(s, "-"))
	}
	s.VisitAll(func(f *pflag.Flag) {
		varName := rx.ReplaceAllStringFunc(f.Name, toUpper)
		envName := "HERD_" + strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_")
		if err := viper.BindPFlag(varName, f); err != nil {
			panic(err)
		}
		if err := viper.BindEnv(varName, envName); err != nil {
			panic(err)
		}
	})
}

func initConfig() {
	viper.AddConfigPath(currentUser.configDir)
	viper.AddConfigPath("/etc/herd")
	viper.SetConfigName("config")

	// Read the configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			bail("Can't read configuration: %s", err)
		}
	}

	// Check configuration variables

	// Limit concurrent ssh connections to parallelism
	if viper.GetInt("Parallel") != 0 && viper.GetInt("Parallel") < viper.GetInt("SshAgentCount") {
		viper.Set("SshAgentCount", viper.GetInt("Parallel"))
	}

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

func setupScriptEngine(executor herd.Executor) (*scripting.ScriptEngine, error) {
	hosts := new(herd.HostSet)
	hosts.SetSortFields(viper.GetStringSlice("Sort"))
	colors := viper.Sub("Colors")
	colorConfig := herd.ColorConfig{}
	if colors != nil {
		colorConfig.LogDebug = colors.GetString("LogDebug")
		colorConfig.LogInfo = colors.GetString("LogInfo")
		colorConfig.LogWarn = colors.GetString("LogWarn")
		colorConfig.LogError = colors.GetString("LogError")
		colorConfig.Command = colors.GetString("Command")
		colorConfig.Summary = colors.GetString("Summary")
		colorConfig.Provider = colors.GetString("Provider")
		colorConfig.HostStdout = colors.GetString("HostStdout")
		colorConfig.HostStderr = colors.GetString("HostStderr")
		colorConfig.HostOK = colors.GetString("HostOK")
		colorConfig.HostFail = colors.GetString("HostFail")
		colorConfig.HostError = colors.GetString("HostError")
		colorConfig.HostCancel = colors.GetString("HostCancel")
	}

	ui := herd.NewSimpleUI(colorConfig, hosts)
	ui.SetOutputMode(viper.Get("Output").(herd.OutputMode))
	ui.SetOutputTimestamp(viper.GetBool("Timestamp"))
	ui.SetPagerEnabled(!viper.GetBool("NoPager"))
	ui.BindLogrus()

	registry := herd.NewRegistry(currentUser.dataDir, currentUser.cacheDir)
	conf := viper.Sub("Providers")
	if conf != nil {
		if err := registry.LoadProviders(conf); err != nil {
			logrus.Error(err.Error())
			ui.End()
			return nil, err
		}
	}
	if !viper.GetBool("NoMagicProviders") {
		registry.LoadMagicProviders()
	}
	if viper.GetBool("Refresh") {
		registry.InvalidateCache()
	}
	if viper.GetBool("NoRefresh") {
		registry.KeepCaches()
	}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("LoadTimeout"))
	defer cancel()
	if err := registry.LoadHosts(ctx, ui.LoadingMessage); err != nil {
		// Do not log this error, registry.LoadHosts() does its own error logging
		if viper.GetBool("StrictLoading") {
			ui.End()
			return nil, err
		}
	}
	if err := registry.LoadHostKeys(ctx, ui.LoadingMessage); err != nil {
		if viper.GetBool("StrictLoading") {
			ui.End()
			return nil, err
		}
	}
	ui.Sync()
	runner := herd.NewRunner(hosts, executor)
	handleSignals(runner)
	runner.SetSplay(viper.GetDuration("Splay"))
	runner.SetParallel(viper.GetInt("Parallel"))
	// If only one of the timeouts was specified, we adjust the other based on batch size
	if viper.IsSet("Timeout") || !viper.IsSet("HostTimeout") {
		runner.SetTimeout(viper.GetDuration("Timeout"))
	}
	if viper.IsSet("HostTimeout") || !viper.IsSet("Timeout") {
		runner.SetHostTimeout(viper.GetDuration("HostTimeout"))
	}
	runner.SetConnectTimeout(viper.GetDuration("ConnectTimeout"))
	return scripting.NewScriptEngine(hosts, ui, registry, runner), nil
}

func historyFile(dir string) string {
	// Read existing history, migrate if needed
	hist, err := os.ReadDir(dir)
	if err != nil {
		logrus.Warnf("Could not read history directory %s: %s", dir, err)
		return filepath.Join(dir, time.Now().Format("2006-01-02_150405.json"))
	}
	mustMigrate := false
	seq := 0
	hist = slices.DeleteFunc(hist, func(e os.DirEntry) bool {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			return true
		}
		parts := strings.Split(e.Name(), "_")
		if len(parts) != 3 {
			mustMigrate = true
			return false
		}
		s, err := strconv.Atoi(parts[0])
		if err != nil {
			logrus.Warnf("Could not parse history file name %s: %s", hist[len(hist)-1].Name(), err)
			return true
		}
		if s > seq {
			seq = s
		}
		return false
	})
	if mustMigrate {
		migrateHistory(dir, hist)
		hist, _ = os.ReadDir(dir)
		hist = slices.DeleteFunc(hist, func(e os.DirEntry) bool {
			return e.IsDir() || !strings.HasSuffix(e.Name(), ".json")
		})
	}

	return filepath.Join(dir, fmt.Sprintf("%d_%s.json", seq+1, time.Now().Format("2006-01-02_150405")))
}

func migrateHistory(dir string, hist []os.DirEntry) {
	logrus.Warnf("Migrating history in %s", dir)
	type entry struct {
		oldName string
		newName string
	}
	entries := make([]entry, 0, len(hist))
	for i, e := range hist {
		parts := strings.Split(e.Name(), "_")
		if len(parts) != 2 && len(parts) != 3 {
			logrus.Warnf("Could not parse history file name %s, skipping migration", e.Name())
			continue
		}
		entries = append(entries, entry{
			oldName: e.Name(),
			newName: fmt.Sprintf("%d_%s_%s", i+1, parts[len(parts)-2], parts[len(parts)-1]),
		})
	}
	for _, e := range entries {
		if e.oldName == e.newName {
			continue
		}
		logrus.Infof("Migrating history: %s => %s\n", e.oldName, e.newName)
		if err := os.Rename(filepath.Join(dir, e.oldName), filepath.Join(dir, e.newName)); err != nil {
			logrus.Warnf("Could not rename %s to %s: %s", e.oldName, e.newName, err)
		}
	}
}
