package tailscale

import (
	"context"
	"strings"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
	"tailscale.com/client/local"
)

func init() {
	herd.RegisterProvider("tailscale", newProvider, magicProvider)
}

func magicProvider() herd.HostProvider {
	if status, err := local.Status(context.Background()); err == nil && status.BackendState == "Running" {
		return &tailscaleProvider{name: "tailscale"}
	}
	return nil
}

type tailscaleProvider struct {
	name   string
	config struct {
		Prefix string
		Domain string
	}
}

func newProvider(name string) herd.HostProvider {
	return &tailscaleProvider{name: name}
}

func (p *tailscaleProvider) Name() string {
	return p.name
}

func (p *tailscaleProvider) Prefix() string {
	return p.config.Prefix
}

func (p *tailscaleProvider) Equivalent(o herd.HostProvider) bool {
	return true
}

func (p *tailscaleProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *tailscaleProvider) Load(ctx context.Context, lm herd.LoadingMessage) (hosts *herd.HostSet, err error) {
	lm(p.name, false, nil)
	defer func() { lm(p.name, true, err) }()

	status, err := local.Status(ctx)
	if err != nil {
		return nil, err
	}

	ret := herd.NewHostSet()
	for _, peer := range status.Peer {
		ips := make([]string, len(peer.TailscaleIPs))
		for i, ip := range peer.TailscaleIPs {
			ips[i] = ip.String()
		}
		attrs := herd.HostAttributes{
			"id":            peer.ID,
			"publickey":     peer.PublicKey,
			"hostname":      peer.HostName,
			"dnsname":       peer.DNSName,
			"os":            peer.OS,
			"userid":        peer.UserID,
			"tailscale_ips": ips,
			"addrs":         peer.Addrs,
			"curaddr":       peer.CurAddr,
			"relay":         peer.Relay,
			"rxbytes":       peer.RxBytes,
			"txbytes":       peer.TxBytes,
			"created":       peer.Created,
			"lastwrite":     peer.LastWrite,
			"lastseen":      peer.LastSeen,
			"lasthandhake":  peer.LastHandshake,
			"exitnode":      peer.ExitNode,
			"active":        peer.Active,
			"apiurl":        peer.PeerAPIURL,
			"shareenode":    peer.ShareeNode,
			"innetworkmap":  peer.InNetworkMap,
			"inmagicsock":   peer.InMagicSock,
			"inengine":      peer.InEngine,
		}
		name := peer.DNSName
		if name == "" {
			name = peer.TailscaleIPs[0].String()
		} else if p.config.Domain != "" {
			name = peer.DNSName[:strings.IndexRune(peer.DNSName, '.')] + "." + p.config.Domain
		}
		ret.AddHost(herd.NewHost(name, peer.TailscaleIPs[0].String(), attrs))
	}
	return ret, nil
}
