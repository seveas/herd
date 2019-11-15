package herd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type HostProvider interface {
	GetHosts(hostnameGlob string) Hosts
}

type Cacher interface {
	Cache(ctx context.Context) error
}

type Providers []HostProvider

var ProviderMakers = make(map[string]func(string, *viper.Viper) (HostProvider, error))
var ProviderMagic = make(map[string]func(Providers) Providers)

func LoadProviders() (Providers, error) {
	ret := make(Providers, 0)

	// Initialize all the magic providers
	for _, callable := range ProviderMagic {
		ret = callable(ret)
	}

	// And now all the explicitely configured ones
	providers := viper.Sub("Providers")
	if providers != nil {
		for key, _ := range providers.AllSettings() {
			ps := providers.Sub(key)
			pname := ps.GetString("Provider")
			maker, ok := ProviderMakers[strings.ToLower(pname)]
			if !ok {
				return nil, fmt.Errorf("No such provider: %s", pname)
			}
			p, err := maker(pname, ps)
			if err != nil {
				return nil, err
			}
			ret = append(ret, p)
		}
	}
	return ret, nil
}

func (p *Providers) Cache() []error {
	if err := os.MkdirAll(viper.GetString("CacheDir"), 0700); err != nil {
		return []error{err}
	}
	ctx := context.Background()
	errs := make([]error, 0)
	todo := 0
	ec := make(chan error)
	for _, pr := range *p {
		if c, ok := pr.(Cacher); ok {
			todo += 1
			go func(pr Cacher) {
				err := c.Cache(ctx)
				if err != nil {
					err = fmt.Errorf("Error contacting %s: %s", pr, err)
				}
				ec <- err
			}(c)
		}
	}
	for todo > 0 {
		err := <-ec
		if err != nil {
			errs = append(errs, err)
		}
		todo -= 1
	}
	return errs
}

func (p *Providers) GetHosts(hostnameGlob string, attributes MatchAttributes) Hosts {
	hosts := make(Hosts, 0)
	seen := make(map[string]int)
	for _, provider := range *p {
		for _, host := range provider.GetHosts(hostnameGlob) {
			if existing, ok := seen[host.Name]; ok {
				hosts[existing].Amend(host)
			} else {
				seen[host.Name] = len(hosts)
				hosts = append(hosts, host)
			}
		}
	}
	ret := make(Hosts, 0)
	for _, host := range hosts {
		if host.Match(hostnameGlob, attributes) {
			ret = append(ret, host)
		}
	}
	return ret
}
