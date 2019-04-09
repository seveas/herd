package herd

import (
	"io"
	"io/ioutil"
	"log"
	"os/user"
	"path"

	"golang.org/x/crypto/ssh"
)

type KnownHostsProvider struct {
	Files []string
	User  user.User
}

func NewKnownHostsProvider() *KnownHostsProvider {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	usr, err := user.Current()
	if err == nil && usr.HomeDir != "" {
		files = append(files, path.Join(usr.HomeDir, ".ssh", "known_hosts"))
	}
	return &KnownHostsProvider{
		Files: files,
	}
}

func (p *KnownHostsProvider) GetHosts(hostnameGlob string, attributes HostAttributes) Hosts {
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
				log.Fatalf("Error parsing known hosts file %s: %s", f, err)
				// FIXME: Can we still parse the rest?
				break
			}
			data = rest
			name := matches[0]
			// FIXME: Save the rest as aliases?
			if idx, ok := seen[name]; ok {
				hosts[idx].PublicKeys = append(hosts[idx].PublicKeys, key)
				continue
			}
			host := NewHost(name, []ssh.PublicKey{key}, HostAttributes{"PublicKeyComment": comment})
			if !host.Match(hostnameGlob, attributes) {
				continue
			}
			seen[host.Name] = len(hosts)
			hosts = append(hosts, host)
		}
	}
	return hosts
}
