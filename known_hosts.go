package herd

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type KnownHostsProvider struct {
	name   string
	hashed bool
	config struct {
		Prefix string
		Files  []string
	}
}

func (p *KnownHostsProvider) Name() string {
	return p.name
}

func (p *KnownHostsProvider) Prefix() string {
	return p.config.Prefix
}

func NewKnownHostsProvider(name string) HostProvider {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	if home, ok := os.LookupEnv("HOME"); ok {
		files = append(files, filepath.Join(home, ".ssh", "known_hosts"))
	} else {
		u, err := user.Current()
		if err == nil && u.HomeDir != "" {
			files = append(files, filepath.Join(u.HomeDir, ".ssh", "known_hosts"))
		}
	}
	p := &KnownHostsProvider{
		name:   name,
		hashed: false,
	}
	p.config.Files = files
	return p
}

func (p *KnownHostsProvider) Equivalent(o HostProvider) bool {
	op, ok := o.(*KnownHostsProvider)
	return ok && reflect.DeepEqual(p.config.Files, op.config.Files)
}

func (p *KnownHostsProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *KnownHostsProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
	seen := make(map[string]int)
	for _, f := range p.config.Files {
		data, err := ioutil.ReadFile(f)
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
				if !p.hashed {
					logrus.Warnf("Hashed hostnames found in %s. Please set `HashKnownHosts no` in your ssh config and delte the hashed entries", f)
					p.hashed = true
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
			host := NewHost(name, HostAttributes{"PublicKeyComment": comment})
			host.AddPublicKey(key)
			seen[host.Name] = len(hosts)
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}
