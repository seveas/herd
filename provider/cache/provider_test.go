package cache

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/seveas/katyusha"

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

func (p *fakeProvider) Equivalent(o katyusha.HostProvider) bool {
	return false
}

func (p *fakeProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (katyusha.Hosts, error) {
	p.loaded++
	h := katyusha.NewHost("test-host", "", katyusha.HostAttributes{"foo": "bar"})
	if p.doError {
		return katyusha.Hosts{h}, fmt.Errorf("You wanted an error")
	}
	return katyusha.Hosts{h}, nil
}

func (p *fakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestCache(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "katyusha-test-cache-")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)
	c := NewFromProvider(&fakeProvider{}).(*Cache)
	c.config.Lifetime = 1 * time.Hour
	c.SetCacheDir(filepath.Join(tmpdir, "cache"))
	hosts, err := c.Load(nil, func(string, bool, error) {})
	if len(hosts) != 1 || err != nil {
		t.Errorf("First cache load did not succeed")
	}
	if c.source.(*fakeProvider).loaded != 1 {
		t.Errorf("First cache load did not appear to happen")
	}
	if c.mustRefresh() {
		t.Errorf("We must immediately refresh the cache")
	}
	hosts, err = c.Load(nil, func(string, bool, error) {})
	if len(hosts) != 1 || err != nil {
		t.Errorf("Second cache load did not succeed")
	}
	if c.source.(*fakeProvider).loaded != 1 {
		t.Errorf("Second cache load went to the backend provider")
	}

	c.source.(*fakeProvider).doError = true
	c.config.File += "-failed"
	c.Invalidate()
	hosts, err = c.Load(nil, func(string, bool, error) {})
	if err == nil {
		t.Errorf("Expected an error from the cache")
	}
	if c.source.(*fakeProvider).loaded != 2 {
		t.Errorf("Fake provider was not called")
	}
	if _, err := os.Stat(filepath.Join(c.config.File)); err == nil {
		t.Errorf("Hosts were cached when the source provider errored out")
	}

}

func TestRelativeFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "katyusha-test-cache-")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)
	r := katyusha.NewRegistry("/foo", filepath.Join(tmpdir, "cache"))
	p := &fakeProvider{}
	c := NewFromProvider(p)
	c.(*Cache).config.Prefix = "cache:"
	r.AddProvider(c)
	if c.(*Cache).config.File != filepath.Join(tmpdir, "cache", "fake.cache") {
		t.Errorf("Proper cache path not set, found %s", c.(*Cache).config.File)
	}

	if err := r.LoadHosts(func(string, bool, error) {}); err != nil {
		t.Errorf("Registry load did not succeed")
	}
	hosts := r.GetHosts("", nil)
	if len(hosts) != 1 {
		panic("Test broken")
	}

	if attr, ok := hosts[0].Attributes["cache:fake:foo"]; !ok || attr.(string) != "bar" {
		t.Errorf("Expected attribute not found in %v", hosts[0].Attributes)
	}
}
