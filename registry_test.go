package herd

import (
	"context"
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
