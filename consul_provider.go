package herd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type ConsulNode struct {
	Name     string `json:"name"`
	Node     *consul.Node
	Services []*consul.CatalogService
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
		return p, nil
	}
	ProviderMagic["consul"] = func(p Providers) Providers {
		if addr, ok := os.LookupEnv("CONSUL_HTTP_ADDR"); ok {
			return append(p, &ConsulProvider{Name: "consul", File: path.Join(viper.GetString("CacheDir"), "consul.cache"), Address: addr, CacheLifetime: 1 * time.Hour})
		}
		return p
	}
}

type ConsulServices []struct {
	name string
	tags []string
}

func (s ConsulServices) String() string {
	data := make([]string, len(s))
	for i, service := range s {
		if len(service.tags) > 0 {
			data[i] = fmt.Sprintf("%s=%s", service.name, strings.Join(service.tags, ","))
		} else {
			data[i] = service.name
		}
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
		if s.name == v {
			for _, t := range tags {
				found := false
				for _, t2 := range s.tags {
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

func (p *ConsulProvider) Cache(mc chan CacheMessage, ctx context.Context) error {
	if info, err := os.Stat(p.File); err == nil && time.Since(info.ModTime()) < p.CacheLifetime {
		return nil
	}

	conf := consul.DefaultConfig()
	conf.Address = p.Address
	client, err := consul.NewClient(conf)
	if err != nil {
		return err
	}
	catalog := client.Catalog()
	datacenters, err := catalog.Datacenters()
	if err != nil {
		return err
	}
	UI.Debugf("Consul datacenters: %v", datacenters)
	nodes := make([]*ConsulNode, 0)
	nc := make(chan map[string]*ConsulNode)
	for _, dc := range datacenters {
		name := fmt.Sprintf("%s@%s", p.Name, dc)
		mc <- CacheMessage{name: name, finished: false, err: nil}
		go func(dc, name string) {
			nodes, err := p.CacheDatacenter(conf, dc)
			mc <- CacheMessage{name: name, finished: true, err: err}
			nc <- nodes
		}(dc, name)
	}
	todo := len(datacenters)
	for todo > 0 {
		for key, node := range <-nc {
			node.Name = key
			nodes = append(nodes, node)
		}
		todo -= 1
	}
	data, err := json.Marshal(nodes)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(p.File+".new", data, 0600); err != nil {
		return err
	}
	if err := os.Rename(p.File+".new", p.File); err != nil {
		return err
	}
	return nil
}

func (p *ConsulProvider) CacheDatacenter(conf *consul.Config, dc string) (map[string]*ConsulNode, error) {
	nodes := make(map[string]*ConsulNode)
	client, err := consul.NewClient(conf)
	catalog := client.Catalog()
	if err != nil {
		return nodes, err
	}
	opts := consul.QueryOptions{Datacenter: dc, WaitTime: 5 * time.Second}
	catalognodes, _, err := catalog.Nodes(&opts)
	if err != nil {
		return nodes, err
	}
	for _, node := range catalognodes {
		nodes[node.Node] = &ConsulNode{Node: node, Services: []*consul.CatalogService{}}
	}
	services, _, err := catalog.Services(&opts)
	if err != nil {
		return nodes, err
	}
	for service, _ := range services {
		servicenodes, _, err := catalog.Service(service, "", &opts)
		if err != nil {
			return nodes, err
		}
		for _, node := range servicenodes {
			nodes[node.Node].Services = append(nodes[node.Node].Services, node)
		}
	}
	return nodes, nil
}

func (p *ConsulProvider) GetHosts(hostnameGlob string) Hosts {
	jp := &JsonProvider{Name: p.Name, File: p.File, PreProcess: func(obj *map[string]interface{}) {
		node := (*obj)["Node"].(map[string]interface{})
		(*obj)["datacenter"] = node["Datacenter"]
		svc := (*obj)["Services"].([]interface{})
		services := make(ConsulServices, len(svc))
		for i, s := range svc {
			service := s.(map[string]interface{})
			services[i].name = service["ServiceName"].(string)
			svct := service["ServiceTags"].([]interface{})
			services[i].tags = make([]string, len(svct))
			for j, t := range svct {
				services[i].tags[j] = t.(string)
			}
		}
		(*obj)["service"] = services
		delete(*obj, "Services")
		delete(*obj, "Node")
	},
	}
	return jp.GetHosts(hostnameGlob)
}
