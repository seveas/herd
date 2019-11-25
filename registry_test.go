package katyusha

import (
	"context"
	"reflect"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r, _ := NewRegistry()
	if len(r.Providers) < 1 {
		t.Errorf("got %d providers, expected at least 1", len(r.Providers))
		return
	}
	if _, ok := r.Providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.Providers[0]))
	}
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
	UI = NewSimpleUI()
	r := Registry{Providers: []HostProvider{&FakeProvider{}, &FakeProvider{}}}
	err := r.Load()
	if err != nil {
		t.Errorf("%t %v", err, err)
		t.Errorf("Could not load hosts: %s", err.Error())
	}
	if len(r.Hosts) != 1 {
		t.Errorf("Hosts returned by multiple providers are not merged, got %d hosts instead of 1", len(r.Hosts))
	}
}
