package herd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"reflect"
	"sort"
	"strings"

	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var availableProviders = make(map[string]func(string) HostProvider)
var magicProviders = make(map[string]func() HostProvider)

func Providers() []string {
	ret := []string{}
	for k := range availableProviders {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func RegisterProvider(name string, constructor func(string) HostProvider, magic func() HostProvider) {
	availableProviders[name] = constructor
	if magic != nil {
		magicProviders[name] = magic
	}
}

type Registry struct {
	providers []HostProvider
	hosts     Hosts
	sort      []string
	dataDir   string
	cacheDir  string
}

type Hosts []*Host

type HostProvider interface {
	Name() string
	Prefix() string
	ParseViper(v *viper.Viper) error
	Load(ctx context.Context, l LoadingMessage) (Hosts, error)
	Equivalent(p HostProvider) bool
}

type DataLoader interface {
	SetDataDir(string) error
}

type Cache interface {
	Source() HostProvider
	Invalidate()
	SetCacheDir(string)
}

func NewRegistry(dataDir, cacheDir string) *Registry {
	return &Registry{
		providers: []HostProvider{},
		hosts:     Hosts{},
		sort:      []string{"name"},
		dataDir:   dataDir,
		cacheDir:  cacheDir,
	}
}

func NewProvider(pname, name string) (HostProvider, error) {
	if pname == "" {
		return nil, fmt.Errorf("No provider specified")
	}
	c, ok := availableProviders[pname]
	if !ok {
		// Try finding a plugin by this name, if we have the provider plugin loaded
		if c, ok := availableProviders["plugin"]; ok {
			if _, err := exec.LookPath(fmt.Sprintf("herd-provider-%s", name)); err == nil {
				return c(name), nil
			}
		}
		return nil, fmt.Errorf("No such provider: %s", pname)
	}
	return c(name), nil
}

func (r *Registry) LoadMagicProviders() {
	for _, fnc := range magicProviders {
		if p := fnc(); p != nil {
			r.AddMagicProvider(p)
		}
	}
}

func (r *Registry) LoadProviders(c *viper.Viper) error {
	rerr := &MultiError{Subject: "Errors loading providers"}

	// And now all the explicitely configured ones
	for key, _ := range c.AllSettings() {
		ps := c.Sub(key)
		pname := ps.GetString("Provider")
		p, err := NewProvider(pname, key)
		if err != nil {
			rerr.Add(fmt.Errorf("Error parsing config for %s: %s", key, err))
		} else {
			err = p.ParseViper(ps)
			if err != nil {
				rerr.Add(fmt.Errorf("Error parsing config for %s: %s", key, err))
			} else {
				r.AddProvider(p)
			}
		}
	}
	if rerr.HasErrors() {
		return rerr
	}
	return nil
}

func (r *Registry) AddProvider(p HostProvider) {
	logrus.Debugf("Adding provider %s", p.Name())
	if c, ok := p.(Cache); ok {
		c.SetCacheDir(r.cacheDir)
	}
	if c, ok := stripCache(p).(DataLoader); ok {
		c.SetDataDir(r.dataDir)
	}
	r.providers = append(r.providers, p)
}

func (r *Registry) AddMagicProvider(p HostProvider) {
	sp := stripCache(p)
	if c, ok := sp.(DataLoader); ok {
		if err := c.SetDataDir(r.dataDir); err != nil {
			return
		}
	}
	for _, pr := range r.providers {
		pr := stripCache(pr)
		if reflect.TypeOf(sp) != reflect.TypeOf(pr) {
			continue
		}
		if sp.Equivalent(pr) {
			return
		}
	}
	r.AddProvider(p)
}

func stripCache(p HostProvider) HostProvider {
	if c, ok := p.(Cache); ok {
		return c.Source()
	}
	return p
}

func (r *Registry) InvalidateCache() {
	for _, p := range r.providers {
		if c, ok := p.(Cache); ok {
			c.Invalidate()
		}
	}
}

func (r *Registry) LoadHosts(lm LoadingMessage) error {
	if len(r.hosts) > 0 {
		return nil
	}
	ctx := context.Background()
	sg := scattergather.New(int64(len(r.providers)))

	for _, p := range r.providers {
		sg.Run(func(ctx context.Context, args ...interface{}) (interface{}, error) {
			p := args[0].(HostProvider)
			hosts, err := p.Load(ctx, lm)
			lm(p.Name(), true, err)
			logrus.Debugf("%d hosts returned from %s", len(hosts), p.Name())
			for _, host := range hosts {
				if p.Prefix() != "" {
					host.Attributes = host.Attributes.prefix(p.Prefix())
				}
				host.Attributes["herd_provider"] = []string{p.Name()}
			}
			return hosts, err
		}, ctx, p)
	}

	untypedHosts, err := sg.Wait()
	lm("", true, nil)
	seen := make(map[string]int)
	allHosts := make(Hosts, 0)

	if err != nil {
		return err
	}

	for _, uh := range untypedHosts {
		hosts := uh.(Hosts)
		for _, host := range hosts {
			if existing, ok := seen[host.Name]; ok {
				allHosts[existing].Amend(host)
			} else {
				seen[host.Name] = len(allHosts)
				allHosts = append(allHosts, host)
			}
		}
	}
	r.hosts = allHosts
	r.hosts.Sort(r.sort)
	return nil
}

func (r *Registry) GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts {
	ret := make(Hosts, 0)
	if strings.HasPrefix(hostnameGlob, "file:") {
		return r.getHostsFromFile(hostnameGlob[5:], attributes)
	}
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

func (r *Registry) getHostsFromFile(fn string, attributes MatchAttributes) Hosts {
	ret := make(Hosts, 0)
	seen := make(map[string]int)
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		logrus.Errorf("Error reading %s: %s", fn, err)
		return ret
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for i, host := range r.hosts {
		seen[host.Name] = i
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if i, ok := seen[line]; ok {
			host := r.hosts[i]
			if host.Match("*", attributes) {
				ret = append(ret, host)
			}
		} else {
			logrus.Warnf("Host %s not found", line)
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
