package herd

import (
	"context"
	"encoding/json"
	"fmt"
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
	ProviderMagic["json"] = func() []HostProvider {
		fn := path.Join(viper.GetString("RootDir"), "inventory.json")
		if _, err := os.Stat(fn); err != nil {
			return []HostProvider{}
		}
		return []HostProvider{&JsonProvider{Name: "inventory", File: fn}}
	}
	ProviderMakers["json"] = func(name string, v *viper.Viper) (HostProvider, error) {
		return &JsonProvider{Name: name, File: viper.GetString("File")}, nil
	}
}

func (p *JsonProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
	data, err := ioutil.ReadFile(p.File)
	if err != nil {
		UI.Errorf("Could not load %s data in %s: %s", p.Name, p.File, err)
		return hosts, err
	}

	if err = json.Unmarshal(data, &hosts); err != nil {
		err = fmt.Errorf("Could not parse %s data in %s: %s", p.Name, p.File, err)
	}
	return hosts, err
}

func (p *JsonProvider) String() string {
	return p.Name
}
