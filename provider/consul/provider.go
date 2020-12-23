package consul

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"github.com/seveas/herd"

	consul "github.com/hashicorp/consul/api"
	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("consul", newConsulProvider, consulProviderMagic)
}

type consulProvider struct {
	name         string
	consulConfig *consul.Config
	config       struct {
		Address string
		Prefix  string
		Timeout time.Duration
	}
}

func newConsulProvider(name string) herd.HostProvider {
	p := &consulProvider{name: name}
	p.config.Timeout = 10 * time.Second
	p.consulConfig = consul.DefaultConfig()
	return p
}

func consulProviderMagic(r *herd.Registry) {
	addr, _ := os.LookupEnv("CONSUL_HTTP_ADDR")
	if addr == "" {
		_, err := net.LookupHost("consul.service.consul")
		if err == nil {
			addr = "http://consul.service.consul:8500"
		}
	}
	if addr != "" {
		p := newConsulProvider("consul").(*consulProvider)
		p.config.Address = addr
		r.AddMagicProvider(herd.NewCacheFromProvider(p))
	}
}

func (p *consulProvider) Name() string {
	return p.name
}

func (p *consulProvider) Prefix() string {
	return p.config.Prefix
}

func (p *consulProvider) Equivalent(o herd.HostProvider) bool {
	return p.config.Address == o.(*consulProvider).config.Address
}

func (p *consulProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *consulProvider) Load(ctx context.Context, mc chan herd.CacheMessage) (herd.Hosts, error) {
	p.consulConfig.Address = p.config.Address
	client, err := consul.NewClient(p.consulConfig)
	if err != nil {
		return herd.Hosts{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()
	catalog := client.Catalog()
	datacenters, err := catalog.Datacenters()
	if err != nil {
		return herd.Hosts{}, err
	}
	logrus.Debugf("Consul datacenters: %v", datacenters)
	sg := scattergather.New(int64(len(datacenters)))
	for _, dc := range datacenters {
		sg.Run(func(ctx context.Context, args ...interface{}) (interface{}, error) {
			dc := args[0].(string)
			name := fmt.Sprintf("%s@%s", p.name, dc)
			mc <- herd.CacheMessage{Name: name, Finished: false, Err: nil}
			hosts, err := p.loadDatacenter(dc)
			mc <- herd.CacheMessage{Name: name, Finished: true, Err: err}
			return hosts, err
		}, ctx, dc)
	}

	untypedResults, err := sg.Wait()
	if err != nil {
		return herd.Hosts{}, err
	}

	hosts := make(herd.Hosts, 0)
	for _, hu := range untypedResults {
		hosts = append(hosts, hu.(herd.Hosts)...)
	}
	return hosts, nil
}

func (p *consulProvider) loadDatacenter(dc string) (herd.Hosts, error) {
	nodePositions := make(map[string]int)
	client, err := consul.NewClient(p.consulConfig)
	catalog := client.Catalog()
	if err != nil {
		return herd.Hosts{}, err
	}
	opts := consul.QueryOptions{Datacenter: dc, WaitTime: 5 * time.Second}
	catalognodes, _, err := catalog.Nodes(&opts)
	if err != nil {
		return herd.Hosts{}, err
	}
	hosts := make(herd.Hosts, len(catalognodes))
	for i, node := range catalognodes {
		nodePositions[node.Node] = i
		hosts[i] = herd.NewHost(node.Node, herd.HostAttributes{"datacenter": node.Datacenter})
	}
	services, _, err := catalog.Services(&opts)
	if err != nil {
		return hosts, err
	}
	for service, _ := range services {
		servicenodes, _, err := catalog.Service(service, "", &opts)
		if err != nil {
			return hosts, err
		}
		for _, service := range servicenodes {
			h := hosts[nodePositions[service.Node]]
			s := []string{}
			si, ok := h.Attributes["service"]
			if ok {
				s = si.([]string)
			}
			h.Attributes["service"] = append(s, service.ServiceName)
			h.Attributes[fmt.Sprintf("service:%s", service.ServiceName)] = service.ServiceTags
		}
	}

	for _, h := range hosts {
		if s, ok := h.Attributes["service"]; ok {
			ss := s.([]string)
			sort.Strings(ss)
			h.Attributes["service"] = ss
		}
	}
	return hosts, nil
}
