package katyusha

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/seveas/katyusha"

	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("json", newProvider, magicProvider)
}

type jsonProvider struct {
	name   string
	hashed bool
	config struct {
		Prefix string
		File   string
	}
}

func newProvider(name string) katyusha.HostProvider {
	return &jsonProvider{name: name}
}

func magicProvider() katyusha.HostProvider {
	p := &jsonProvider{name: "inventory"}
	p.config.File = "inventory.json"
	return p
}

func (p *jsonProvider) Name() string {
	return p.name
}

func (p *jsonProvider) Prefix() string {
	return p.config.Prefix
}

func (p *jsonProvider) Equivalent(o katyusha.HostProvider) bool {
	return p.config.File == o.(*jsonProvider).config.File
}

func (p *jsonProvider) SetDataDir(dir string) error {
	if !filepath.IsAbs(p.config.File) {
		p.config.File = filepath.Join(dir, p.config.File)
		_, err := os.Stat(p.config.File)
		return err
	}
	return nil
}

func (p *jsonProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *jsonProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (katyusha.Hosts, error) {
	data, err := ioutil.ReadFile(p.config.File)
	if err != nil {
		return nil, err
	}

	var hosts katyusha.Hosts
	err = json.Unmarshal(data, &hosts)
	return nil, err
}

var _ katyusha.DataLoader = &jsonProvider{}
