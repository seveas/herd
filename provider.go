package katyusha

type HostProvider interface {
	GetHosts(hostnameGlob string, attributes HostAttributes) Hosts
}

type Providers []HostProvider

func LoadProviders(c AppConfig) Providers {
	ret := make(Providers, 0)

	// Always load the known hosts provider
	khp := NewKnownHostsProvider()
	ret = append(ret, khp)

	// Load the consul provider if ...?

	// Load the puppetdb provider if ...?

	return ret
}

func (p *Providers) GetHosts(hostnameGlob string, attributes HostAttributes) Hosts {
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
