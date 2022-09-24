package known_hosts

import (
	"context"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func init() {
	herd.RegisterProvider("known_hosts", newProvider, magicProvider)
}

type knownHostsProvider struct {
	name   string
	config struct {
		Prefix string
		Files  []string
	}
}

func (p *knownHostsProvider) Name() string {
	return p.name
}

func (p *knownHostsProvider) Prefix() string {
	return p.config.Prefix
}

func newProvider(name string) herd.HostProvider {
	return &knownHostsProvider{name: name}
}

func magicProvider() herd.HostProvider {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	if home, ok := os.LookupEnv("HOME"); ok {
		files = append(files, filepath.Join(home, ".ssh", "known_hosts"))
	} else {
		u, err := user.Current()
		if err == nil && u.HomeDir != "" {
			files = append(files, filepath.Join(u.HomeDir, ".ssh", "known_hosts"))
		}
	}
	p := &knownHostsProvider{name: "known_hosts"}
	p.config.Files = files
	return p
}

func (p *knownHostsProvider) Equivalent(o herd.HostProvider) bool {
	return reflect.DeepEqual(p.config.Files, o.(*knownHostsProvider).config.Files)
}

func (p *knownHostsProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *knownHostsProvider) Load(ctx context.Context, lm herd.LoadingMessage) (herd.Hosts, error) {
	hosts := make(herd.Hosts, 0)
	seen := make(map[string]int)
	for _, f := range p.config.Files {
		hashed := false
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		for {
			_, matches, key, comment, rest, err := ssh.ParseKnownHosts(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				logrus.Warnf("Error parsing known hosts file %s: %s", f, err)
				data = rest
				continue
			}
			data = rest
			name := matches[0]
			if strings.HasPrefix(name, "|") {
				if !hashed {
					logrus.Warnf("Hashed hostnames found in %s. Please set `HashknownHosts no` in your ssh config and delete the hashed entries", f)
					hashed = true
				}
				continue
			}
			if idx, ok := seen[name]; ok {
				// -1 means: seen but did not match
				// FIXME: if we ever match on key attributes, this is wrong.
				if idx != -1 {
					hosts[idx].AddPublicKey(key)
				}
				continue
			}
			host := herd.NewHost(name, "", herd.HostAttributes{"PublicKeyComment": comment})
			host.AddPublicKey(key)
			seen[host.Name] = len(hosts)
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}
