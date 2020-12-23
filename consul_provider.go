// +build !no_consul

package herd

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	availableProviders["consul"] = NewConsulProvider
	magicProviders["consul"] = func(r *Registry) {
		addr, _ := os.LookupEnv("CONSUL_HTTP_ADDR")
		if addr == "" {
			_, err := net.LookupHost("consul.service.consul")
			if err == nil {
				addr = "http://consul.service.consul:8500"
			}
		}
		if addr != "" {
			p := NewConsulProvider("consul").(*ConsulProvider)
			p.config.Address = addr
			r.AddMagicProvider(NewCacheFromProvider(p))
		}
	}
}

type ConsulProvider struct {
	name   string
	config struct {
		Address string
		Prefix  string
		Timeout time.Duration
	}
}

func NewConsulProvider(name string) HostProvider {
	return &ConsulProvider{name: name}
}

func (p *ConsulProvider) Name() string {
	return p.name
}

func (p *ConsulProvider) Prefix() string {
	return p.config.Prefix
}

func (p *ConsulProvider) Equivalent(o HostProvider) bool {
	return p.config.Address == o.(*ConsulProvider).config.Address
}

func (p *ConsulProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *ConsulProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	conf := consul.DefaultConfig()
	conf.Address = p.config.Address
	client, err := consul.NewClient(conf)
	if err != nil {
		return Hosts{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()
	catalog := client.Catalog()
	datacenters, err := catalog.Datacenters()
	if err != nil {
		return Hosts{}, err
	}
	logrus.Debugf("Consul datacenters: %v", datacenters)
	hosts := make(Hosts, 0)
	rc := make(chan loadresult)
	for _, dc := range datacenters {
		name := fmt.Sprintf("%s@%s", p.name, dc)
		mc <- CacheMessage{Name: name, Finished: false, Err: nil}
		go func(dc, name string) {
			hosts, err := p.loadDatacenter(conf, dc)
			mc <- CacheMessage{Name: name, Finished: true, Err: err}
			rc <- loadresult{hosts: hosts, err: err}
		}(dc, name)
	}
	todo := len(datacenters)
	errs := &MultiError{}
	for todo > 0 {
		r := <-rc
		if r.err != nil {
			errs.Add(r.err)
		}
		hosts = append(hosts, r.hosts...)
		todo -= 1
	}

	if !errs.HasErrors() {
		return hosts, nil
	}
	return hosts, errs
}

func (p *ConsulProvider) loadDatacenter(conf *consul.Config, dc string) (Hosts, error) {
	nodePositions := make(map[string]int)
	client, err := consul.NewClient(conf)
	catalog := client.Catalog()
	if err != nil {
		return Hosts{}, err
	}
	opts := consul.QueryOptions{Datacenter: dc, WaitTime: 5 * time.Second}
	catalognodes, _, err := catalog.Nodes(&opts)
	if err != nil {
		return Hosts{}, err
	}
	hosts := make(Hosts, len(catalognodes))
	for i, node := range catalognodes {
		nodePositions[node.Node] = i
		hosts[i] = NewHost(node.Node, HostAttributes{"datacenter": node.Datacenter})
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
