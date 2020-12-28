package server

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/provider/plugin/common"

	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
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
	logrus.SetOutput(ioutil.Discard)
	logrus.AddHook(&logrusHook{logger: logger})
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

type logrusHook struct {
	logger common.Logger
}

func (h *logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logrusHook) Fire(entry *logrus.Entry) error {
	h.logger.EmitLogMessage(entry.Level, entry.Message)
	return nil
}

var _ common.Provider = &providerImpl{}
