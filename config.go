package katyusha

import (
	"os"
	"os/user"
	"path"
	"time"
)

type LogLevel uint

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
)

type RunnerConfig struct {
	ConnectTimeout time.Duration
	HostTimeout    time.Duration
	Timeout        time.Duration
	Parallel       int
}

type UIConfig struct {
	Formatter Formatter
	LogLevel  LogLevel
}

type AppConfig struct {
	List        bool
	ListOneline bool
	ScriptFile  string
	Interactive bool
	HistoryDir  string
	Runner      RunnerConfig
	UI          UIConfig
}

func NewAppConfig() AppConfig {
	c := AppConfig{Runner: RunnerConfig{}}
	c.SetDefaults()
	return c
}

func (c *AppConfig) SetDefaults() {
	c.List = false
	c.ListOneline = false
	c.Interactive = false
	usr, err := user.Current()
	if err == nil && usr.HomeDir != "" {
		c.HistoryDir = path.Join(usr.HomeDir, ".katyusha", "history")
	} else {
		c.HistoryDir, _ = os.Getwd()
	}
	c.Runner.ConnectTimeout = 3 * time.Second
	c.Runner.HostTimeout = 10 * time.Second
	c.Runner.Timeout = 60 * time.Second
	c.Runner.Parallel = 0
	c.UI.Formatter = NewPrettyFormatter()
	c.UI.LogLevel = DEBUG
}
