package herd

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(dataDir("0"), cacheDir("0"))
	if len(r.providers) > 0 {
		t.Errorf("got %d providers, expected none", len(r.providers))
		return
	}
}

func TestMagicProviders(t *testing.T) {
	defer os.Setenv("HOME", realUserHome)

	os.Setenv("HOME", homeDir("1"))
	r := NewRegistry(dataDir("1"), cacheDir("1"))
	r.LoadMagicProviders()
	if len(r.providers) != 1 {
		t.Errorf("Got %d providers, expected 1", len(r.providers))
		t.Errorf("%v", r.providers)
	}
	if _, ok := r.providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.providers[0]))
	}

	os.Setenv("HOME", homeDir("2"))
	r = NewRegistry(dataDir("2"), cacheDir("2"))
	r.LoadMagicProviders()
	if len(r.providers) != 2 {
		t.Errorf("Got %d providers, expected 2", len(r.providers))
	}
	if _, ok := r.providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.providers[0]))
	}
}

type FakeProvider struct {
	baseProvider `mapstructure:",squash"`
}

func (p *FakeProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	return Hosts{NewHost("hostname", HostAttributes{})}, nil
}

func (p *FakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestGetHosts(t *testing.T) {
	r := Registry{providers: []HostProvider{&FakeProvider{}, &FakeProvider{}}}
	err := r.LoadHosts(nil)
	if err != nil {
		t.Errorf("%t %v", err, err)
		t.Errorf("Could not load hosts: %s", err.Error())
	}
	if len(r.hosts) != 1 {
		t.Errorf("Hosts returned by multiple providers are not merged, got %d hosts instead of 1", len(r.hosts))
	}
}

func TestRelativeFiles(t *testing.T) {
	r := NewRegistry(dataDir("2"), cacheDir("2"))
	p := &PlainTextProvider{File: "inventory", baseProvider: baseProvider{Name: "ittest"}}
	r.AddProvider(p)
	if p.File != filepath.Join(dataDir("2"), "inventory") {
		t.Errorf("Filepath did not get interpreted relative to dataDir")
	}
	c := &Cache{baseProvider: baseProvider{Name: "itcache"}, Source: p}
	r.AddProvider(c)
	if c.File != filepath.Join(cacheDir("2"), "itcache.cache") {
		t.Errorf("Proper cache path not set, found %s", c.File)
	}
	c2 := &Cache{baseProvider: baseProvider{Name: "itcache"}, Source: p, File: "it-cache.cache"}
	r.AddProvider(c2)
	if c2.File != filepath.Join(cacheDir("2"), "it-cache.cache") {
		t.Errorf("Proper cache path not set, found %s", c2.File)
	}
}
