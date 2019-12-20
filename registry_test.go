package katyusha

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if len(r.providers) > 0 {
		t.Errorf("got %d providers, expected none", len(r.providers))
		return
	}
}

func TestMagicProviders(t *testing.T) {
	defer os.Setenv("HOME", realUserHome)

	os.Setenv("HOME", filepath.Join(testDataRoot, "homes", "1"))
	r := NewRegistry()
	r.LoadMagicProviders()
	if len(r.providers) != 1 {
		t.Errorf("Got %d providers, expected 1", len(r.providers))
	}
	if _, ok := r.providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.providers[0]))
	}

	os.Setenv("HOME", filepath.Join(testDataRoot, "homes", "2"))
}

type FakeProvider struct {
}

func (p *FakeProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	return Hosts{NewHost("hostname", HostAttributes{})}, nil
}

func (p *FakeProvider) String() string {
	return "fake"
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
