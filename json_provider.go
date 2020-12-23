package katyusha

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JsonProvider struct {
	name   string
	hashed bool
	config struct {
		Prefix string
		File   string
	}
}

func NewJsonProvider(name string) HostProvider {
	return &JsonProvider{name: name}
}

func (p *JsonProvider) Name() string {
	return p.name
}

func (p *JsonProvider) Prefix() string {
	return p.config.Prefix
}

func (p *JsonProvider) Equivalent(o HostProvider) bool {
	return p.config.File == o.(*JsonProvider).config.File
}

func (p *JsonProvider) SetDataDir(dir string) {
	if !filepath.IsAbs(p.config.File) {
		p.config.File = filepath.Join(dir, p.config.File)
	}
}

func (p *JsonProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *JsonProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
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
