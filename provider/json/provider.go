package katyusha

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/seveas/katyusha"

	"github.com/sirupsen/logrus"
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

func (p *jsonProvider) Load(ctx context.Context, mc chan katyusha.CacheMessage) (katyusha.Hosts, error) {
	hosts := make(katyusha.Hosts, 0)
	data, err := ioutil.ReadFile(p.config.File)
	if err != nil {
		logrus.Errorf("Could not load %s data in %s: %s", p.name, p.config.File, err)
		return hosts, err
	}

	if err = json.Unmarshal(data, &hosts); err != nil {
		logrus.Errorf("Could not parse %s data in %s: %s", p.name, p.config.File, err)
	}
	return hosts, err
}

var _ katyusha.DataLoader = &jsonProvider{}
