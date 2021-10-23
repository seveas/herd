package example

import (
	"context"
	"fmt"

	"github.com/seveas/katyusha"

	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("example", newProvider, nil)
}

type exampleProvider struct {
	name   string
	config struct {
		Prefix string
		Color  string
	}
}

func newProvider(name string) katyusha.HostProvider {
	return &exampleProvider{name: name}
}

func (p *exampleProvider) Name() string {
	return p.name
}

func (p *exampleProvider) Prefix() string {
	return p.config.Prefix
}

func (p *exampleProvider) Equivalent(o katyusha.HostProvider) bool {
	return true
}

func (p *exampleProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *exampleProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (katyusha.Hosts, error) {
	nhosts := 5
	hosts := make(katyusha.Hosts, nhosts)
	for i := 0; i < nhosts; i++ {
		attrs := katyusha.HostAttributes{
			"static_attribute":  "static_value",
			"dynamic_attribute": fmt.Sprintf("dynamic_value_%d", i),
			"config_color":      p.config.Color,
		}
		hosts[i] = katyusha.NewHost(fmt.Sprintf("host-%d.example.com", i), "", attrs)
	}
	return hosts, nil
}
