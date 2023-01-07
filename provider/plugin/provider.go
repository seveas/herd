package plugin

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/plugin/common"

	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("plugin", newPlugin, nil)
}

type pluginProvider struct {
	name   string
	plugin common.ProviderPluginImpl
	logger *logForwarder
	config struct {
		Command  string
		Prefix   string
		Checksum string
		checksum []byte
	}
}

func newPlugin(name string) herd.HostProvider {
	p := &pluginProvider{name: name, logger: &logForwarder{}}
	if path, err := exec.LookPath(fmt.Sprintf("herd-provider-%s", name)); err == nil {
		p.config.Command = path
	}
	if path, err := exec.LookPath(fmt.Sprintf("herd-provider-%s.exe", name)); err == nil {
		p.config.Command = path
	}
	return p
}

func (p *pluginProvider) Name() string {
	return p.name
}

func (p *pluginProvider) Prefix() string {
	return p.config.Prefix
}

func (p *pluginProvider) ParseViper(v *viper.Viper) error {
	if err := v.Unmarshal(&p.config); err != nil {
		return err
	}
	if p.config.Checksum == "" {
		h := crypto.SHA256.New()
		if fd, err := os.Open(p.config.Command); err == nil {
			if _, err := io.Copy(h, fd); err == nil {
				p.config.checksum = h.Sum(nil)
			}
		}
		logrus.Debugf("Checksum for %s: %s", p.config.Command, hex.EncodeToString(p.config.checksum))
	} else {
		cs, err := hex.DecodeString(p.config.Checksum)
		if err != nil {
			return err
		}
		p.config.checksum = cs
	}
	if err := p.connect(); err != nil {
		return err
	}
	return p.plugin.Configure(v.AllSettings())
}

func (p *pluginProvider) Equivalent(o herd.HostProvider) bool {
	return p.config.Command == o.(*pluginProvider).config.Command
}

func (p *pluginProvider) SetDataDir(dir string) error {
	if p.plugin == nil {
		return errors.New("SetDataDir called before plugin was connected")
	}
	return p.plugin.SetDataDir(dir)
}

func (p *pluginProvider) SetCacheDir(dir string) {
	if p.plugin == nil {
		logrus.Errorf("Invalidate called before plugin was connected")
		return
	}
	p.plugin.SetCacheDir(dir)
}

func (p *pluginProvider) Keep() {
	if p.plugin == nil {
		logrus.Errorf("Invalidate called before plugin was connected")
		return
	}
	p.plugin.Keep()
}

func (p *pluginProvider) Invalidate() {
	if p.plugin == nil {
		logrus.Errorf("Invalidate called before plugin was connected")
		return
	}
	p.plugin.Invalidate()
}

func (p *pluginProvider) Source() herd.HostProvider {
	return p
}

func (p *pluginProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	if p.plugin == nil {
		return nil, errors.New("SetDataDir called before plugin was connected")
	}
	p.logger.lm = lm
	return p.plugin.Load(ctx)
}

func (p *pluginProvider) connect() error {
	pluginMap := map[string]plugin.Plugin{
		"provider": &common.ProviderPlugin{},
	}
	client := plugin.NewClient(&plugin.ClientConfig{
		Managed:          true,
		HandshakeConfig:  common.Handshake,
		Plugins:          pluginMap,
		Cmd:              exec.Command(p.config.Command),
		Logger:           common.NewLogrusLogger(logrus.StandardLogger(), fmt.Sprintf("plugin-%s", p.name)),
		SyncStdout:       os.Stdout,
		SyncStderr:       os.Stderr,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		SecureConfig:     &plugin.SecureConfig{Hash: crypto.SHA256.New(), Checksum: p.config.checksum},
	}) //#nosec G204 -- Cmd is user-supplied by design

	rpcClient, err := client.Client()
	if err != nil {
		return err
	}
	raw, err := rpcClient.Dispense("provider")
	if err != nil {
		return err
	}
	p.plugin = raw.(common.ProviderPluginImpl)
	if err := p.plugin.SetLogger(p.logger); err != nil {
		return err
	}
	return nil
}

type logForwarder struct {
	lm herd.LoadingMessage
}

func (l *logForwarder) LoadingMessage(name string, done bool, err error) {
	l.lm(name, done, err)
}

func (l *logForwarder) EmitLogMessage(level logrus.Level, message string) {
	logrus.StandardLogger().Log(level, message)
}

// Static checks to make sure we implement the interfaces we want to implement
var (
	_ common.Logger   = &logForwarder{}
	_ herd.Cache      = &pluginProvider{}
	_ herd.DataLoader = &pluginProvider{}
)
