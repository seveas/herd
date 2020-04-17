package katyusha

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type KnownHostsProvider struct {
	Name  string
	Files []string
}

func NewKnownHostsProvider(name string) HostProvider {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	home, err := homedir.Dir()
	if err == nil {
		files = append(files, filepath.Join(home, ".ssh", "known_hosts"))
	}
	return &KnownHostsProvider{Name: name, Files: files}
}

func (p *KnownHostsProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(p)
}

func (p *KnownHostsProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
	seen := make(map[string]int)
	for _, f := range p.Files {
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

func (p *KnownHostsProvider) String() string {
	return "ssh_known_hosts"
}
