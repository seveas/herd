package ssh

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var localUser string
var discoveredConfigs []*Config
var clientVersion = "SSH-2.0-Herd-0.7.1" // FIXME

func init() {
	var fn string
	if home, ok := os.LookupEnv("HOME"); ok {
		fn = filepath.Join(home, ".ssh", "config")
	}
	if u, err := user.Current(); err == nil {
		localUser = u.Username
		if fn == "" && u.HomeDir != "" {
			fn = filepath.Join(u.HomeDir, ".ssh", "config")
		}
	}
	if configs, err := parseConfig(fn); err != nil {
		discoveredConfigs = configs
	}
	defaultConfig.ClientConfig.User = localUser
}

type StrictHostKeyChecking int

const (
	No = iota
	Yes
	Ask
	AcceptNew
)

type Config struct {
	Globs                 []string
	Port                  int
	StrictHostKeyChecking StrictHostKeyChecking
	VerifyHostKeyDns      bool
	ClientConfig          *ssh.ClientConfig
	IdentityFile          string
}

var defaultConfig = &Config{
	Port:                  22,
	StrictHostKeyChecking: AcceptNew,
	VerifyHostKeyDns:      false,
	ClientConfig: &ssh.ClientConfig{
		ClientVersion: clientVersion,
		User:          localUser,
		Timeout:       3 * time.Second,
	},
}

func NewConfig(globs []string) *Config {
	return &Config{
		Globs:                 globs,
		Port:                  22,
		StrictHostKeyChecking: AcceptNew,
		VerifyHostKeyDns:      false,
		ClientConfig: &ssh.ClientConfig{
			ClientVersion: clientVersion,
			User:          localUser,
			Timeout:       3 * time.Second,
		},
	}
}

func (c *Config) readOpenSSHConfig(name string) {
	for i := len(discoveredConfigs) - 1; i >= 0; i-- {
		config := discoveredConfigs[i]
		for _, g := range config.Globs {
			if ok, err := filepath.Match(g, name); ok && err == nil {
				c.updateFromConfig(config)
				break
			}
		}
	}
}

func (c *Config) updateFromConfig(o *Config) {
	if o.Port != defaultConfig.Port {
		c.Port = o.Port
	}
	if o.StrictHostKeyChecking != defaultConfig.StrictHostKeyChecking {
		c.StrictHostKeyChecking = o.StrictHostKeyChecking
	}
	if o.VerifyHostKeyDns != defaultConfig.VerifyHostKeyDns {
		c.VerifyHostKeyDns = o.VerifyHostKeyDns
	}
	if o.IdentityFile != defaultConfig.IdentityFile {
		c.IdentityFile = o.IdentityFile
	}
	if o.ClientConfig.User != defaultConfig.ClientConfig.User {
		c.ClientConfig.User = o.ClientConfig.User
	}
}

// Parse an openssh config file
var splitWhitespace = regexp.MustCompile("\\s")

func parseConfig(file string) ([]*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configs := make([]*Config, 0)

	// Anything before the first Host section is global
	config := NewConfig([]string{"*"})
	seen := make(map[string]bool)

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		// Ignore comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := splitWhitespace.Split(line, 2)
		if len(parts) != 2 {
			logrus.Errorf("Ignoring invalid ssh config line: %s", line)
			continue
		}
		key := strings.ToLower(parts[0])
		val := parts[1]

		// New host section, add the existing section to the returned configs
		if key == "host" {
			configs = append(configs, config)
			config = NewConfig(splitWhitespace.Split(val, -1))
			seen = make(map[string]bool)
			continue
		}

		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = true

		switch key {
		case "user":
			config.ClientConfig.User = val
		case "port":
			config.Port, _ = strconv.Atoi(val)
		case "verifyhostkeydns":
			config.VerifyHostKeyDns = strings.ToLower("val") == "yes"
		case "stricthostkeychecking":
			switch strings.ToLower(val) {
			case "yes":
				config.StrictHostKeyChecking = Yes
			case "no":
				config.StrictHostKeyChecking = No
			case "ask":
				// We cannot ask the user a  question and thus treat ask the same as yes
				fallthrough
			case "accept-new":
				config.StrictHostKeyChecking = AcceptNew
			}
		case "identityfile":
			config.IdentityFile = val

		}
	}
	return append(configs, config), nil
}

// Find all variables relevant for a host, first match wins
func ConfigForHost(name string) *Config {
	config := NewConfig(nil)
	config.readPuttyConfig(name)
	config.readOpenSSHConfig(name)

	return config
}

var hostAliases = make(map[string]string)

func RegisterHostAlias(host, alias string) {
	hostAliases[host] = alias
}
