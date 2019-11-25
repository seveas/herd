package katyusha

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Registry struct {
	Providers []HostProvider
	Hosts     Hosts
}

type HostProvider interface {
	Load(ctx context.Context, mc chan CacheMessage) (Hosts, error)
	String() string
}

type PostProcessor interface {
	PostProcess(Hosts)
}

type CacheMessage struct {
	name     string
	finished bool
	err      error
}

// These are populated by init() functions in the providers' files
var ProviderMakers = make(map[string]func(string, *viper.Viper) (HostProvider, error))
var ProviderMagic = make(map[string]func() []HostProvider)

func NewRegistry() (*Registry, error) {
	ret := &Registry{
		Providers: []HostProvider{},
		Hosts:     Hosts{},
	}
	rerr := &MultiError{}

	// Initialize all the magic providers
	for _, callable := range ProviderMagic {
		p := callable()
		ret.Providers = append(ret.Providers, p...)
	}

	// And now all the explicitely configured ones
	providers := viper.Sub("Providers")
	if providers != nil {
		for key, _ := range providers.AllSettings() {
			ps := providers.Sub(key)
			pname := ps.GetString("Provider")
			maker, ok := ProviderMakers[strings.ToLower(pname)]
			if !ok {
				rerr.Add(fmt.Errorf("No such provider: %s", pname))
				continue
			}
			p, err := maker(pname, ps)
			if err != nil {
				rerr.Add(err)
				continue
			}
			ret.Providers = append(ret.Providers, p)
		}
	}
	if len(rerr.Errors) != 0 {
		return nil, rerr
	}
	return ret, nil
}

type loadresult struct {
	provider string
	hosts    []*Host
	err      error
}

func (r *Registry) Load() error {
	if len(r.Hosts) > 0 {
		return nil
	}
	ctx := context.Background()
	mc := make(chan CacheMessage)
	rc := make(chan loadresult)
	defer close(mc)
	defer close(rc)

	st := time.Now()
	for _, p := range r.Providers {
		go func(pr HostProvider) {
			hosts, err := pr.Load(ctx, mc)
			rc <- loadresult{provider: pr.String(), hosts: hosts, err: err}
		}(p)
	}

	caches := make([]string, 0)
	rerr := &MultiError{}
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	todo := len(r.Providers)
	seen := make(map[string]int)
	hosts := make(Hosts, 0)

	for todo > 0 {
		select {
		case msg := <-mc:
			if msg.err != nil && msg.err.Error() != "" {
				err := fmt.Errorf("Error contacting %s: %s", msg.name, msg.err)
				UI.Errorf("%s", err.Error())
			}
			if msg.finished {
				UI.Debugf("Cache updated for %s", msg.name)
				for i, v := range caches {
					if v == msg.name {
						caches = append(caches[:i], caches[i+1:]...)
						break
					}
				}
			} else {
				caches = append(caches, msg.name)
			}
			UI.CacheProgress(st, caches)
		case r := <-rc:
			UI.Debugf("%d hosts returned from %s", len(r.hosts), r.provider)
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
		case <-ticker.C:
			if len(caches) > 0 {
				UI.CacheProgress(st, caches)
			}
		}
	}
	r.Hosts = hosts
	r.Hosts.Sort()
	if len(rerr.Errors) == 0 {
		return nil
	}
	return rerr
}

func (r *Registry) GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts {
	ret := make(Hosts, 0)
	for _, host := range r.Hosts {
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
