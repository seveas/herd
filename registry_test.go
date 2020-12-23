package herd

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

type fakeProvider struct {
}

func (p *fakeProvider) Name() string {
	return "fake"
}

func (p *fakeProvider) Prefix() string {
	return "fake:"
}

func (p *fakeProvider) Equivalent(o HostProvider) bool {
	return false
}

func (p *fakeProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	h := NewHost("test-host", HostAttributes{"foo": "bar"})
	return Hosts{h}, nil
}

func (p *fakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(dataDir("0"), cacheDir("0"))
	if len(r.providers) > 0 {
		t.Errorf("got %d providers, expected none", len(r.providers))
	}
}

func TestMagicProviders(t *testing.T) {
	defer os.Setenv("HOME", realUserHome)

	os.Setenv("HOME", homeDir("1"))
	r := NewRegistry(dataDir("1"), cacheDir("1"))
	expect := 1
	r.LoadMagicProviders()
	if len(r.providers) != expect {
		t.Errorf("Got %d providers, expected %d", len(r.providers), expect)
	}
	if _, ok := r.providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.providers[0]))
	}

	os.Setenv("HOME", homeDir("2"))
	r = NewRegistry(dataDir("2"), cacheDir("2"))
	expect = 2
	r.LoadMagicProviders()
	if len(r.providers) != expect {
		t.Errorf("Got %d providers, expected %d", len(r.providers), expect)
	}
	if _, ok := r.providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.providers[0]))
	}
}

func TestGetHosts(t *testing.T) {
	r := Registry{providers: []HostProvider{&fakeProvider{}, &fakeProvider{}}}
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
	p := &PlainTextProvider{name: "ittest"}
	p.config.File = "inventory"
	r.AddProvider(p)
	if p.config.File != filepath.Join(dataDir("2"), "inventory") {
		t.Errorf("Filepath did not get interpreted relative to dataDir")
	}
}
