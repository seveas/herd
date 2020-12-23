package katyusha

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var availableProviders = map[string]func(string) HostProvider{
	"http":        NewHttpProvider,
	"json":        NewJsonProvider,
	"plain":       NewPlainTextProvider,
	"known_hosts": NewKnownHostsProvider,
}

func Providers() []string {
	ret := []string{}
	for k := range availableProviders {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

var magicProviders = map[string]func(*Registry){}

func RegisterProvider(name string, constructor func(string) HostProvider, magic func(*Registry)) {
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
	Load(ctx context.Context, mc chan CacheMessage) (Hosts, error)
	Equivalent(p HostProvider) bool
}

type DataLoader interface {
	SetDataDir(string)
}

type Cache interface {
	Source() HostProvider
	Invalidate()
	SetCacheDir(string)
}

type CacheMessage struct {
	Name     string
	Finished bool
	Err      error
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
		return nil, fmt.Errorf("No such provider: %s", pname)
	}
	return c(name), nil
}

func (r *Registry) LoadMagicProviders() {
	// We always want these to be done first, so they're not implementing magic providerness themselves
	r.AddMagicProvider(NewKnownHostsProvider("known_hosts"))
	fn := filepath.Join(r.dataDir, "inventory")
	if _, err := os.Stat(fn); err == nil {
		p := NewPlainTextProvider("inventory")
		p.(*PlainTextProvider).config.File = "inventory"
		r.AddMagicProvider(p)
	}
	fn += ".json"
	if _, err := os.Stat(fn); err == nil {
		p := NewJsonProvider("inventory")
		p.(*JsonProvider).config.File = "inventory.json"
		r.AddMagicProvider(p)
	}
	// And now we do the other magic ones
	for _, fnc := range magicProviders {
		fnc(r)
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
	if c, ok := p.(DataLoader); ok {
		c.SetDataDir(r.dataDir)
	}
	r.providers = append(r.providers, p)
}

func (r *Registry) AddMagicProvider(p HostProvider) {
	sp := stripCache(p)
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

type loadresult struct {
	provider HostProvider
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
		go func(pr HostProvider, ctx context.Context) {
			hosts, err := pr.Load(ctx, mc)
			rc <- loadresult{provider: pr, hosts: hosts, err: err}
		}(p, ctx)
	}

	rerr := &MultiError{Subject: "Errors querying providers"}
	todo := len(r.providers)
	seen := make(map[string]int)
	hosts := make(Hosts, 0)

	for todo > 0 {
		r := <-rc
		logrus.Debugf("%d hosts returned from %s", len(r.hosts), r.provider.Name())
		if r.err != nil {
			rerr.Add(r.err)
		}
		for _, host := range r.hosts {
			if r.provider.Prefix() != "" {
				host.Attributes = host.Attributes.prefix(r.provider.Prefix())
			}
			host.Attributes["katyusha_provider"] = []string{r.provider.Name()}
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
