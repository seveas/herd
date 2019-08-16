package katyusha

import (
	"io"
	"io/ioutil"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

type KnownHostsProvider struct {
	Files []string
}

func NewKnownHostsProvider() *KnownHostsProvider {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	home, err := homedir.Dir()
	if err == nil {
		files = append(files, path.Join(home, ".ssh", "known_hosts"))
	}
	return &KnownHostsProvider{
		Files: files,
	}
}

func (p *KnownHostsProvider) GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts {
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
				UI.Warnf("Error parsing known hosts file %s: %s", f, err)
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
			if !host.Match(hostnameGlob, attributes) {
				seen[host.Name] = -1
				continue
			}
			seen[host.Name] = len(hosts)
			hosts = append(hosts, host)
		}
	}
	return hosts
}
