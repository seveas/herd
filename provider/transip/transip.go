package transip

import (
	"context"
	"fmt"

	"github.com/seveas/katyusha"

	"github.com/spf13/viper"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/vps"
)

func init() {
	katyusha.RegisterProvider("transip", newProvider, nil)
}

type transipProvider struct {
	name   string
	config struct {
		Prefix string
		User   string
		Key    string
	}
}

func newProvider(name string) katyusha.HostProvider {
	return &transipProvider{}
}

func (p *transipProvider) Name() string {
	return p.name
}

func (p *transipProvider) Prefix() string {
	return p.config.Prefix
}

func (p *transipProvider) Equivalent(o katyusha.HostProvider) bool {
	op := o.(*transipProvider)
	return p.config.User == op.config.User
}

func (p *transipProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *transipProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (hosts katyusha.Hosts, err error) {
	lm(p.name, false, nil)
	defer func() { lm(p.name, true, err) }()

	client, err := gotransip.NewClient(gotransip.ClientConfiguration{
		AccountName:    p.config.User,
		PrivateKeyPath: p.config.Key,
	})
	if err != nil {
		return nil, err
	}
	repo := &vps.Repository{Client: client}
	vpss, err := repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Unable to get VPS list: %s", err)
	}
	ret := make(katyusha.Hosts, len(vpss))
	for i, vps := range vpss {
		attrs := katyusha.HostAttributes{
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
		ret[i] = katyusha.NewHost(vps.Name, vps.IPAddress, attrs)
	}
	return ret, nil
}
