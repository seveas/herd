package transip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/vps"
)

func init() {
	herd.RegisterProvider("transip", newProvider, nil)
}

type transipProvider struct {
	name   string
	config struct {
		Prefix string
		User   string
		Key    string
	}
}

func newProvider(name string) herd.HostProvider {
	return &transipProvider{name: name}
}

func (p *transipProvider) Name() string {
	return p.name
}

func (p *transipProvider) Prefix() string {
	return p.config.Prefix
}

func (p *transipProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*transipProvider)
	return p.config.User == op.config.User
}

func (p *transipProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

type ctxRoundTripper struct {
	ctx context.Context
	rt  http.RoundTripper
}

func (r *ctxRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.rt.RoundTrip(req.WithContext(r.ctx))
}

func (p *transipProvider) Load(ctx context.Context, lm herd.LoadingMessage) (hosts *herd.HostSet, err error) {
	lm(p.name, false, nil)
	defer func() { lm(p.name, true, err) }()

	hc := &http.Client{Transport: &ctxRoundTripper{ctx: ctx, rt: http.DefaultTransport}}
	client, err := gotransip.NewClient(gotransip.ClientConfiguration{
		AccountName:    p.config.User,
		PrivateKeyPath: p.config.Key,
		HTTPClient:     hc,
	})
	if err != nil {
		return nil, err
	}
	repo := &vps.Repository{Client: client}
	vpss, err := repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Unable to get VPS list: %s", err)
	}
	ret := herd.NewHostSet()
	for _, vps := range vpss {
		attrs := herd.HostAttributes{
			"uuid":             vps.UUID,
			"description":      vps.Description,
			"productname":      vps.ProductName,
			"operatingsystem":  vps.OperatingSystem,
			"disksize":         vps.DiskSize,
			"memorysize":       vps.MemorySize,
			"cpus":             vps.CPUs,
			"status":           vps.Status,
			"ipaddress":        vps.IPAddress,
			"macaddress":       vps.MacAddress,
			"currentsnapshots": vps.CurrentSnapshots,
			"maxsnapshots":     vps.MaxSnapshots,
			"islocked":         vps.IsLocked,
			"isblocked":        vps.IsBlocked,
			"iscustomerlocked": vps.IsCustomerLocked,
			"availabilityzone": vps.AvailabilityZone,
			"tags":             vps.Tags,
		}
		ret.AddHost(herd.NewHost(vps.Name, vps.IPAddress, attrs))
	}
	return ret, nil
}
