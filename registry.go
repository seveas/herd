package herd

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Registry struct {
	providers []HostProvider
	hosts     Hosts
	sort      []string
	dataDir   string
}

type Hosts []*Host

type HostProvider interface {
	Load(ctx context.Context, mc chan CacheMessage) (Hosts, error)
	String() string
}

type CacheMessage struct {
	name     string
	finished bool
	err      error
}

// These are populated by init() functions in the providers' files
var providerMakers = make(map[string]func(string, string, *viper.Viper) (HostProvider, error))
var providerMagic = make(map[string]func(string) []HostProvider)

func NewRegistry(dataDir string) *Registry {
	return &Registry{
		providers: []HostProvider{},
		hosts:     Hosts{},
		sort:      []string{"name"},
		dataDir:   dataDir,
	}
}

func (r *Registry) LoadMagicProviders() {
	// Initialize all the magic providers
	for name, callable := range providerMagic {
		for _, p := range callable(r.dataDir) {
			r.AddProvider(p)
		}
	}
}

func (r *Registry) LoadProviders(c *viper.Viper) error {
	rerr := &MultiError{Subject: "Errors loading providers"}

	// And now all the explicitely configured ones
	for key, _ := range c.AllSettings() {
		ps := c.Sub(key)
		pname := ps.GetString("Provider")
		maker, ok := providerMakers[strings.ToLower(pname)]
		if !ok {
			rerr.Add(fmt.Errorf("No such provider: %s", pname))
			continue
		}
		p, err := maker(r.dataDir, pname, ps)
		if err != nil {
			rerr.Add(err)
		} else {
			r.AddProvider(p)
		}
	}
	if rerr.HasErrors() {
		return rerr
	}
	return nil
}

func (r *Registry) AddProvider(p HostProvider) {
	r.providers = append(r.providers, p)
}

type loadresult struct {
	provider string
	hosts    []*Host
	err      error
}

func (r *Registry) LoadHosts(mc chan CacheMessage) error {
	if len(r.hosts) > 0 {
		return nil
	}
	ctx := context.Background()
	rc := make(chan loadresult)
	if mc == nil {
		mc = make(chan CacheMessage)
		go func() {
			for range mc {
			}
		}()
	}
	defer close(mc)
	defer close(rc)

	for _, p := range r.providers {
		go func(pr HostProvider) {
			hosts, err := pr.Load(ctx, mc)
			rc <- loadresult{provider: pr.String(), hosts: hosts, err: err}
		}(p)
	}

	rerr := &MultiError{Subject: "Errors querying providers"}
	todo := len(r.providers)
	seen := make(map[string]int)
	hosts := make(Hosts, 0)

	for todo > 0 {
		r := <-rc
		logrus.Debugf("%d hosts returned from %s", len(r.hosts), r.provider)
		if r.err != nil {
			rerr.Add(r.err)
		}
		for _, host := range r.hosts {
			if existing, ok := seen[host.Name]; ok {
				hosts[existing].Amend(host)
			} else {
				seen[host.Name] = len(hosts)
				hosts = append(hosts, host)
			}
		}
		todo -= 1
	}
	r.hosts = hosts
	r.hosts.Sort(r.sort)
	if !rerr.HasErrors() {
		return nil
	}
	return rerr
}

func (r *Registry) GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts {
	ret := make(Hosts, 0)
	for _, host := range r.hosts {
		if host.Match(hostnameGlob, attributes) {
			ret = append(ret, host)
		}
	}
	if len(ret) == 0 && len(attributes) == 0 {
		if _, err := net.LookupHost(hostnameGlob); err == nil {
			ret = append(ret, NewHost(hostnameGlob, HostAttributes{}))
		}
	}
	return ret
}

func (r *Registry) SetSortFields(s []string) {
	r.sort = s
}

func (hosts Hosts) String() string {
	var ret strings.Builder
	for i, h := range hosts {
		if i > 0 {
			ret.WriteString(", ")
		}
		ret.WriteString(h.Name)
	}
	return ret.String()
}

func (h Hosts) Sort(fields []string) {
	if len(h) < 2 || len(fields) < 1 {
		return
	}
	// Most common and default case
	if len(fields) == 1 && fields[0] == "name" {
		sort.Slice(h, func(i, j int) bool { return h[i].Name < h[j].Name })
	} else {
		sort.Slice(h, func(i, j int) bool {
			return h[i].less(h[j], fields)
		})
	}
}

func (h Hosts) Uniq() Hosts {
	if len(h) < 2 {
		return h
	}
	src, dst := 1, 0
	for src < len(h) {
		if h[src].Name != h[dst].Name {
			dst += 1
			if dst != src {
				h[dst] = h[src]
			}
		}
		src += 1
	}
	return h[:dst+1]
}

func (h Hosts) maxLen() int {
	hlen := 0
	for _, host := range h {
		if len(host.Name) > hlen {
			hlen = len(host.Name)
		}
	}
	return hlen
}
