package consul

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/cache"

	consul "github.com/hashicorp/consul/api"
	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("consul", newProvider, magicProvider)
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

func newProvider(name string) herd.HostProvider {
	p := &consulProvider{name: name}
	p.config.Timeout = 10 * time.Second
	p.consulConfig = consul.DefaultConfig()
	return p
}

func magicProvider() herd.HostProvider {
	addr, _ := os.LookupEnv("CONSUL_HTTP_ADDR")
	if addr == "" {
		_, err := net.LookupHost("consul.service.consul")
		if err == nil {
			addr = "http://consul.service.consul:8500"
		}
	}
	if addr != "" {
		p := newProvider("consul").(*consulProvider)
		p.config.Address = addr
		return cache.NewFromProvider(p)
	}
	return nil
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

func (p *consulProvider) Load(ctx context.Context, lm herd.LoadingMessage) (herd.Hosts, error) {
	p.consulConfig.Address = p.config.Address
	lm(p.name, false, nil)
	client, err := consul.NewClient(p.consulConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()
	catalog := client.Catalog()
	datacenters, err := catalog.Datacenters()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Consul datacenters: %v", datacenters)
	sg := scattergather.New(int64(len(datacenters)))
	for _, dc := range datacenters {
		sg.Run(func(ctx context.Context, args ...interface{}) (interface{}, error) {
			dc := args[0].(string)
			name := fmt.Sprintf("%s@%s", p.name, dc)
			lm(name, false, nil)
			hosts, err := p.loadDatacenter(dc)
			lm(name, true, err)
			return hosts, err
		}, ctx, dc)
	}

	untypedResults, err := sg.Wait()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	opts := consul.QueryOptions{Datacenter: dc, WaitTime: 5 * time.Second}
	catalognodes, _, err := catalog.Nodes(&opts)
	if err != nil {
		return nil, err
	}
	hosts := make(herd.Hosts, len(catalognodes))
	for i, node := range catalognodes {
		nodePositions[node.Node] = i
		hosts[i] = herd.NewHost(node.Node, herd.HostAttributes{"datacenter": node.Datacenter})
	}
	services, _, err := catalog.Services(&opts)
	if err != nil {
		return nil, err
	}
	for service, _ := range services {
		servicenodes, _, err := catalog.Service(service, "", &opts)
		if err != nil {
			return nil, err
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
