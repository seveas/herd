package ci_cache

import (
	"context"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("ci_cache", newProvider, nil)
}

type ciProvider struct {
	name   string
	config struct {
		Prefix     string
		CacheDir   string
		Keep       bool
		Invalidate bool
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

func (p *ciProvider) SetCacheDir(dir string) {
	p.config.CacheDir = dir
}

func (p *ciProvider) Keep() {
	p.config.Keep = true
}

func (p *ciProvider) Invalidate() {
	p.config.Invalidate = true
}

func (p *ciProvider) Source() herd.HostProvider {
	return p
}

func (p *ciProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	lm(p.name, false, nil)
	hosts := herd.NewHostSet()
	hosts.AddHost(herd.NewHost("ci", "", herd.HostAttributes{"cachedir": p.config.CacheDir, "keep": p.config.Keep, "invalidate": p.config.Invalidate}))
	return hosts, nil
}

var _ herd.Cache = &ciProvider{}
