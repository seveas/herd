package herd

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PlainTextProvider struct {
	name   string
	config struct {
		File   string
		Prefix string
	}
}

func NewPlainTextProvider(name string) HostProvider {
	return &PlainTextProvider{name: name}
}

func (p *PlainTextProvider) Name() string {
	return p.name
}

func (p *PlainTextProvider) Prefix() string {
	return p.config.Prefix
}

func (p *PlainTextProvider) Equivalent(o HostProvider) bool {
	return p.config.File == o.(*PlainTextProvider).config.File
}

func (p *PlainTextProvider) SetDataDir(dir string) {
	if !filepath.IsAbs(p.config.File) {
		p.config.File = filepath.Join(dir, p.config.File)
	}
}
func (p *PlainTextProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *PlainTextProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
	data, err := ioutil.ReadFile(p.config.File)
	if err != nil {
		logrus.Errorf("Could not load %s data in %s: %s", p.name, p.config.File, err)
		return hosts, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		host := NewHost(line, HostAttributes{})
		hosts = append(hosts, host)
	}
	return hosts, nil
}

var _ DataLoader = &PlainTextProvider{}
