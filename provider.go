package katyusha

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

type HostProvider interface {
	GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts
}

type Providers []HostProvider

func LoadProviders() Providers {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	home, err := homedir.Dir()
	if err == nil {
		files = append(files, path.Join(home, ".ssh", "known_hosts"))
	}
	ret := Providers{
		&KnownHostsProvider{
			Files: files,
		},
		&CliProvider{},
	}

	// Load the consul provider if ...?

	// Load the puppetdb provider if ...?

	return ret
}

func (p *Providers) GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts {
	ret := make(Hosts, 0)
	seen := make(map[string]int)
	for _, provider := range *p {
		for _, host := range provider.GetHosts(hostnameGlob, attributes) {
			if existing, ok := seen[host.Name]; ok {
				ret[existing].Amend(host)
			} else {
				seen[host.Name] = len(ret)
				ret = append(ret, host)
			}
		}
	}
	return ret
}
