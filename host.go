package herd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
)

// HostAttributes represent attributes on a host, they can be of any type, but querying is limited to strings,
// booleans, numbers, nil and slices of these values.
type HostAttributes map[string]any

func (h HostAttributes) prefix(prefix string) HostAttributes {
	ret := make(map[string]any)
	for k, v := range h {
		ret[prefix+k] = v
	}
	return ret
}

// Host represents a remote host. It can be instantiated manually, but is
// usually fetched from one or more Providers, which can all contribute to the
// hosts attributes.
type Host struct {
	Name       string
	Address    string
	Attributes HostAttributes
	Connection io.Closer `yaml:"-" json:"-"`
	LastResult *Result   `yaml:",omitempty" json:",omitempty"`
	publicKeys []ssh.PublicKey
	csum       uint32
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

func (h *Host) MarshalJSON() ([]byte, error) {
	if len(h.publicKeys) > 0 {
		keys := make([][]byte, len(h.publicKeys))
		for i, key := range h.publicKeys {
			keys[i] = key.Marshal()
		}
		h.Attributes["__publicKeys"] = keys
	}
	data, err := json.Marshal(host(*h))
	delete(h.Attributes, "__publicKeys")
	return data, err
}

// NewHost initializes a host object, which also initializes any internal data,
// without which SSH connections will not be possible.
func NewHost(name, address string, attributes HostAttributes) *Host {
	h := &Host{
		Name:       name,
		Address:    address,
		Attributes: attributes,
	}
	h.init()
	return h
}

// Set all the defults and initialize ssh configuration for the host
func (h *Host) init() {
	h.csum = crc32.ChecksumIEEE([]byte(h.Name))
	h.publicKeys = make([]ssh.PublicKey, 0)

	if h.Attributes == nil {
		h.Attributes = make(HostAttributes)
	}
	parts := strings.SplitN(h.Name, ".", 2)
	h.Attributes["hostname"] = parts[0]
	if len(parts) == 2 {
		h.Attributes["domainname"] = parts[1]
	} else {
		h.Attributes["domainname"] = ""
	}
	if keys, ok := h.Attributes["__publicKeys"]; ok {
		for _, k := range keys.([]any) {
			b, err := base64.StdEncoding.DecodeString(k.(string))
			if err != nil {
				logrus.Errorf("Unable to decode marshaledpublic key for %s: %s", h.Name, err)
				continue
			}
			key, err := ssh.ParsePublicKey(b)
			if err != nil {
				logrus.Errorf("Unable to parse public key for %s: %s", h.Name, err)
			}
			h.AddPublicKey(key)
		}
		delete(h.Attributes, "__publicKeys")
	}
}

func (h Host) String() string {
	return fmt.Sprintf("Host{Name: %s, Keys: %d, Attributes: %s}", h.Name, len(h.publicKeys), h.Attributes)
}

// AddPublicKey adds a public key to a host, it will be used by the SSH client
// to verify the host's identity.
func (h *Host) AddPublicKey(k ssh.PublicKey) {
	h.publicKeys = append(h.publicKeys, k)
}

func (h *Host) PublicKeys(keyTypes ...string) []ssh.PublicKey {
	if len(keyTypes) == 0 {
		return h.publicKeys
	}
	keys := make([]ssh.PublicKey, 0, len(h.publicKeys))
	for _, k := range h.publicKeys {
		if slices.Contains(keyTypes, k.Type()) {
			keys = append(keys, k)
		}
	}
	return keys
}

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
	r := h.LastResult
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
	if h2.Attributes["herd_provider"] != nil {
		if h.Attributes["herd_provider"] == nil {
			h.Attributes["herd_provider"] = make([]string, 0)
		}
		h.Attributes["herd_provider"] = append(h.Attributes["herd_provider"].([]string), h2.Attributes["herd_provider"].([]string)[0])
	}
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

func (h *Host) less(h2 *Host, attributes []string) bool {
	for _, attr := range attributes {
		v1, ok1 := h.GetAttribute(attr)
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
	return h.Name < h2.Name
}
