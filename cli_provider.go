package herd

import (
	"golang.org/x/crypto/ssh"
	"net"
)

type CliProvider struct {
}

func NewCliProvider() *CliProvider {
	return &CliProvider{}
}

func (p *CliProvider) GetHosts(name string, attributes HostAttributes) Hosts {
	if _, err := net.LookupHost(name); err != nil {
		return Hosts{}
	}
	return Hosts{NewHost(name, []ssh.PublicKey{}, HostAttributes{})}

}
