package herd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"

	"github.com/seveas/herd/sshagent"
)

var localUser string
var extConfig *sshConfig

// Parse SSH configuration during startup, so the host initializer can access
// it always.
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
	extConfig, _ = parseSshConfig(fn)
}

// Hosts can have attributes of any types, but querying is limited to strings,
// booleans, numbers, nil and slices of these values.
type HostAttributes map[string]interface{}

func (h HostAttributes) prefix(prefix string) HostAttributes {
	ret := make(map[string]interface{})
	for k, v := range h {
		ret[prefix+k] = v
	}
	return ret
}

// A host represents a remote host. It can be instantiated manually, but is
// usually fetched from one or more Providers, which can all contribute to the
// hosts attributes.
type Host struct {
	Name       string
	Address    string
	Port       int
	Attributes HostAttributes
	publicKeys []ssh.PublicKey
	sshBanner  string
	sshConfig  *ssh.ClientConfig
	extConfig  map[string]string
	connection *ssh.Client
	lastResult *Result
	csum       uint32
	sshAgent   *sshagent.Agent
}

type host Host

func (h *Host) UnmarshalJSON(data []byte) error {
	var h2 host
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(&h2); err != nil {
		return err
	}
	for k, v := range h2.Attributes {
		if n, ok := v.(json.Number); ok {
			if i, err := n.Int64(); err == nil {
				h2.Attributes[k] = i
			} else {
				h2.Attributes[k], _ = n.Float64()
			}
		}
	}
	*h = Host(h2)
	h.init()
	return nil
}

// Hosts should be initialized with this function, which also initializes any
// internal data, without which SSH connections will not be possible.
func NewHost(name, address string, attributes HostAttributes) *Host {
	h := &Host{
		Name:       name,
		Address:    address,
		Port:       22,
		Attributes: attributes,
	}
	h.init()
	return h
}

func (h *Host) sshKeys() ([]ssh.Signer, error) {
	path, _ := h.extConfig["identityfile"]
	path, err := h.expandSshTokens(path)
	if err != nil {
		logrus.Errorf("Could not parse identify file path %s: %s", path, err)
		return []ssh.Signer{}, err
	}
	return h.sshAgent.SignersForPath(path), nil
}

// Set all the defults and initialize ssh configuration for the host
func (h *Host) init() {
	h.extConfig = extConfig.configForHost(h.Name)
	h.publicKeys = make([]ssh.PublicKey, 0)
	h.sshConfig = &ssh.ClientConfig{
		ClientVersion:   "SSH-2.0-Herd-0.1",
		User:            localUser,
		Timeout:         3 * time.Second,
		HostKeyCallback: h.hostKeyCallback,
		BannerCallback:  h.bannerCallback,
	}
	h.sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeysCallback(h.sshKeys)}
	h.csum = crc32.ChecksumIEEE([]byte(h.Name))
	parts := strings.SplitN(h.Name, ".", 2)
	if h.Attributes == nil {
		h.Attributes = make(HostAttributes)
	}
	h.Attributes["hostname"] = parts[0]
	if len(parts) == 2 {
		h.Attributes["domainname"] = parts[1]
	} else {
		h.Attributes["domainname"] = ""
	}
	if h.Port == 0 {
		h.Port = 22
	}
	if user, ok := h.extConfig["user"]; ok {
		h.sshConfig.User = user
	}
	if port, ok := h.extConfig["port"]; ok {
		if p, err := strconv.Atoi(port); err == nil {
			h.Port = p
		}
	}
}

func (host Host) String() string {
	return fmt.Sprintf("Host{Name: %s, Keys: %d, Attributes: %s, Config: %v}", host.Name, len(host.publicKeys), host.Attributes, host.sshConfig)
}

func (host Host) expandSshTokens(input string) (string, error) {
	if !strings.ContainsRune(input, '%') {
		return input, nil
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	home, ok := os.LookupEnv("HOME")
	if !ok {
		home = u.HomeDir
	}
	if input[0] == '~' {
		input = home + input[1:]
	}
	re := regexp.MustCompile("%[%CdhikLlnprTu]")
	output := re.ReplaceAllStringFunc(input, func(token string) string {
		switch token {
		case "%":
			return "%"
		case "C":
			panic("OOPS")
		case "d":
			return home
		// Does not quite match openssh, but the best we can do
		case "h", "k", "n":
			return host.Name
		case "i":
			return u.Uid
		case "L":
			var name string
			name, err = os.Hostname()
			return strings.Split(name, ".")[0]
		case "l":
			var name string
			name, err = os.Hostname()
			return name
		case "p":
			return fmt.Sprintf("%d", host.Port)
		case "r":
			return host.sshConfig.User
		case "%T":
			return "NONE"
		case "%u":
			return localUser
		}
		err = fmt.Errorf("Don't know what to return for %s", token)
		return ""
	})
	return output, err
}

// Adds a public key to a host. Used by the ssh know hosts provider, but can be
// used by any other code as well.
func (h *Host) AddPublicKey(k ssh.PublicKey) {
	h.publicKeys = append(h.publicKeys, k)
	algos := []string{}
	for _, k := range h.publicKeys {
		algos = append(algos, k.Type())
	}
	h.sshConfig.HostKeyAlgorithms = algos
}

func (h *Host) PublicKeys() []ssh.PublicKey {
	return h.publicKeys
}

func (h *Host) address() string {
	if h.Address == "" {
		return net.JoinHostPort(h.Name, strconv.Itoa(h.Port))
	}
	return net.JoinHostPort(h.Address, strconv.Itoa(h.Port))
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
		value, ok := h.GetAttribute(name)
		if !ok && !attribute.Negate {
			return false
		}
		if ok && !attribute.Match(value) {
			return false
		}
	}
	return true
}

func (h *Host) GetAttribute(key string) (interface{}, bool) {
	value, ok := h.Attributes[key]
	if ok {
		return value, ok
	}
	r := h.lastResult
	if r == nil {
		r = &Result{ExitStatus: -1}
	}
	switch key {
	case "name":
		return h.Name, true
	case "random":
		return h.csum, true
	case "address":
		return h.Address, true
	case "stdout":
		return string(r.Stdout), true
	case "stderr":
		return string(r.Stderr), true
	case "exitstatus":
		return r.ExitStatus, true
	case "err":
		return r.Err, true
	}
	return nil, false
}

func (h *Host) Amend(h2 *Host) {
	if h.Address == "" {
		h.Address = h2.Address
	}
	h.Attributes["herd_provider"] = append(h.Attributes["herd_provider"].([]string), h2.Attributes["herd_provider"].([]string)[0])
	for attr, value := range h2.Attributes {
		if attr == "herd_provider" {
			continue
		}
		h.Attributes[attr] = value
	}
	for _, k := range h2.publicKeys {
		h.AddPublicKey(k)
	}
}

func (h *Host) hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	// Do we have the key?
	bkey := key.Marshal()
	for _, pkey := range h.publicKeys {
		if bytes.Equal(bkey, pkey.Marshal()) {
			return nil
		}
	}

	// We don't have the key, but is it in DNS?
	if x, ok := h.extConfig["verifyhostkeydns"]; ok && strings.ToLower(x) == "yes" {
		if dnsVerify(h.Name, key) {
			h.AddPublicKey(key)
			return nil
		}
	}

	// We couldn't verify the key, what should we do?
	check, ok := h.extConfig["stricthostkeychecking"]
	if !ok || check == "" {
		// We default to accept-new instead of ask, as we cannot ask the user a
		// question and thus treat ask the same as yes
		check = "accept-new"
	}

	switch strings.ToLower(check) {
	case "accept-new":
		logrus.Warnf("ssh: no known host key for %s, accepting new key", h.Name)
		fallthrough
	case "no":
		h.AddPublicKey(key)
		return nil
	default:
		return fmt.Errorf("ssh: no host key found for %s", h.Name)
	}
}

func (h *Host) bannerCallback(message string) error {
	h.sshBanner = message
	return nil
}

type TimeoutError struct {
	message string
}

func (e TimeoutError) Error() string {
	return e.message
}

func (host *Host) keyScan(ctx context.Context, keyTypes []string) *Result {
	logrus.Debugf("Scanning keys on %s (%s)", host.Name, host.address())
	host.extConfig["stricthostkeychecking"] = "no"
	conf := *host.sshConfig
	ctx, cancel := context.WithTimeout(ctx, host.sshConfig.Timeout)
	defer cancel()
	go func() {
		for _, keyType := range keyTypes {
			found := false
			for _, k := range host.publicKeys {
				if k.Type() == keyType {
					found = true
					break
				}
			}
			if !found {
				logrus.Debugf("Don't have a %s key for %s, checking whether the host has it", keyType, host.Name)
				conf.HostKeyAlgorithms = []string{keyType}
				client, err := ssh.Dial("tcp", host.address(), &conf)
				if err != nil {
					logrus.Debugf("Error checking %s key on %s: %s", keyType, host.Name, err)
				} else {
					client.Close()
				}
			}
		}
		cancel()
	}()
	<-ctx.Done()
	return &Result{Host: host}
}

func (host *Host) connect(ctx context.Context) (*ssh.Client, error) {
	if host.connection != nil {
		return host.connection, nil
	}
	logrus.Debugf("Connecting to %s (%s)", host.Name, host.address())
	ctx, cancel := context.WithTimeout(ctx, host.sshConfig.Timeout)
	defer cancel()
	var client *ssh.Client
	ec := make(chan error)
	go func() {
		var err error
		client, err = ssh.Dial("tcp", host.address(), host.sshConfig)
		ec <- err
	}()
	select {
	case <-ctx.Done():
		host.connection = nil
		return nil, TimeoutError{"Timed out while connecting to server"}
	case err := <-ec:
		host.connection = client
		return client, err
	}
}

func (host *Host) disconnect() {
	if host.connection != nil {
		logrus.Debugf("Disconnecting from %s", host.address())
		host.connection.Close()
		host.connection = nil
	}
}

type byteWriter interface {
	Write([]byte) (int, error)
	Bytes() []byte
}

type lineWriterBuffer struct {
	oc      chan OutputLine
	host    *Host
	stderr  bool
	buf     *bytes.Buffer
	lineBuf []byte
}

func newLineWriterBuffer(host *Host, stderr bool, oc chan OutputLine) *lineWriterBuffer {
	return &lineWriterBuffer{
		buf:     bytes.NewBuffer([]byte{}),
		lineBuf: []byte{},
		host:    host,
		oc:      oc,
		stderr:  stderr,
	}
}

func (buf *lineWriterBuffer) Write(p []byte) (int, error) {
	n, err := buf.buf.Write(p)
	buf.lineBuf = bytes.Join([][]byte{buf.lineBuf, p}, []byte{})
	for {
		idx := bytes.Index(buf.lineBuf, []byte("\n"))
		if idx == -1 {
			break
		}
		buf.oc <- OutputLine{Host: buf.host, Data: buf.lineBuf[:idx+1], Stderr: buf.stderr}
		buf.lineBuf = buf.lineBuf[idx+1:]
	}
	return n, err
}

func (buf *lineWriterBuffer) Bytes() []byte {
	return buf.buf.Bytes()
}

func (host *Host) Run(ctx context.Context, command string, oc chan OutputLine) *Result {
	if strings.HasPrefix(command, "herd:keyscan:") {
		parts := strings.Split(command, ":")
		keyTypes := strings.Split(parts[2], ",")
		return host.keyScan(ctx, keyTypes)
	}
	now := time.Now()
	r := &Result{Host: host, StartTime: now, EndTime: now, ElapsedTime: 0, ExitStatus: -1}
	host.lastResult = r
	var stdout, stderr byteWriter
	if oc != nil {
		stdout = newLineWriterBuffer(host, false, oc)
		stderr = newLineWriterBuffer(host, true, oc)
	} else {
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	}

	if err := ctx.Err(); err != nil {
		r.Err = err
		return r
	}
	client, err := host.connect(ctx)
	if err != nil {
		r.EndTime = time.Now()
		r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()
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
	return r
}

func (h1 *Host) less(h2 *Host, attributes []string) bool {
	for _, attr := range attributes {
		v1, ok1 := h1.GetAttribute(attr)
		v2, ok2 := h2.GetAttribute(attr)
		// Sort nodes that are missing the attribute last
		if ok1 && !ok2 {
			return true
		}
		if !ok1 && ok2 {
			return false
		}
		if !ok1 && !ok2 {
			continue
		}
		// Compare the string values, this way we don't need to check a matrix of types
		s1, err1 := cast.ToStringE(v1)
		s2, err2 := cast.ToStringE(v2)
		if err1 != nil || err2 != nil {
			continue
		}
		// When equal, continue to the next field
		if s1 == s2 {
			continue
		}
		return s1 < s2
	}
	return h1.Name < h2.Name
}
