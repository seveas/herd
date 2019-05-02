package cmd

import (
	"fmt"
	"io"
	"path"

	"github.com/chzyer/readline"
	"github.com/seveas/herd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var interactiveCmd = &cobra.Command{
	Use:                   "interactive [filters]",
	Short:                 "Interactive shell",
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

	// Enter interactive mode
	il := &InteractiveLoop{Runner: runner}
	il.Run()
	runner.End()

	return nil
}

type InteractiveLoop struct {
	Runner *herd.Runner
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
		herd.UI.Errorf("Unable to start interactive mode: %s", err)
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
			herd.UI.Errorf(err.Error())
			break
		}
		if line == "exit" {
			break
		}
		commands, err := herd.ParseCode(line + "\n")
		if err != nil {
			herd.UI.Errorf(err.Error())
			continue
		}
		for _, command := range commands {
			herd.UI.Debugf("%s", command)
			command.Execute(l.Runner)
			rl.SetPrompt(l.Prompt())
		}
	}
}

func (l *InteractiveLoop) Prompt() string {
	return fmt.Sprintf("herd [%d hosts] $ ", len(l.Runner.Hosts))
}

func (l *InteractiveLoop) AutoComplete() readline.AutoCompleter {
	p := readline.PcItem
	return readline.NewPrefixCompleter(
		p("set",
			p("Timeout"),
			p("HostTimeout"),
			p("ConnectTimeout"),
			p("Parallel"),
		),
		p("add hosts"),
		p("remove hosts"),
		p("list hosts",
			p("oneline"),
		),
		p("run"),
	)
}
