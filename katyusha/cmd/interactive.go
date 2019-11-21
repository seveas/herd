package cmd

import (
	"fmt"
	"io"
	"path"

	"github.com/seveas/katyusha"
	"github.com/seveas/readline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive [glob [filters] [<+|-> glob [filters]...]]",
	Short: "Interactive shell for running commands on a set of hosts",
	Long: `With Katyusha's interactive shell, you can easily run multiple commands, and
manipulate the host list between commands. You can even use the result of
previous commands as filters.`,
	RunE:                  runInteractive,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func runInteractive(cmd *cobra.Command, args []string) error {
	filters, rest := splitArgs(cmd, args)
	if len(rest) > 0 {
		return fmt.Errorf("Command provided, but interactive mode doesn't support that")
	}
	commands, err := filterCommands(filters)
	if err != nil {
		return err
	}
	runner := runCommands(commands, false)
	// Nil return means provider problems
	if runner == nil {
		return nil
	}

	// Enter interactive mode
	il := &InteractiveLoop{Runner: runner}
	il.Run()
	runner.End()

	return nil
}

type InteractiveLoop struct {
	Runner *katyusha.Runner
}

func (l *InteractiveLoop) Run() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          l.Prompt(),
		AutoComplete:    l.AutoComplete(),
		HistoryFile:     path.Join(viper.GetString("HistoryDir"), "interactive"),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		katyusha.UI.Errorf("Unable to start interactive mode: %s", err)
		return
	}
	defer rl.Close()
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		} else if err != nil {
			katyusha.UI.Errorf(err.Error())
			break
		}
		if line == "exit" {
			break
		}
		commands, err := katyusha.ParseCode(line + "\n")
		if err != nil {
			katyusha.UI.Errorf(err.Error())
			continue
		}
		for _, command := range commands {
			katyusha.UI.Debugf("%s", command)
			command.Execute(l.Runner)
			rl.SetPrompt(l.Prompt())
		}
	}
}

func (l *InteractiveLoop) Prompt() string {
	return fmt.Sprintf("katyusha [%d hosts] $ ", len(l.Runner.Hosts))
}

func (l *InteractiveLoop) AutoComplete() readline.AutoCompleter {
	p := readline.PcItem
	return readline.NewPrefixCompleter(
		p("set",
			p("Timeout"),
			p("HostTimeout"),
			p("ConnectTimeout"),
			p("Parallel"),
			p("Output"),
			p("LogLevel"),
		),
		p("add hosts"),
		p("remove hosts"),
		p("list hosts",
			p("oneline"),
		),
		p("run"),
	)
}
