package cache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

type fakeProvider struct {
	loaded  int
	doError bool
}

func (p *fakeProvider) Name() string {
	return "fake"
}

func (p *fakeProvider) Prefix() string {
	return "fake:"
}

func (p *fakeProvider) Equivalent(o herd.HostProvider) bool {
	return false
}

func (p *fakeProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	p.loaded++
	hosts := new(herd.HostSet)
	hosts.AddHost(herd.NewHost("test-host", "", herd.HostAttributes{"foo": "bar"}))
	if p.doError {
		return hosts, fmt.Errorf("You wanted an error")
	}
	return hosts, nil
}

func (p *fakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestCache(t *testing.T) {
	tmpdir := t.TempDir()
	c := NewFromProvider(&fakeProvider{}).(*Cache)
	c.config.Lifetime = 1 * time.Hour
	c.SetCacheDir(filepath.Join(tmpdir, "cache"))
	hosts, err := c.Load(t.Context(), func(string, bool, error) {})
	if hosts.Len() != 1 || err != nil {
		t.Errorf("First cache load did not succeed")
	}
	if c.source.(*fakeProvider).loaded != 1 {
		t.Errorf("First cache load did not appear to happen")
	}
	if c.mustRefresh() {
		t.Errorf("We must immediately refresh the cache")
	}
	hosts, err = c.Load(t.Context(), func(string, bool, error) {})
	if hosts.Len() != 1 || err != nil {
		t.Errorf("Second cache load did not succeed")
	}
	if c.source.(*fakeProvider).loaded != 1 {
		t.Errorf("Second cache load went to the backend provider")
	}

	c.source.(*fakeProvider).doError = true
	c.config.File += "-failed"
	c.Invalidate()
	_, err = c.Load(t.Context(), func(string, bool, error) {})
	if err == nil {
		t.Errorf("Expected an error from the cache")
	}
	if c.source.(*fakeProvider).loaded != 2 {
		t.Errorf("Fake provider was not called")
	}
	if _, err := os.Stat(c.config.File); err == nil {
		t.Errorf("Hosts were cached when the source provider errored out")
	}
}

func TestRelativeFiles(t *testing.T) {
	tmpdir := t.TempDir()
	r := herd.NewRegistry("/foo", filepath.Join(tmpdir, "cache"))
	p := &fakeProvider{}
	c := NewFromProvider(p)
	c.(*Cache).config.Prefix = "cache:"
	r.AddProvider(c)
	if c.(*Cache).config.File != filepath.Join(tmpdir, "cache", "fake.cache") {
		t.Errorf("Proper cache path not set, found %s", c.(*Cache).config.File)
	}

	if err := r.LoadHosts(t.Context(), func(string, bool, error) {}); err != nil {
		t.Errorf("Registry load did not succeed")
	}
	hosts := r.Search("", nil, nil, 0)
	if hosts.Len() != 1 {
		panic("Test broken")
	}

	h := hosts.Get(0)
	if attr, ok := h.Attributes["cache:fake:foo"]; !ok || attr.(string) != "bar" {
		t.Errorf("Expected attribute not found in %v", h.Attributes)
	}
}
