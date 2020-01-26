// +build !windows

package herd

import (
	"context"
	"io"
	"github.com/spf13/viper"
)

func findPageant() io.ReadWriter {
	return nil
}

func puttyConfig(host string) map[string]string {
	return make(map[string]string)
}

type PuttyProvider struct {
	Name string
}

func (p *PuttyProvider) String() string {
	return p.Name
}

func (p *PuttyProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func (p *PuttyProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	return make(Hosts, 0), nil
}