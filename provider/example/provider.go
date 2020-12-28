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

type plainTextProvider struct {
	name   string
	config struct {
		Prefix string
		Color  string
	}
}

func newProvider(name string) herd.HostProvider {
	return &plainTextProvider{name: name}
}

func (p *plainTextProvider) Name() string {
	return p.name
}

func (p *plainTextProvider) Prefix() string {
	return p.config.Prefix
}

func (p *plainTextProvider) Equivalent(o herd.HostProvider) bool {
	return true
}

func (p *plainTextProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *plainTextProvider) Load(ctx context.Context, lm herd.LoadingMessage) (herd.Hosts, error) {
	nhosts := 5
	hosts := make(herd.Hosts, nhosts)
	for i := 0; i < nhosts; i++ {
		attrs := herd.HostAttributes{
			"static_attribute":  "static_value",
			"dynamic_attribute": fmt.Sprintf("dynamic_value_%d", i),
			"config_color":      p.config.Color,
		}
		hosts[i] = herd.NewHost(fmt.Sprintf("host-%d.example.com", i), attrs)
	}
	return hosts, nil
}
