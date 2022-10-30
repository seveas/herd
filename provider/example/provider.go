package example

import (
	"context"
	"fmt"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("example", newProvider, nil)
}

type exampleProvider struct {
	name   string
	config struct {
		Prefix string
		Color  string
	}
}

func newProvider(name string) herd.HostProvider {
	return &exampleProvider{name: name}
}

func (p *exampleProvider) Name() string {
	return p.name
}

func (p *exampleProvider) Prefix() string {
	return p.config.Prefix
}

func (p *exampleProvider) Equivalent(o herd.HostProvider) bool {
	return true
}

func (p *exampleProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *exampleProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	nhosts := 5
	hosts := new(herd.HostSet)
	for i := 0; i < nhosts; i++ {
		attrs := herd.HostAttributes{
			"static_attribute":  "static_value",
			"dynamic_attribute": fmt.Sprintf("dynamic_value_%d", i),
			"config_color":      p.config.Color,
		}
		hosts.AddHost(herd.NewHost(fmt.Sprintf("host-%d.example.com", i), "", attrs))
	}
	return hosts, nil
}
