package ci_dataloader

import (
	"context"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("ci_dataloader", newProvider, nil)
}

type ciProvider struct {
	name   string
	config struct {
		Prefix  string
		DataDir string
	}
}

func newProvider(name string) herd.HostProvider {
	return &ciProvider{name: name}
}

func (p *ciProvider) Name() string {
	return p.name
}

func (p *ciProvider) Prefix() string {
	return p.config.Prefix
}

func (p *ciProvider) Equivalent(o herd.HostProvider) bool {
	return true
}

func (p *ciProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *ciProvider) SetDataDir(dir string) error {
	p.config.DataDir = dir
	return nil
}

func (p *ciProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	lm(p.name, false, nil)
	hosts := herd.NewHostSet()
	hosts.AddHost(herd.NewHost("ci", "", herd.HostAttributes{"datadir": p.config.DataDir}))
	return hosts, nil
}

var _ herd.DataLoader = &ciProvider{}
