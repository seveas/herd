package google

import (
	"context"
	"fmt"
	"reflect"

	"github.com/seveas/herd"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func init() {
	herd.RegisterProvider("google", newProvider, nil)
}

type googleProvider struct {
	name   string
	zones  map[string]*computepb.Zone
	config struct {
		Prefix  string
		Key     string
		Project string
		Zones   []string
	}
}

func newProvider(name string) herd.HostProvider {
	return &googleProvider{
		name:  name,
		zones: make(map[string]*computepb.Zone),
	}
}

func (p *googleProvider) Name() string {
	return p.name
}

func (p *googleProvider) Prefix() string {
	return p.config.Prefix
}

func (p *googleProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*googleProvider)
	return p.config.Key == op.config.Key &&
		p.config.Project == op.config.Project &&
		reflect.DeepEqual(p.config.Zones, op.config.Zones)
}

func (p *googleProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *googleProvider) Load(ctx context.Context, lm herd.LoadingMessage) (hosts herd.Hosts, err error) {
	lm(p.name, false, nil)
	defer func() { lm(p.name, true, err) }()

	if len(p.config.Zones) == 0 {
		if err := p.setZones(); err != nil {
			return nil, err
		}
	}
	logrus.Debugf("GCP zones: %v", p.config.Zones)
	sg := scattergather.New(int64(len(p.config.Zones)))
	for _, region := range p.config.Zones {
		sg.Run(func(ctx context.Context, args ...interface{}) (interface{}, error) {
			zone := args[0].(string)
			name := fmt.Sprintf("%s@%s", p.name, zone)
			lm(name, false, nil)
			hosts, err := p.loadZone(zone)
			lm(name, true, err)
			return hosts, err
		}, ctx, region)
	}

	untypedResults, err := sg.Wait()
	if err != nil {
		return nil, err
	}

	hosts = make(herd.Hosts, 0)
	for _, hu := range untypedResults {
		hosts = append(hosts, hu.(herd.Hosts)...)
	}
	return hosts, nil
}

func (p *googleProvider) setZones() error {
	ctx := context.Background()
	client, err := compute.NewZonesRESTClient(ctx, option.WithCredentialsFile(p.config.Key))
	if err != nil {
		return err
	}
	p.config.Zones = make([]string, 0)
	req := &computepb.ListZonesRequest{Project: p.config.Project}
	it := client.List(ctx, req)
	for {
		zone, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if *zone.Status == computepb.Zone_UP {
			p.config.Zones = append(p.config.Zones, *zone.Name)
			p.zones[*zone.Name] = zone
		}
	}
	return nil
}

func bv(b *bool) bool {
	return b != nil && *b
}
func sv(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func iv(i *uint64) int {
	if i == nil {
		return 0
	}
	return int(*i)
}
func (p *googleProvider) loadZone(zone string) (herd.Hosts, error) {
	ctx := context.Background()
	region := p.zones[zone].Region
	client, err := compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile(p.config.Key))
	if err != nil {
		return nil, err
	}
	req := &computepb.ListInstancesRequest{
		Project: p.config.Project,
		Zone:    zone,
	}
	it := client.List(ctx, req)
	ret := make(herd.Hosts, 0)
	for {
		inst, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		attrs := herd.HostAttributes{
			"can_ip_forward":   bv(inst.CanIpForward),
			"cpuplatform":      sv(inst.CpuPlatform),
			"description":      sv(inst.Description),
			"id":               iv(inst.Id),
			"fingerprint":      sv(inst.Fingerprint),
			"kind":             sv(inst.Kind),
			"labelfingerprint": sv(inst.LabelFingerprint),
			"last_start":       sv(inst.LastStartTimestamp),
			"last_stop":        sv(inst.LastStopTimestamp),
			"last_suspend":     sv(inst.LastSuspendedTimestamp),
			"machinetype":      sv(inst.MachineType),
			"mincpuplatform":   sv(inst.MinCpuPlatform),
			"instancename":     sv(inst.Name),
			"startrestricted":  bv(inst.StartRestricted),
			"statusmessage":    sv(inst.StatusMessage),
			"region":           region,
			"zone":             zone,
		}
		if inst.Labels != nil {
			for k, v := range inst.Labels {
				attrs[k] = v
			}
		}
		if inst.Tags != nil {
			attrs["tags"] = inst.Tags
		}
		name := *inst.Name
		if inst.Hostname != nil {
			name = *inst.Hostname
		}
		ret = append(ret, herd.NewHost(name, *inst.NetworkInterfaces[0].NetworkIP, attrs))
	}

	return ret, nil
}
