package katyusha

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"net"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type HostAttributes map[string]interface{}

type MatchAttribute struct {
	Name        string
	FuzzyTyping bool
	Negate      bool
	Regex       bool
	Value       interface{}
}

type MatchValue interface {
	Match(m MatchAttribute) bool
}

func (m MatchAttribute) String() string {
	c1, c2 := '=', '='
	if m.Negate {
		c1 = '!'
	}
	if m.Regex {
		c2 = '~'
	}
	return fmt.Sprintf("%s %c%c %s", m.Name, c1, c2, m.Value)
}

func (m MatchAttribute) Match(value interface{}) (matches bool) {
	defer func() {
		if m.Negate {
			matches = !matches
		}
	}()
	if m.Value == value {
		return true
	}
	if m.Regex {
		svalue, ok := value.(string)
		return ok && m.Value.(*regexp.Regexp).MatchString(svalue)
	}
	if v, ok := value.(MatchValue); ok {
		return v.Match(m)
	}
	if m.FuzzyTyping {
		if bvalue, ok := value.(bool); ok && (m.Value == "true" || m.Value == "false") {
			return bvalue == (m.Value == "true")
		}
		if m.Value == "nil" {
			return value == nil
		}
		myival, err := strconv.ParseInt(m.Value.(string), 0, 64)
		if err != nil {
			return false
		}
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int() == myival
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(v.Uint()) == myival
		}
	}
	// Let's be gentle on all the int types in attributes
	if myival, ok := m.Value.(int64); ok {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int() == myival
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(v.Uint()) == myival
		}
	}
	return false
}

type MatchAttributes []MatchAttribute

type Host struct {
	Name       string
	Attributes HostAttributes
	PublicKeys []ssh.PublicKey   `json:"-"`
	SshBanner  string            `json:"-"`
	SshConfig  *ssh.ClientConfig `json:"-"`
	Connection *ssh.Client       `json:"-"`
	LastResult *Result           `json:"-"`
	Csum       uint32            `json:"-"`
}

func NewHost(name string, attributes HostAttributes) *Host {
	h := &Host{
		Name:       name,
		PublicKeys: make([]ssh.PublicKey, 0),
		Attributes: attributes,
		SshConfig: &ssh.ClientConfig{
			ClientVersion: "SSH-2.0-Katyusha-0.1",
			Auth:          []ssh.AuthMethod{ssh.PublicKeysCallback(SshAgentKeys)},
			Timeout:       3 * time.Second,
		},
		Connection: nil,
		LastResult: nil,
		Csum:       crc32.ChecksumIEEE([]byte(name)),
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
	} else {
		h.Attributes["domainname"] = ""
	}
	return h
}

func (host Host) String() string {
	return fmt.Sprintf("Host{Name: %s, Keys: %d, Attributes: %s, Config: %v}", host.Name, len(host.PublicKeys), host.Attributes, host.SshConfig)
}

func (h *Host) AddPublicKey(k ssh.PublicKey) {
	h.PublicKeys = append(h.PublicKeys, k)
}

func (h *Host) Address() string {
	return fmt.Sprintf("%s:22", h.Name)
}

var _regexpType = reflect.TypeOf(regexp.MustCompile(""))
var _stringType = reflect.TypeOf("")

func (h *Host) Match(hostnameGlob string, attributes MatchAttributes) bool {

	if hostnameGlob != "" {
		ok, err := filepath.Match(hostnameGlob, h.Name)
		if !ok || err != nil {
			return false
		}
	}

	for _, attribute := range attributes {
		name := attribute.Name
		value, ok := h.Attributes[name]
		if !ok {
			if h.LastResult != nil {
				if name == "stdout" {
					ok = true
					value = string(h.LastResult.Stdout)
				} else if name == "stderr" {
					ok = true
					value = string(h.LastResult.Stderr)
				} else if name == "exitstatus" {
					ok = true
					value = h.LastResult.ExitStatus
				} else if name == "err" {
					ok = true
					value = h.LastResult.Err
				}
			}
		}
		if !ok && !attribute.Negate {
			return false
		}
		if ok && !attribute.Match(value) {
			return false
		}
	}
	return true
}

func (h *Host) Amend(h2 *Host) {
	for attr, value := range h2.Attributes {
		h.Attributes[attr] = value
	}
	for _, k := range h2.PublicKeys {
		h.AddPublicKey(k)
	}
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

func (host *Host) Run(ctx context.Context, command string) Result {
	now := time.Now()
	r := Result{Host: host.Name, StartTime: now, EndTime: now, ElapsedTime: 0, ExitStatus: -1}
	var stdout, stderr ByteWriter
	if viper.GetString("Output") == "line" {
		prefix := fmt.Sprintf("%-*s  ", ctx.Value("hostnamelen").(int), host.Name)
		stdout = UI.NewLineWriterBuffer(host, prefix, false)
		stderr = UI.NewLineWriterBuffer(host, prefix, true)
		defer func() {
			if r.Err != nil {
				stderr.(*LineWriterBuffer).WriteStatus(r)
			} else {
				stdout.(*LineWriterBuffer).WriteStatus(r)
			}
		}()
	} else {
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	}

	if err := ctx.Err(); err != nil {
		r.Err = err
		return r
	}
	client, err := host.Connect(ctx)
	if err != nil {
		r.Err = err
		return r
	}
	sess, err := client.NewSession()
	if err != nil {
		r.Err = err
		return r
	}
	defer sess.Close()

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
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()
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
	return r
}
