package herd

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

type HostProvider interface {
	GetHosts(hostnameGlob string) Hosts
}

type Cacher interface {
	Cache(ctx *context.Context) error
}

type Providers []HostProvider

func LoadProviders() (Providers, error) {
	files := []string{"/etc/ssh/ssh_known_hosts"}
	home, err := homedir.Dir()
	if err == nil {
		files = append(files, path.Join(home, ".ssh", "known_hosts"))
	}
	ret := Providers{
		&KnownHostsProvider{
			Files: files,
		},
		&CliProvider{},
	}
	providers := viper.Sub("Providers")
	if providers != nil {
		for key, _ := range providers.AllSettings() {
			ps := providers.Sub(key)
			pname := ps.GetString("Provider")
			var provider HostProvider
			if pname == "json" {
				provider = &JsonProvider{}
			} else if pname == "http" {
				provider = &HttpProvider{
					Name: key,
					File: path.Join(viper.GetString("CacheDir"), key+".cache"),
				}
			}
			err := ps.Unmarshal(provider)
			if err != nil {
				return nil, err
			}
			ret = append(ret, provider)
		}
	}
	return ret, nil
}

func (p *Providers) Cache() []error {
	if err := os.MkdirAll(viper.GetString("CacheDir"), 0700); err != nil {
		return []error{err}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	errs := make([]error, 0)
	todo := 0
	ec := make(chan error)
	for _, pr := range *p {
		if c, ok := pr.(Cacher); ok {
			todo += 1
			go func(pr Cacher) {
				err := c.Cache(&ctx)
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
