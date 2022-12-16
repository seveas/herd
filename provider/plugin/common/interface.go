package common

import (
	"context"

	"github.com/seveas/herd"

	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Logger interface {
	LoadingMessage(name string, done bool, err error)
	EmitLogMessage(level logrus.Level, message string)
}

type ProviderPluginImpl interface {
	SetLogger(Logger) error
	Configure(map[string]interface{}) error
	Load(ctx context.Context) (*herd.HostSet, error)
	SetDataDir(string) error
	SetCacheDir(string)
	Invalidate()
	Keep()
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "HERD",
	MagicCookieValue: "plugin",
}

type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ProviderPluginImpl
}

func (p *ProviderPlugin) GRPCServer(b *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterProviderPluginServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: b,
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(ctx context.Context, b *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: NewProviderPluginClient(c),
		broker: b,
		ctx:    ctx,
	}, nil
}
