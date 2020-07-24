package herd

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JsonProvider struct {
	BaseProvider `mapstructure:",squash"`
	File         string
}

func NewJsonProvider(name string) HostProvider {
	return &JsonProvider{BaseProvider: BaseProvider{Name: name}}
}

func (p *JsonProvider) Equivalent(o HostProvider) bool {
	if c, ok := o.(*Cache); ok {
		o = c.Source
	}
	op, ok := o.(*JsonProvider)
	return ok && p.File == op.File
}

func (p *JsonProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(p)
}

func (p *JsonProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
	data, err := ioutil.ReadFile(p.File)
	if err != nil {
		logrus.Errorf("Could not load %s data in %s: %s", p.Name, p.File, err)
		return hosts, err
	}

	if err = json.Unmarshal(data, &hosts); err != nil {
		logrus.Errorf("Could not parse %s data in %s: %s", p.Name, p.File, err)
	}
	for _, h := range hosts {
		h.init()
	}
	return hosts, err
}
