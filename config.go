package herd

import (
	"os"
	"os/user"
	"path"
	"time"
)

type RunnerConfig struct {
	ConnectTimeout time.Duration
	HostTimeout    time.Duration
	Timeout        time.Duration
	Parallel       int
}

type AppConfig struct {
	List        bool
	ListOneline bool
	Interactive bool
	Formatter   Formatter
	HistoryDir  string
	Runner      RunnerConfig
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
	c.Formatter = NewPrettyFormatter()
	usr, err := user.Current()
	if err == nil && usr.HomeDir != "" {
		c.HistoryDir = path.Join(usr.HomeDir, ".herd", "history")
	} else {
		c.HistoryDir, _ = os.Getwd()
	}
	c.Runner.ConnectTimeout = 3 * time.Second
	c.Runner.HostTimeout = 10 * time.Second
	c.Runner.Timeout = 60 * time.Second
	c.Runner.Parallel = 0
}
