package common

import (
	"context"

	"github.com/seveas/katyusha"

	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Logger interface {
	LoadingMessage(name string, done bool, err error)
	EmitLogMessage(level logrus.Level, message string)
}

type Provider interface {
	Configure(map[string]interface{}) error
	Load(ctx context.Context, logger Logger) (katyusha.Hosts, error)
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "KATYUSHA",
	MagicCookieValue: "plugin",
}

type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Provider
}

func (p *ProviderPlugin) GRPCServer(b *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterProviderServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: b,
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(ctx context.Context, b *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: NewProviderClient(c),
		broker: b,
		ctx:    ctx,
	}, nil
}
