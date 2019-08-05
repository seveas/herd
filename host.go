package herd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type HostAttributes map[string]interface{}

type Host struct {
	Name       string
	PublicKeys []ssh.PublicKey   `json:"-"`
	Attributes HostAttributes    `json:"-"`
	SshBanner  string            `json:"-"`
	SshConfig  *ssh.ClientConfig `json:"-"`
	Connection *ssh.Client       `json:"-"`
	LastResult *Result           `json:"-"`
}

type Hosts []*Host

func (hosts Hosts) String() string {
	var ret strings.Builder
	for i, h := range hosts {
		if i > 0 {
			ret.WriteString(", ")
		}
		ret.WriteString(h.Name)
	}
	return ret.String()
}

func (h Hosts) SortAndUniq() Hosts {
	if len(h) < 2 {
		return h
	}
	sort.Slice(h, func(i, j int) bool { return h[i].Name < h[j].Name })
	src, dst := 1, 0
	for src < len(h) {
		if h[src].Name != h[dst].Name {
			dst += 1
			if dst != src {
				h[dst] = h[src]
			}
		}
		src += 1
	}
	return h[:dst+1]
}

func NewHost(name string, pubKeys []ssh.PublicKey, attributes HostAttributes) *Host {
	h := &Host{
		Name:       name,
		PublicKeys: pubKeys,
		Attributes: attributes,
		SshConfig: &ssh.ClientConfig{
			ClientVersion: "SSH-2.0-Herd-0.1",
			Auth:          []ssh.AuthMethod{ssh.PublicKeysCallback(SshAgentKeys)},
			Timeout:       3 * time.Second,
		},
		Connection: nil,
		LastResult: nil,
	}
	h.SshConfig.HostKeyCallback = h.HostKeyCallback
	h.SshConfig.BannerCallback = h.BannerCallback
	u, err := user.Current()
	if err == nil {
		h.SshConfig.User = u.Username
	}
	parts := strings.SplitN(name, ".", 2)
	h.Attributes["hostname"] = parts[0]
	if len(parts) == 2 {
		h.Attributes["domainname"] = parts[1]
	}
	return h
}

func (host Host) String() string {
	return fmt.Sprintf("Host{Name: %s, Keys: %d, Attributes: %s, Config: %v}", host.Name, len(host.PublicKeys), host.Attributes, host.SshConfig)
}

func (h *Host) Address() string {
	return fmt.Sprintf("%s:22", h.Name)
}

var _regexpType = reflect.TypeOf(regexp.MustCompile(""))
var _stringType = reflect.TypeOf("")

func (h *Host) Match(hostnameGlob string, attributes HostAttributes) bool {
	hostMatched := false

	if hostnameGlob != "" {
		ok, err := filepath.Match(hostnameGlob, h.Name)
		if !ok || err != nil {
			return false
		}
		hostMatched = true
	}

	for attr, value := range attributes {
		val, ok := h.Attributes[attr]
		if !ok {
			if h.LastResult != nil {
				if attr == "stdout" {
					val = string(h.LastResult.Stdout)
				} else if attr == "stderr" {
					val = string(h.LastResult.Stderr)
				} else if attr == "exitstatus" {
					val = h.LastResult.ExitStatus
				} else {
					return false
				}
			} else {
				return false
			}
		}
		t1, t2 := reflect.TypeOf(value), reflect.TypeOf(val)
		if t1 != t2 && !(t1 == _regexpType && t2 == _stringType) {
			return false
		}
		if t1 == _regexpType {
			if !value.(*regexp.Regexp).MatchString(val.(string)) {
				return false
			}
		} else {
			if val != value {
				return false
			}
		}
	}

	return hostMatched
}

func (h *Host) Amend(h2 *Host) {
	for attr, value := range h2.Attributes {
		h.Attributes[attr] = value
	}
	// FIXME merge keys and ssh config
}

func (h *Host) HostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	if len(h.PublicKeys) == 0 {
		UI.Warnf("Warning: no known host key for %s, accepting any key", h.Name)
		return nil
	}
	bkey := key.Marshal()
	for _, pkey := range h.PublicKeys {
		if bytes.Equal(bkey, pkey.Marshal()) {
			return nil
		}
	}
	return fmt.Errorf("ssh: no matching host key")
}

func (h *Host) BannerCallback(message string) error {
	h.SshBanner = message
	return nil
}

type TimeoutError struct {
	message string
}

func (e TimeoutError) Error() string {
	return e.message
}

func (host *Host) Connect(ctx context.Context) (*ssh.Client, error) {
	if host.Connection != nil {
		return host.Connection, nil
	}
	UI.Debugf("Connecting to %s", host.Address())
	ctx, cancel := context.WithTimeout(ctx, host.SshConfig.Timeout)
	defer cancel()
	var client *ssh.Client
	ec := make(chan error)
	go func() {
		var err error
		client, err = ssh.Dial("tcp", host.Address(), host.SshConfig)
		ec <- err
	}()
	select {
	case <-ctx.Done():
		host.Connection = nil
		return nil, TimeoutError{"Timed out while connecting to server"}
	case err := <-ec:
		host.Connection = client
		return client, err
	}
}

func (host *Host) Disconnect() {
	if host.Connection != nil {
		UI.Debugf("Disconnecting from %s", host.Address())
		host.Connection.Close()
		host.Connection = nil
	}
}

func (host *Host) Run(ctx context.Context, command string, c chan Result) {
	r := Result{Host: host.Name, StartTime: time.Now(), ExitStatus: -1}

	if err := ctx.Err(); err != nil {
		r.Err = err
		c <- r
		return
	}
	client, err := host.Connect(ctx)
	if err != nil {
		r.Err = err
		c <- r
		return
	}
	sess, err := client.NewSession()
	if err != nil {
		r.Err = err
		c <- r
		return
	}
	defer sess.Close()

	var stdout, stderr ByteWriter
	if viper.GetString("Output") == "line" {
		prefix := fmt.Sprintf("%-*s  ", ctx.Value("hostnamelen").(int), host.Name)
		stdout = NewLineWriterBuffer(prefix, false)
		stderr = NewLineWriterBuffer(prefix, true)
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
		// FIXME: no error is ever returned, but the signal is not sent to the process either.
		// https://github.com/openssh/openssh-portable/commit/cd98925c6405e972dc9f211afc7e75e838abe81c
		// - OpenSSH 7.9 or newer required
		sess.Signal(ssh.SIGKILL)
		r.Err = TimeoutError{"Timed out while executing command"}
	case err := <-ec:
		r.Err = err
	}

	r.EndTime = time.Now()
	if r.Err != nil {
		if err, ok := r.Err.(*ssh.ExitError); ok {
			r.ExitStatus = err.ExitStatus()
		}
	} else {
		r.ExitStatus = 0
	}
	r.Stdout = stdout.Bytes()
	r.Stderr = stderr.Bytes()
	host.LastResult = &r
	c <- r
}
