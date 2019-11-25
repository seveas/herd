package katyusha

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type ConsulProvider struct {
	Name          string
	Address       string
	File          string
	CacheLifetime time.Duration
}

func init() {
	ProviderMakers["consul"] = func(name string, v *viper.Viper) (HostProvider, error) {
		p := &ConsulProvider{
			Name:          name,
			File:          path.Join(viper.GetString("CacheDir"), name+".cache"),
			CacheLifetime: 1 * time.Hour,
		}
		err := v.Unmarshal(p)
		if err != nil {
			return nil, err
		}
		return &Cache{File: p.File, Lifetime: p.CacheLifetime, Provider: p}, nil
	}
	ProviderMagic["consul"] = func() []HostProvider {
		addr, ok := os.LookupEnv("CONSUL_HTTP_ADDR")
		if !ok {
			return []HostProvider{}
		}
		p := &ConsulProvider{
			Name:          "consul",
			File:          path.Join(viper.GetString("CacheDir"), "consul.cache"),
			Address:       addr,
			CacheLifetime: 1 * time.Hour,
		}
		return []HostProvider{&Cache{File: p.File, Lifetime: p.CacheLifetime, Provider: p}}
	}
}

type ConsulService struct {
	Name string
	Tags []string
}

type ConsulServices []ConsulService

func (s ConsulService) String() string {
	if len(s.Tags) > 0 {
		return fmt.Sprintf("%s=%s", s.Name, strings.Join(s.Tags, ","))
	}
	return s.Name
}

func (s ConsulServices) String() string {
	data := make([]string, len(s))
	for i, service := range s {
		data[i] = service.String()
	}
	return strings.Join(data, ";")
}

func (s ConsulServices) Match(m MatchAttribute) bool {
	// FIXME support regexes
	v, ok := m.Value.(string)
	var tags []string
	if !ok {
		return false
	}
	if strings.ContainsRune(v, ':') {
		parts := strings.SplitN(v, ":", 2)
		v = parts[0]
		tags = strings.Split(parts[1], ",")
	}
	for _, s := range s {
		if s.Name == v {
			for _, t := range tags {
				found := false
				for _, t2 := range s.Tags {
					if t == t2 || t == "" {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
			return true
		}
	}
	return false
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
	UI.Debugf("Consul datacenters: %v", datacenters)
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
			errs.AddHidden(r.err)
		}
		hosts = append(hosts, r.hosts...)
		todo -= 1
	}
	if len(errs.Errors) != 0 {
		return hosts, errs
	}
	return hosts, nil
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
			s := ConsulServices{}
			si, ok := hosts[nodePositions[service.Node]].Attributes["service"]
			if ok {
				s = si.(ConsulServices)
			}
			svc := ConsulService{Name: service.ServiceName, Tags: service.ServiceTags}
			hosts[nodePositions[service.Node]].Attributes["service"] = append(s, svc)
		}
	}
	return hosts, nil
}

func (p *ConsulProvider) PostProcess(h Hosts) {
	for i, host := range h {
		if services, ok := host.Attributes["service"]; ok {
			data, _ := json.Marshal(services)
			svc := ConsulServices{}
			json.Unmarshal(data, &svc)
			h[i].Attributes["service"] = svc
		}
	}
}
