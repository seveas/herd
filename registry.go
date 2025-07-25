package herd

import (
	"bufio"
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
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
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
	providers    []HostProvider
	hosts        *HostSet
	globPrefixes map[string]func(string, *HostSet) (*HostSet, error)
	dataDir      string
	cacheDir     string
	magicLoaded  bool
}

type HostProvider interface {
	Name() string
	Prefix() string
	ParseViper(v *viper.Viper) error
	Load(ctx context.Context, l LoadingMessage) (*HostSet, error)
	Equivalent(p HostProvider) bool
}

type HostKeyProvider interface {
	LoadHostKeys(ctx context.Context, l LoadingMessage) (map[string][]ssh.PublicKey, error)
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
		globPrefixes: map[string]func(string, *HostSet) (*HostSet, error){
			"file:": fileFilter,
		},
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
	r.magicLoaded = true
}

func (r *Registry) LoadHostKeys(ctx context.Context, lm LoadingMessage) error {
	if r.magicLoaded {
		return nil
	}
	sg := scattergather.New[map[string][]ssh.PublicKey](int64(len(r.providers)))
	for _, fnc := range magicProviders {
		p := fnc()
		if hp, ok := p.(HostKeyProvider); ok {
			hp := hp
			sg.Run(context.Background(), func() (map[string][]ssh.PublicKey, error) {
				return hp.LoadHostKeys(ctx, lm)
			})
		}
	}
	allKeys, err := sg.Wait()
	if err != nil {
		return err
	}
	r.hosts.addHostKeys(allKeys)
	return nil
}

func (r *Registry) LoadProviders(c *viper.Viper) error {
	rerr := &MultiError{Subject: "Errors loading providers"}

	// And now all the explicitly configured ones
	for key := range c.AllSettings() {
		ps := c.Sub(key)
		pname := ps.GetString("Provider")
		if ps.IsSet("Enabled") && !ps.GetBool("Enabled") {
			logrus.Debugf("Skipping disabled provider %s", key)
			continue
		}
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
	sg.KeepAllResults(true)

	for _, p := range r.providers {
		sg.Run(ctx, func() (*HostSet, error) {
			hosts, err := p.Load(ctx, lm)
			lm(p.Name(), true, err)
			if err != nil && hosts == nil {
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
	r.hosts = MergeHostSets(hostSets)
	lm("", true, err)
	return err
}

func (r *Registry) Search(hostnameGlob string, attributes MatchAttributes, sampled []string, count int) *HostSet {
	ret := r.hosts
	for glob, fnc := range r.globPrefixes {
		if strings.HasPrefix(hostnameGlob, glob) {
			var err error
			ret, err = fnc(hostnameGlob[len(glob):], ret)
			if err != nil {
				logrus.Errorf("Error applying filter %s: %s", glob, err)
				ret = &HostSet{}
			}
			hostnameGlob = "*"
			break
		}
	}

	ret = ret.Search(hostnameGlob, attributes)
	if len(ret.hosts) == 0 && attributes != nil && len(attributes) == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if _, err := net.DefaultResolver.LookupHost(ctx, hostnameGlob); err == nil {
			ret = &HostSet{hosts: []*Host{NewHost(hostnameGlob, "", HostAttributes{})}}
		}
	}
	if len(sampled) != 0 {
		ret.Sort()
		ret.Sample(sampled, count)
	}
	return ret
}

func (r *Registry) AddGlobPrefix(prefix string, fnc func(string, *HostSet) (*HostSet, error)) {
	if _, ok := r.globPrefixes[prefix]; ok {
		logrus.Warnf("Overwriting glob prefix `%s`", prefix)
	}
	r.globPrefixes[prefix] = fnc
}

func (r *Registry) GlobPrefixes() []string {
	ret := make([]string, 0, len(r.globPrefixes))
	for k := range r.globPrefixes {
		ret = append(ret, k)
	}
	return ret
}

func fileFilter(fn string, hs *HostSet) (*HostSet, error) {
	seen := make(map[string]bool)
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		seen[strings.TrimSpace(scanner.Text())] = true
	}
	hs = hs.Filter(func(h *Host) bool {
		if _, ok := seen[h.Name]; !ok {
			return false
		}
		delete(seen, h.Name)
		return true
	})
	// We synthesize hosts that were not yet seen in the inventory
	for h := range seen {
		hs.hosts = append(hs.hosts, NewHost(h, "", HostAttributes{}))
	}
	return hs, nil
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
