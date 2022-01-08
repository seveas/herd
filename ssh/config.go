package ssh

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var localUser *user.User
var discoveredConfigs []*config
var clientVersion = "SSH-2.0-Herd-" + herd.Version()

// FIXME this needs to be in NewExecutor
func init() {
	localUser = &user.User{HomeDir: "/", Username: "nobody", Uid: "65535"}
	if u, err := user.Current(); err == nil {
		localUser = u
	}
	if home, ok := os.LookupEnv("HOME"); ok {
		localUser.HomeDir = home
	}
	fn := filepath.Join(localUser.HomeDir, ".ssh", "config")
	if configs, err := parseConfig(fn); err == nil {
		discoveredConfigs = configs
	}
	defaultConfig.clientConfig.User = localUser.Username
}

type strictHostKeyChecking int

const (
	no = iota
	yes
	ask
	acceptNew
)

type config struct {
	globs                 []string
	port                  int
	strictHostKeyChecking strictHostKeyChecking
	verifyHostKeyDns      bool
	identityFile          string
	clientConfig          *ssh.ClientConfig
}

// FIXME no duplication. Can be done by initializing with NewConfig in the initializer.
var defaultConfig = &config{
	port:                  22,
	strictHostKeyChecking: acceptNew,
	verifyHostKeyDns:      false,
	clientConfig: &ssh.ClientConfig{
		ClientVersion: clientVersion,
		Timeout:       3 * time.Second,
	},
}

func NewConfig(globs []string) *config {
	return &config{
		globs:                 globs,
		port:                  22,
		strictHostKeyChecking: acceptNew,
		verifyHostKeyDns:      false,
		clientConfig: &ssh.ClientConfig{
			ClientVersion: clientVersion,
			Timeout:       3 * time.Second,
			User:          localUser.Username,
		},
	}
}

func (c *config) readOpenSSHConfig(name string) {
	for i := len(discoveredConfigs) - 1; i >= 0; i-- {
		config := discoveredConfigs[i]
		for _, g := range config.globs {
			if ok, err := filepath.Match(g, name); ok && err == nil {
				c.updateFromConfig(config)
				break
			}
		}
	}
}

func (c *config) updateFromConfig(o *config) {
	if o.port != defaultConfig.port {
		c.port = o.port
	}
	if o.strictHostKeyChecking != defaultConfig.strictHostKeyChecking {
		c.strictHostKeyChecking = o.strictHostKeyChecking
	}
	if o.verifyHostKeyDns != defaultConfig.verifyHostKeyDns {
		c.verifyHostKeyDns = o.verifyHostKeyDns
	}
	if o.identityFile != defaultConfig.identityFile {
		c.identityFile = o.identityFile
	}
	if o.clientConfig.User != defaultConfig.clientConfig.User {
		c.clientConfig.User = o.clientConfig.User
	}
}

// Parse an openssh config file
var splitWhitespace = regexp.MustCompile("\\s")

func parseConfig(file string) ([]*config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configs := make([]*config, 0)

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
		val := strings.TrimSpace(parts[1])

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
			config.clientConfig.User = val
		case "port":
			config.port, _ = strconv.Atoi(val)
		case "verifyhostkeydns":
			config.verifyHostKeyDns = strings.ToLower(val) == "yes"
		case "stricthostkeychecking":
			switch strings.ToLower(val) {
			case "yes":
				config.strictHostKeyChecking = yes
			case "no":
				config.strictHostKeyChecking = no
			case "ask":
				// We cannot ask the user a  question and thus treat ask the same as yes
				fallthrough
			case "accept-new":
				config.strictHostKeyChecking = acceptNew
			}
		case "identityfile":
			config.identityFile = val
		}
	}
	return append(configs, config), nil
}

func (c *config) expandSshTokens(input string, hostname, username string, port int) (string, error) {
	if input == "" {
		return input, nil
	}
	if input[0] == '~' {
		input = localUser.HomeDir + input[1:]
	}
	if !strings.ContainsRune(input, '%') {
		return input, nil
	}
	var err error
	re := regexp.MustCompile("%[%CdhikLlnprTu]")
	output := re.ReplaceAllStringFunc(input, func(token string) string {
		switch token {
		case "%":
			return "%"
		case "C":
			err = errors.New("%C is not supported")
			return ""
		case "d":
			return localUser.HomeDir
		// Does not quite match openssh, but the best we can do
		case "h", "k", "n":
			var name string
			name, err = os.Hostname()
			return name
		case "i":
			return localUser.Uid
		case "L":
			var name string
			name, err = os.Hostname()
			return strings.Split(name, ".")[0]
		case "l":
			var name string
			name, err = os.Hostname()
			return name
		case "p":
			return fmt.Sprintf("%d", port)
		case "r":
			return username
		case "%T":
			return "NONE"
		case "%u":
			return localUser.Username
		}
		err = fmt.Errorf("Don't know what to return for %s", token)
		return ""
	})
	return output, err
}

// Find all variables relevant for a host, first match wins
func configForHost(host *herd.Host) *config {
	config := NewConfig(nil)
	config.readPuttyConfig(host.Name)
	config.readOpenSSHConfig(host.Name)
	if path, err := config.expandSshTokens(config.identityFile, host.Name, config.clientConfig.User, config.port); err == nil {
		config.identityFile = path
	}
	algos := []string{}
	for _, k := range host.PublicKeys() {
		algos = append(algos, k.Type())
	}
	if len(algos) != 0 {
		config.clientConfig.HostKeyAlgorithms = algos
	}
	return config
}

var hostAliases = make(map[string]string)

func RegisterHostAlias(host, alias string) {
	hostAliases[host] = alias
}
