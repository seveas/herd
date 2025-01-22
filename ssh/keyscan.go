package ssh

import (
	"context"
	"errors"
	"net"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type KeyScanExecutor struct {
	keyTypes       []string
	user           user.User
	connectTimeout time.Duration
}

func NewKeyScanExecutor(keyTypes []string, user user.User) (herd.Executor, error) {
	return &KeyScanExecutor{
		keyTypes: keyTypes,
		user:     user,
	}, nil
}

func (e *KeyScanExecutor) SetConnectTimeout(t time.Duration) {
	e.connectTimeout = t
}

func (e *KeyScanExecutor) Run(ctx context.Context, host *herd.Host, cmd string, oc chan herd.OutputLine) *herd.Result {
	hostKeyReceived := errors.New("host key received")
	address := host.Address
	if address == "" {
		address = host.Name
	}
	config, err := configForHost(host, &e.user)
	if err != nil {
		return &herd.Result{Host: host.Name, Err: err}
	}
	config.strictHostKeyChecking = no
	cc := config.clientConfig
	cc.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		host.AddPublicKey(key)
		return hostKeyReceived
	}
	address = net.JoinHostPort(address, strconv.Itoa(config.port))
	logrus.Debugf("Scanning keys on %s (%s)", host.Name, address)
	ctx, cancel := context.WithTimeout(ctx, cc.Timeout)
	defer cancel()
	go func() {
		for _, keyType := range e.keyTypes {
			found := false
			for _, k := range host.PublicKeys() {
				if k.Type() == keyType {
					found = true
					break
				}
			}
			if !found {
				logrus.Debugf("Don't have an %s key for %s, checking whether the host has it", keyType, host.Name)
				cc.HostKeyAlgorithms = strings.Split(keyType, ",")
				_, err := ssh.Dial("tcp", address, cc)
				if err != nil && !strings.HasSuffix(err.Error(), "host key received") {
					logrus.Debugf("Error checking %s key on %s: %s", keyType, host.Name, err)
				}
			}
		}
		cancel()
	}()
	<-ctx.Done()
	return &herd.Result{Host: host.Name}
}

func (e *KeyScanExecutor) Disconnect() {
}

var _ herd.Executor = &KeyScanExecutor{}
