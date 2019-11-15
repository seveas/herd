package katyusha

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
)

type JsonProvider struct {
	Name string
	File string
}

func init() {
	ProviderMagic["json"] = func(p Providers) Providers {
		fn := path.Join(viper.GetString("RootDir"), "inventory.json")
		if _, err := os.Stat(fn); err != nil {
			return p
		}
		return append(p, &JsonProvider{Name: "inventory", File: fn})
	}
	ProviderMakers["json"] = func(name string, v *viper.Viper) (HostProvider, error) {
		return &JsonProvider{Name: name, File: viper.GetString("File")}, nil
	}
}

func (p *JsonProvider) GetHosts(hostnameGlob string) Hosts {
	hosts := make(Hosts, 0)
	data, err := ioutil.ReadFile(p.File)
	if err != nil {
		UI.Errorf("Could not load %s data in %s: %s", p.Name, p.File, err)
		return hosts
	}
	var objects []map[string]interface{}

	if err = json.Unmarshal(data, &objects); err != nil {
		UI.Errorf("Could not parse %s data in %s: %s", p.Name, p.File, err)
		return hosts
	}
	for _, obj := range objects {
		if p.PreProcess != nil {
			p.PreProcess(&obj)
		}
		hostname := ""
		var ok bool
		if val, exists := obj["name"]; exists {
			if hostname, ok = val.(string); !ok {
				UI.Debugf("Error in json: hostname should be a string, not %v", val)
				continue
			}
		}
		if val, exists := obj["hostname"]; exists {
			if hostname, ok = val.(string); !ok {
				UI.Debugf("Error in json: hostname should be a string, not %v", val)
				continue
			}
		}
		if hostname == "" {
			UI.Debugf("Error in json: object without hostname found %v", obj)
			continue
		}
		host := NewHost(hostname, obj)
		if host.Match(hostnameGlob, MatchAttributes{}) {
			hosts = append(hosts, host)
		}
	}
	return hosts
}
