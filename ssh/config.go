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

var clientVersion = "SSH-2.0-Herd-" + herd.Version()

type strictHostKeyChecking int

const (
	no = iota
	yes
	ask
	acceptNew
)

type config struct {
	user   user.User
	blocks []*configBlock
}

type configBlock struct {
	globs                 []string
	port                  int
	strictHostKeyChecking strictHostKeyChecking
	verifyHostKeyDns      bool
	identityFile          string
	clientConfig          *ssh.ClientConfig
}

func newConfig(u user.User) *config {
	c := &config{
		user: u,
	}
	return c
}

func (c *config) readOpenSSHConfig() error {
	fn := filepath.Join(c.user.HomeDir, ".ssh", "config")
	blocks, err := parseConfig(fn)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	c.blocks = blocks
	blocks, err = parseConfig("/etc/ssh/ssh_config")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	c.blocks = append(c.blocks, blocks...)
	return nil
}

func newConfigBlock(globs []string) *configBlock {
	return &configBlock{
		globs:                 globs,
		port:                  22,
		strictHostKeyChecking: acceptNew,
		verifyHostKeyDns:      false,
		clientConfig: &ssh.ClientConfig{
			ClientVersion: clientVersion,
			Timeout:       3 * time.Second,
		},
	}
}

func (c *configBlock) updateFromBlock(o *configBlock) {
	if c.port == 22 {
		c.port = o.port
	}
	if c.strictHostKeyChecking == acceptNew {
		c.strictHostKeyChecking = o.strictHostKeyChecking
	}
	if !c.verifyHostKeyDns {
		c.verifyHostKeyDns = o.verifyHostKeyDns
	}
	if c.identityFile == "" {
		c.identityFile = o.identityFile
	}
	if c.clientConfig.User == "" {
		c.clientConfig.User = o.clientConfig.User
	}
}

// Parse an openssh config file
var splitWhitespace = regexp.MustCompile(`\s`)

func parseConfig(file string) ([]*configBlock, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configs := make([]*configBlock, 0)

	// Anything before the first Host section is global
	config := newConfigBlock([]string{"*"})
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
			config = newConfigBlock(splitWhitespace.Split(val, -1))
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

func (c *config) expandSshTokens(input, hostname string, b *configBlock) (string, error) {
	if input == "" {
		return input, nil
	}
	if input[0] == '~' {
		input = c.user.HomeDir + input[1:]
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
			return c.user.HomeDir
		// Does not quite match openssh, but the best we can do
		case "h", "k", "n":
			return hostname
		case "i":
			return c.user.Uid
		case "L":
			var name string
			name, err = os.Hostname()
			return strings.Split(name, ".")[0]
		case "l":
			var name string
			name, err = os.Hostname()
			return name
		case "p":
			return fmt.Sprintf("%d", b.port)
		case "r":
			return b.clientConfig.User
		case "%T":
			return "NONE"
		case "%u":
			return c.user.Username
		}
		err = fmt.Errorf("Don't know what to return for %s", token)
		return ""
	})
	return output, err
}

// Find all variables relevant for a host, first match wins
func (c *config) forHost(host *herd.Host) *configBlock {
	b := newConfigBlock(nil)
	b.readPuttyConfig(host.Name)
	for i := len(c.blocks) - 1; i >= 0; i-- {
		config := c.blocks[i]
		for _, g := range config.globs {
			if ok, err := filepath.Match(g, host.Name); ok && err == nil {
				b.updateFromBlock(config)
				break
			}
		}
	}

	if b.clientConfig.User == "" {
		b.clientConfig.User = c.user.Username
	}
	if path, err := c.expandSshTokens(b.identityFile, host.Name, b); err == nil {
		b.identityFile = path
	}
	algos := []string{}
	for _, k := range host.PublicKeys() {
		algos = append(algos, k.Type())
	}
	if len(algos) != 0 {
		b.clientConfig.HostKeyAlgorithms = algos
	}
	return b
}

var hostAliases = make(map[string]string)

func RegisterHostAlias(host, alias string) {
	hostAliases[host] = alias
}
