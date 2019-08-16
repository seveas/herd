package herd

import (
	"net"
)

type CliProvider struct {
	Name string
}

func (p *CliProvider) GetHosts(name string, attributes MatchAttributes) Hosts {
	if _, err := net.LookupHost(name); err != nil {
		return Hosts{}
	}
	return Hosts{NewHost(name, HostAttributes{})}

}
