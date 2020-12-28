package server

import (
	"context"

	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/provider/plugin/common"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/viper"

	_ "github.com/seveas/katyusha/provider/example"
)

type providerImpl struct {
	provider katyusha.HostProvider
}

func (p *providerImpl) Configure(values map[string]interface{}) error {
	v := viper.New()
	for k, s := range values {
		v.SetDefault(k, s)
	}
	return p.provider.ParseViper(v)
}

func (p *providerImpl) Load(ctx context.Context, logger common.Logger) (katyusha.Hosts, error) {
	return p.provider.Load(ctx, logger.LoadingMessage)
}

func ProviderPluginServer(name string) error {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	provider, err := katyusha.NewProvider(name, name)
	if err != nil {
		return err
	}

	p := &providerImpl{
		provider: provider,
	}
	pluginMap := map[string]plugin.Plugin{
		"provider": &common.ProviderPlugin{Impl: p},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.Handshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
	return nil
}

var _ common.Provider = &providerImpl{}
