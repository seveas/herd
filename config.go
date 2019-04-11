package herd

type AppConfig struct {
	List        bool
	ListOneline bool
	Interactive bool
	Formatter   Formatter
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
}
