package herd

import (
	"context"
	"fmt"
	"os"
	"path"
	"sort"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ConsulProvider struct {
	Name          string
	Address       string
	File          string
	CacheLifetime time.Duration
}

func init() {
	providerMakers["consul"] = func(dataDir, name string, v *viper.Viper) (HostProvider, error) {
		p := &ConsulProvider{
			Name:          name,
			File:          path.Join(dataDir, "cache", name+".cache"),
			CacheLifetime: 1 * time.Hour,
		}
		err := v.Unmarshal(p)
		if err != nil {
			return nil, err
		}
		return &Cache{File: p.File, Lifetime: p.CacheLifetime, Provider: p}, nil
	}
	providerMagic["consul"] = func(dataDir string) []HostProvider {
		addr, ok := os.LookupEnv("CONSUL_HTTP_ADDR")
		if !ok {
			return []HostProvider{}
		}
		p := &ConsulProvider{
			Name:          "consul",
			File:          path.Join(dataDir, "cache", "consul.cache"),
			Address:       addr,
			CacheLifetime: 1 * time.Hour,
		}
		return []HostProvider{&Cache{File: p.File, Lifetime: p.CacheLifetime, Provider: p}}
	}
}

func (p *ConsulProvider) String() string {
	return p.Name
}

func (p *ConsulProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	conf := consul.DefaultConfig()
	conf.Address = p.Address
	client, err := consul.NewClient(conf)
	if err != nil {
		return Hosts{}, err
	}
	catalog := client.Catalog()
	datacenters, err := catalog.Datacenters()
	if err != nil {
		return Hosts{}, err
	}
	logrus.Debugf("Consul datacenters: %v", datacenters)
	hosts := make(Hosts, 0)
	rc := make(chan loadresult)
	for _, dc := range datacenters {
		name := fmt.Sprintf("%s@%s", p.Name, dc)
		mc <- CacheMessage{name: name, finished: false, err: nil}
		go func(dc, name string) {
			hosts, err := p.LoadDatacenter(conf, dc)
			mc <- CacheMessage{name: name, finished: true, err: err}
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

func (p *ConsulProvider) LoadDatacenter(conf *consul.Config, dc string) (Hosts, error) {
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
