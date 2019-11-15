package katyusha

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	ProviderMagic["cli"] = func(p Providers) Providers {
		return append(p, &CliProvider{})
	}
	ProviderMagic["plain"] = func(p Providers) Providers {
		fn := path.Join(viper.GetString("RootDir"), "inventory")
		if _, err := os.Stat(fn); err != nil {
			return p
		}
		return append(p, &PlainTextProvider{Name: "inventory", File: fn})
	}
	ProviderMakers["plain"] = func(name string, v *viper.Viper) (HostProvider, error) {
		return &PlainTextProvider{Name: name, File: viper.GetString("File")}, nil
	}
}

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
	Name string
	File string
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
