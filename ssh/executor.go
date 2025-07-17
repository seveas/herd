package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Executor struct {
	agent          *agentPool
	knownHosts     ssh.HostKeyCallback
	user           user.User
	connectTimeout time.Duration
	disconnect     bool
}

func NewExecutor(agentCount int, agentTimeout time.Duration, user user.User, disconnect bool) (herd.Executor, error) {
	agent, err := newAgentPool(agentCount, agentTimeout)
	if err != nil {
		return nil, err
	}

	files := []string{}
	if _, err := os.Stat("/etc/ssh/ssh_known_hosts"); err == nil {
		files = append(files, "/etc/ssh/ssh_known_hosts")
	}
	if home, ok := os.LookupEnv("HOME"); ok {
		f := filepath.Join(home, ".ssh", "known_hosts")
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}

	knownHosts, err := knownhosts.New(files...)
	if err != nil {
		return nil, err
	}

	return &Executor{
		agent:      agent,
		user:       user,
		knownHosts: knownHosts,
		disconnect: disconnect,
	}, nil
}

func (e *Executor) SetConnectTimeout(t time.Duration) {
	e.connectTimeout = t
}

func (e *Executor) Run(ctx context.Context, host *herd.Host, command string, oc chan herd.OutputLine) *herd.Result {
	now := time.Now()
	r := &herd.Result{Host: host.Name, StartTime: now, EndTime: now, ElapsedTime: 0, ExitStatus: -1}
	defer func() {
		r.EndTime = time.Now()
		r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()
	}()

	if err := ctx.Err(); err != nil {
		r.Err = err
		return r
	}
	connection, err := e.connect(ctx, host)
	if err != nil {
		r.Err = err
		return r
	}
	sess, err := connection.NewSession()
	if err != nil {
		r.Err = err
		return r
	}
	defer sess.Close()

	var stdout, stderr byteWriter
	if oc != nil {
		stdout = newLineWriterBuffer(host, false, oc)
		stderr = newLineWriterBuffer(host, true, oc)
	} else {
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	}

	sess.Stdout = stdout
	sess.Stderr = stderr
	ec := make(chan error)

	go func() {
		ec <- sess.Run(command)
	}()

	select {
	case <-ctx.Done():
		terr := herd.TimeoutError{Message: "Timed out while executing command"}
		if err := sess.Signal(ssh.SIGKILL); err != nil {
			terr.Message = fmt.Sprintf("%s, and killing the session failed: %s", terr.Message, err)
		}
		r.Err = terr
	case err := <-ec:
		r.Err = err
	}
	if r.Err != nil {
		if err, ok := r.Err.(*ssh.ExitError); ok {
			r.ExitStatus = err.ExitStatus()
		}
	} else {
		r.ExitStatus = 0
	}
	if e.disconnect {
		_ = connection.Close()
		host.Connection = nil
	}
	r.Stdout = stdout.Bytes()
	r.Stderr = stderr.Bytes()
	return r
}

func (e *Executor) connect(ctx context.Context, host *herd.Host) (*ssh.Client, error) {
	if host.Connection != nil {
		return host.Connection.(*ssh.Client), nil
	}
	config, err := configForHost(host, &e.user)
	if err != nil {
		return nil, err
	}
	cc := config.clientConfig
	cc.Timeout = e.connectTimeout
	cc.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return e.hostKeyCallback(host, config.port, remote, key, config)
	}
	if config.identityFile != "" {
		cc.Auth = []ssh.AuthMethod{ssh.PublicKeysCallback(e.agent.SignersForPathCallback(config.identityFile))}
	} else {
		cc.Auth = []ssh.AuthMethod{ssh.PublicKeysCallback(e.agent.Signers)}
	}
	cc.Auth = append(cc.Auth, ssh.KeyboardInteractive(e.emptyPasswordCallback))
	address := host.Address
	if address == "" {
		address = host.Name
	}
	address = net.JoinHostPort(address, strconv.Itoa(config.port))
	logrus.Debugf("Connecting to %s (%s) as %s with key %s", host.Name, address, cc.User, config.identityFile)

	ctx, cancel := context.WithTimeout(ctx, e.connectTimeout+time.Second/2)
	defer cancel()
	var client *ssh.Client
	ec := make(chan error)
	go func() {
		var err error
		client, err = ssh.Dial("tcp", address, cc)
		ec <- err
	}()
	select {
	case <-ctx.Done():
		return nil, herd.TimeoutError{Message: "Timed out while connecting to server"}
	case err := <-ec:
		if err == nil {
			host.Connection = client
		}
		return client, err
	}
}

func (e *Executor) hostKeyCallback(host *herd.Host, port int, remote net.Addr, key ssh.PublicKey, c *config) error {
	// Do we have the key?
	bkey := key.Marshal()
	for _, pkey := range host.PublicKeys() {
		if bytes.Equal(bkey, pkey.Marshal()) {
			return nil
		}
	}

	// We don't have the key, but is it in DNS?
	if c.verifyHostKeyDns && verifyHostKeyDns(host.Name, key) {
		host.AddPublicKey(key)
		return nil
	}

	// We don't have a key, but is it in known_hosts?
	err := e.knownHosts(net.JoinHostPort(host.Name, strconv.Itoa(port)), remote, key)
	if err == nil {
		host.AddPublicKey(key)
		return nil
	}

	if errors.Is(err, &knownhosts.RevokedError{}) {
		return fmt.Errorf("ssh: host key for %s revoked", host.Name)
	}

	var ke *knownhosts.KeyError
	if errors.As(err, &ke) {
		for _, k := range ke.Want {
			if k.Key.Type() == key.Type() {
				logrus.Errorf("ssh: host key mismatch for %s, expected %v, got %v", host.Name, ke.Want[0].Key.Type(), key.Type())
				return fmt.Errorf("ssh: host key mismatch for %s", host.Name)
			}
		}
	}

	switch c.strictHostKeyChecking {
	case acceptNew:
		logrus.Warnf("ssh: no known host key for %s, accepting new key", host.Name)
		fallthrough
	case no:
		host.AddPublicKey(key)
		return nil
	default:
		return fmt.Errorf("ssh: no host key found for %s", host.Name)
	}
}

func (e *Executor) emptyPasswordCallback(name, instruction string, questions []string, echos []bool) (answers []string, err error) {
	// All we support is an empty challenge, which does not require a response
	// but can be added by some 2fa stacks if the 2fa part is bypassed
	if name == "" && instruction == "" && len(questions) == 0 {
		return []string{}, nil
	}
	return make([]string, len(questions)), fmt.Errorf("keyboard-interactive authentication not supported")
}

var _ herd.Executor = &Executor{}
