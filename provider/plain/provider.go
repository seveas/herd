package plain

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("plain", newProvider, magicProvider)
}

type plainTextProvider struct {
	name   string
	config struct {
		File   string
		Prefix string
	}
}

func newProvider(name string) herd.HostProvider {
	return &plainTextProvider{name: name}
}

func magicProvider() herd.HostProvider {
	p := &plainTextProvider{name: "inventory"}
	p.config.File = "inventory"
	return p
}

func (p *plainTextProvider) Name() string {
	return p.name
}

func (p *plainTextProvider) Prefix() string {
	return p.config.Prefix
}

func (p *plainTextProvider) Equivalent(o herd.HostProvider) bool {
	return p.config.File == o.(*plainTextProvider).config.File
}

func (p *plainTextProvider) SetDataDir(dir string) error {
	if !filepath.IsAbs(p.config.File) {
		p.config.File = filepath.Join(dir, p.config.File)
		_, err := os.Stat(p.config.File)
		return err
	}
	return nil
}

func (p *plainTextProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *plainTextProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	hosts := herd.NewHostSet()
	data, err := os.ReadFile(p.config.File)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		hosts.AddHost(herd.NewHost(line, "", herd.HostAttributes{}))
	}
	return hosts, nil
}

var _ herd.DataLoader = &plainTextProvider{}
