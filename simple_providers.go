package katyusha

import (
	"io/ioutil"
	"net"
	"strings"
)

type CliProvider struct {
	Name string
}

func (p *CliProvider) GetHosts(name string) Hosts {
	if _, err := net.LookupHost(name); err != nil {
		return Hosts{}
	}
	return Hosts{NewHost(name, HostAttributes{})}

}

type PlainTextProvider struct {
	Name       string
	File       string
	PreProcess func(*map[string]interface{})
}

func (p *PlainTextProvider) GetHosts(hostnameGlob string) Hosts {
	hosts := make(Hosts, 0)
	data, err := ioutil.ReadFile(p.File)
	if err != nil {
		UI.Errorf("Could not load %s data in %s: %s", p.Name, p.File, err)
		return hosts
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		host := NewHost(line, HostAttributes{})
		hosts = append(hosts, host)
	}
	return hosts
}
