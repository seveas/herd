package katyusha

import (
	"os"
	"os/user"
	"path"
)

type AppConfig struct {
	List        bool
	ListOneline bool
	Interactive bool
	Formatter   Formatter
	HistoryDir  string
}

func NewAppConfig() AppConfig {
	c := AppConfig{}
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
		c.HistoryDir = path.Join(usr.HomeDir, ".katyusha", "history")
	} else {
		c.HistoryDir, _ = os.Getwd()
	}
}
