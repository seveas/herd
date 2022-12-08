package herd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"syscall"

	"github.com/hashicorp/go-plugin"
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
	hosts     *HostSet
	dataDir   string
	cacheDir  string
}

type HostProvider interface {
	Name() string
	Prefix() string
	ParseViper(v *viper.Viper) error
	Load(ctx context.Context, l LoadingMessage) (*HostSet, error)
	Equivalent(p HostProvider) bool
}

type DataLoader interface {
	SetDataDir(string) error
}

type Cache interface {
	Source() HostProvider
	Invalidate()
	Keep()
	SetCacheDir(string)
}

func NewRegistry(dataDir, cacheDir string) *Registry {
	return &Registry{
		providers: []HostProvider{},
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

	// And now all the explicitly configured ones
	for key := range c.AllSettings() {
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
		_ = c.SetDataDir(r.dataDir)
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

func (r *Registry) KeepCaches() {
	for _, p := range r.providers {
		if c, ok := p.(Cache); ok {
			c.Keep()
		}
	}
}

func (r *Registry) LoadHosts(ctx context.Context, lm LoadingMessage) error {
	if r.hosts != nil {
		return fmt.Errorf("Hosts have already been loaded")
	}
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for sig := range ch {
			plugin.CleanupClients()
			if sn, ok := sig.(syscall.Signal); ok {
				os.Exit(128 + int(sn))
			} else {
				os.Exit(1)
			}
		}
	}()
	defer func() {
		plugin.CleanupClients()
		signal.Stop(ch)
		close(ch)
		signal.Reset()
	}()

	sg := scattergather.New[*HostSet](int64(len(r.providers)))

	for _, p := range r.providers {
		p := p
		sg.Run(ctx, func() (*HostSet, error) {
			hosts, err := p.Load(ctx, lm)
			lm(p.Name(), true, err)
			if err != nil {
				return hosts, err
			}
			logrus.Debugf("%d hosts returned from %s", len(hosts.hosts), p.Name())
			for _, host := range hosts.hosts {
				if p.Prefix() != "" {
					host.Attributes = host.Attributes.prefix(p.Prefix())
				}
				host.Attributes["herd_provider"] = []string{p.Name()}
			}
			return hosts, err
		})
	}

	hostSets, err := sg.Wait()
	lm("", true, err)
	if err != nil {
		return err
	}
	r.hosts = MergeHostSets(hostSets)
	return nil
}

func (r *Registry) Search(hostnameGlob string, attributes MatchAttributes, sampled []string, count int) *HostSet {
	if strings.HasPrefix(hostnameGlob, "file:") {
		return r.getHostsFromFile(hostnameGlob[5:], attributes)
	}
	ret := r.hosts.Search(hostnameGlob, attributes)
	if len(ret.hosts) == 0 && attributes != nil && len(attributes) == 0 {
		if _, err := net.LookupHost(hostnameGlob); err == nil {
			ret = &HostSet{hosts: []*Host{NewHost(hostnameGlob, "", HostAttributes{})}}
		}
	}
	if len(sampled) != 0 {
		ret.Sample(sampled, count)
	}
	return ret
}

func (r *Registry) getHostsFromFile(fn string, attributes MatchAttributes) *HostSet {
	hosts := make([]*Host, 0)
	seen := make(map[string]int)
	data, err := os.ReadFile(fn)
	if err != nil {
		logrus.Errorf("Error reading %s: %s", fn, err)
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for i, host := range r.hosts.hosts {
		seen[host.Name] = i
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if i, ok := seen[line]; ok {
			host := r.hosts.hosts[i]
			if host.Match("*", attributes) {
				hosts = append(hosts, host)
			}
		} else {
			logrus.Warnf("Host %s not found", line)
		}
	}
	return &HostSet{hosts: hosts}
}

func (r *Registry) Settings() (string, map[string]interface{}) {
	providers := make([]string, len(r.providers))
	for i, p := range r.providers {
		providers[i] = p.Name()
	}
	return "Registry", map[string]interface{}{
		"Providers": providers,
	}
}
