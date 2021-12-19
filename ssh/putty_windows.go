package ssh

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func (c *Config) readPuttyConfig(name string) {
	if alias, ok := hostAliases[name]; ok {
		name = alias
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, fmt.Sprintf(`Software\SimonTatham\PuTTY\Sessions\%s`, name), registry.QUERY_VALUE)
	if err != nil {
		k, err = registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\Sessions\Default%20Settings`, registry.QUERY_VALUE)
	}
	if err != nil {
		return
	}
	if iv, _, err := k.GetIntegerValue("PortNumber"); err != nil {
		c.Port = int(iv)
	}
	if sv, _, err := k.GetStringValue("UserName"); err != nil {
		c.ClientConfig.User = sv
	}
}
