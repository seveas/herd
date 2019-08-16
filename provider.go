package katyusha

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

type HostProvider interface {
	GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts
}

type Providers []HostProvider

func LoadProviders() (Providers, error) {
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
	providers := viper.Sub("Providers")
	if providers != nil {
		for key, _ := range providers.AllSettings() {
			ps := providers.Sub(key)
			pname := ps.GetString("Provider")
			var provider HostProvider
			if pname == "json" {
				provider = &JsonProvider{}
			}
			err := ps.Unmarshal(provider)
			if err != nil {
				return nil, err
			}
			ret = append(ret, provider)
		}
	}
	return ret, nil
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
