package server

import (
	"context"
	"fmt"
	"io"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/plugin/common"

	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type pluginImpl struct {
	provider herd.HostProvider
	logger   common.Logger
}

func (p *pluginImpl) SetLogger(l common.Logger) error {
	p.logger = l
	return nil
}

func (p *pluginImpl) Configure(values map[string]interface{}) error {
	v := viper.New()
	for k, s := range values {
		v.SetDefault(k, s)
	}
	return p.provider.ParseViper(v)
}

func stripCache(p herd.HostProvider) herd.HostProvider {
	if c, ok := p.(herd.Cache); ok {
		return c.Source()
	}
	return p
}

func (p *pluginImpl) SetDataDir(dir string) error {
	if p, ok := stripCache(p.provider).(herd.DataLoader); ok {
		return p.SetDataDir(dir)
	}
	return nil
}

func (p *pluginImpl) SetCacheDir(dir string) {
	if c, ok := p.provider.(herd.Cache); ok {
		c.SetCacheDir(dir)
	}
}

func (p *pluginImpl) Invalidate() {
	if c, ok := p.provider.(herd.Cache); ok {
		c.Invalidate()
	}
}

func (p *pluginImpl) Keep() {
	if c, ok := p.provider.(herd.Cache); ok {
		c.Keep()
	}
}

func (p *pluginImpl) Load(ctx context.Context) (*herd.HostSet, error) {
	return p.provider.Load(ctx, p.logger.LoadingMessage)
}

func ProviderPluginServer(name string) error {
	provider, err := herd.NewProvider(name, name)
	if err != nil {
		return err
	}

	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)
	p := &pluginImpl{
		provider: provider,
	}
	pluginMap := map[string]plugin.Plugin{
		"provider": &common.ProviderPlugin{Impl: p},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.Handshake,
		Plugins:         pluginMap,
		Logger:          common.NewLogrusLogger(logrus.StandardLogger(), fmt.Sprintf("plugin-server-%s", name)),
		GRPCServer:      plugin.DefaultGRPCServer,
	})
	return nil
}

var _ common.ProviderPluginImpl = &pluginImpl{}
