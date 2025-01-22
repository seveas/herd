package ssh

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/seveas/herd"

	"github.com/kevinburke/ssh_config"
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
	port                  int
	strictHostKeyChecking strictHostKeyChecking
	verifyHostKeyDns      bool
	identityFile          string
	clientConfig          *ssh.ClientConfig
}

func configForHost(host *herd.Host, user *user.User) (*config, error) {
	c := &config{
		port:                  22,
		strictHostKeyChecking: acceptNew,
		verifyHostKeyDns:      false,
		clientConfig: &ssh.ClientConfig{
			ClientVersion: clientVersion,
			Timeout:       3 * time.Second,
			HostKeyAlgorithms: []string{
				// We make this match openssh's default host key algorithms minus the sk- variants
				ssh.CertAlgoED25519v01,
				ssh.CertAlgoECDSA256v01,
				ssh.CertAlgoECDSA384v01,
				ssh.CertAlgoECDSA521v01,
				ssh.CertAlgoRSASHA512v01,
				ssh.CertAlgoRSASHA256v01,
				ssh.KeyAlgoED25519,
				ssh.KeyAlgoECDSA256,
				ssh.KeyAlgoECDSA384,
				ssh.KeyAlgoECDSA521,
				ssh.KeyAlgoRSASHA512,
				ssh.KeyAlgoRSASHA256,
				ssh.KeyAlgoRSA,
			},
		},
	}
	readPuttyConfig(c, host.Name)
	// Use only algorithms we have keys for, if we have keys
	algos := []string{}
	for _, k := range host.PublicKeys() {
		algos = append(algos, k.Type())
	}
	if len(algos) != 0 {
		c.clientConfig.HostKeyAlgorithms = algos
	}
	// Set user from ssh_config, or use the current user
	c.clientConfig.User = ssh_config.Get(host.Name, "user")
	if c.clientConfig.User == "" {
		c.clientConfig.User = user.Username
	}
	// Set identity file, and parse its tokens
	idf := ssh_config.Get(host.Name, "identityfile")
	if idf != "" {
		if path, err := expandSshTokens(idf, host, user, c); err != nil {
			return nil, err
		} else if _, err = os.Stat(path); err == nil {
			c.identityFile = path
		}
	}
	port := ssh_config.Get(host.Name, "port")
	if port != "" {
		if porti, err := strconv.Atoi(port); err != nil {
			return nil, fmt.Errorf("Invalid port number: %s", port)
		} else {
			c.port = porti
		}
	}
	c.verifyHostKeyDns = ssh_config.Get(host.Name, "verifyhostkeydns") == "yes"
	switch strings.ToLower(ssh_config.Get(host.Name, "stricthostkeychecking")) {
	case "yes":
		c.strictHostKeyChecking = yes
	case "no":
		c.strictHostKeyChecking = no
	case "ask":
		// We cannot ask the user a  question and thus treat ask the same as yes
		fallthrough
	case "accept-new":
		c.strictHostKeyChecking = acceptNew
	}

	return c, nil
}

func expandSshTokens(input string, host *herd.Host, user *user.User, c *config) (string, error) {
	if input == "" {
		return input, nil
	}
	if input[0] == '~' {
		input = user.HomeDir + input[1:]
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
			return user.HomeDir
		// Does not quite match openssh, but the best we can do
		case "h", "k", "n":
			return host.Name
		case "i":
			return user.Uid
		case "L":
			var name string
			name, err = os.Hostname()
			return strings.Split(name, ".")[0]
		case "l":
			var name string
			name, err = os.Hostname()
			return name
		case "p":
			return fmt.Sprintf("%d", c.port)
		case "r":
			return c.clientConfig.User
		case "T":
			return "NONE"
		case "u":
			return user.Username
		}
		err = fmt.Errorf("Don't know what to return for %s", token)
		return ""
	})
	return output, err
}

var hostAliases = make(map[string]string)

func RegisterHostAlias(host, alias string) {
	hostAliases[host] = alias
}
