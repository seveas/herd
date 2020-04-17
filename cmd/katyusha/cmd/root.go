package cmd

import (
	"fmt"
	"os"
	"path"
	"regexp"
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
	rootCmd.PersistentFlags().StringArray("output-filter", []string{}, "Only output results for hosts matching this filter")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "Log level")
	rootCmd.PersistentFlags().StringP("sort", "s", "name", "Sort hosts before running commands")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Don't show status for succesful commands with no output")
	viper.BindPFlag("Timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("HostTimeout", rootCmd.PersistentFlags().Lookup("host-timeout"))
	viper.BindPFlag("ConnectTimeout", rootCmd.PersistentFlags().Lookup("connect-timeout"))
	viper.BindPFlag("Parallel", rootCmd.PersistentFlags().Lookup("parallel"))
	viper.BindPFlag("Output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("LogLevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("Sort", rootCmd.PersistentFlags().Lookup("sort"))
	viper.BindPFlag("Quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		bail("%s", err)
	}

	// We only need to set defaults for things that don't have a flag bound to them
	root := path.Join(home, ".katyusha")
	viper.Set("RootDir", root)
	viper.SetDefault("HistoryDir", path.Join(root, "history"))
	viper.SetDefault("CacheDir", path.Join(root, "cache"))
	viper.SetDefault("Formatter", "pretty")

	viper.AddConfigPath(root)
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
	formatter, ok := katyusha.Formatters[viper.GetString("Formatter")]
	if !ok {
		bail("Unknown formatter: %s. Known formatters: pretty", viper.GetString("Formatter"))
	}

	level, err := logrus.ParseLevel(viper.GetString("LogLevel"))
	if err != nil {
		bail("Unknown loglevel: %s. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR", viper.GetString("LogLevel"))
	}
	logrus.SetLevel(level)

	outputModes := map[string]bool{"all": true, "host": true, "line": true, "pager": true}
	if _, ok := outputModes[viper.GetString("Output")]; !ok {
		bail("Unknown output mode: %s. Known modes: all, host, line, pager", viper.GetString("Output"))
	}
	filters, err := rootCmd.PersistentFlags().GetStringArray("output-filter")
	commands, err := filterCommands(filters)
	if err != nil {
		bail("Invalid filters: %v", filters)
	}

	// Set up the UI
	katyusha.UI = katyusha.NewSimpleUI(formatter)
	outputFilters := make([]katyusha.MatchAttributes, len(commands))
	for i, c := range commands {
		outputFilters[i] = c.(scripting.AddHostsCommand).Attributes
	}
	if viper.GetBool("quiet") {
		if viper.GetString("Output") == "line" {
			outputFilters = append(outputFilters, katyusha.MatchAttributes{{Name: "err", Value: nil, Negate: true}, {Name: "stderr", Value: regexp.MustCompile("\\S"), Regex: true, Negate: true}})
		} else {
			outputFilters = append(outputFilters, katyusha.MatchAttributes{{Name: "err", Value: nil, Negate: true}})
			outputFilters = append(outputFilters, katyusha.MatchAttributes{{Name: "stdout", Value: regexp.MustCompile("\\S"), Regex: true}})
			outputFilters = append(outputFilters, katyusha.MatchAttributes{{Name: "stderr", Value: regexp.MustCompile("\\S"), Regex: true}})
		}
	}
	katyusha.UI.SetOutputFilter(outputFilters)

	logrus.SetOutput(katyusha.UI)
	logrus.SetFormatter(formatter)
}

func bail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
