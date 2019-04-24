package main

import (
	"fmt"
	"io"
	"path"

	"github.com/chzyer/readline"
	"github.com/seveas/herd"
)

type InteractiveLoop struct {
	Config *herd.AppConfig
	Runner *herd.Runner
}

func (l *InteractiveLoop) Run() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          l.Prompt(),
		AutoComplete:    l.AutoComplete(),
		HistoryFile:     path.Join(l.Config.HistoryDir, "interactive"),
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
		commands, err := herd.ParseCode(line+"\n", l.Config)
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
			p("timeout"),
			p("hosttimeout"),
			p("connecttimeout"),
			p("parallel"),
		),
		p("add hosts"),
		p("remove hosts"),
		p("list hosts",
			p("oneline"),
		),
		p("run"),
	)
}
