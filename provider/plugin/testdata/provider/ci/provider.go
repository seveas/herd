package ci

import (
	"context"
	"errors"
	"fmt"

	"github.com/seveas/katyusha"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("ci", newProvider, nil)
}

type ciProvider struct {
	name   string
	config struct {
		Prefix string
		Mode   string
	}
}

func newProvider(name string) katyusha.HostProvider {
	return &ciProvider{name: name}
}

func (p *ciProvider) Name() string {
	return p.name
}

func (p *ciProvider) Prefix() string {
	return p.config.Prefix
}

func (p *ciProvider) Equivalent(o katyusha.HostProvider) bool {
	return true
}

func (p *ciProvider) ParseViper(v *viper.Viper) error {
	err := v.Unmarshal(&p.config)
	if err != nil {
		return err
	}
	switch p.config.Mode {
	case "config-error":
		return errors.New("Simulated configuration error")
	case "config-panic":
		panic("Simulated configuration panic")
	}
	return nil
}

func (p *ciProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (katyusha.Hosts, error) {
	lm(p.name, false, nil)
	for _, level := range logrus.AllLevels {
		if level > logrus.FatalLevel {
			logrus.StandardLogger().Logf(level, "%s message", level)
		}
	}
	switch p.config.Mode {
	case "normal":
		nhosts := 5
		hosts := make(katyusha.Hosts, nhosts)
		for i := 0; i < nhosts; i++ {
			attrs := katyusha.HostAttributes{
				"static_attribute":  "static_value",
				"dynamic_attribute": fmt.Sprintf("dynamic_value_%d", i),
			}
			hosts[i] = katyusha.NewHost(fmt.Sprintf("host-%d.example.com", i), attrs)
		}
		return hosts, nil
	case "empty":
		return nil, nil
	case "error":
		return nil, errors.New("Simulated load error")
	case "panic":
		panic("Simulated load panic")
	}
	return nil, fmt.Errorf("Unknown provider mode: %s", p.config.Mode)
}
