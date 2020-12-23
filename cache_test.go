package katyusha

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

type fakeProvider struct {
	loaded int
}

func (p *fakeProvider) Name() string {
	return "fake"
}

func (p *fakeProvider) Prefix() string {
	return ""
}

func (p *fakeProvider) Equivalent(o HostProvider) bool {
	return false
}

func (p *fakeProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	p.loaded++
	return Hosts{NewHost("test-host", HostAttributes{})}, nil
}

func (p *fakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestCache(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "katyusha-test-cache-")
	mc := make(chan CacheMessage, 20)
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)
	c := NewCacheFromProvider(&fakeProvider{}).(*Cache)
	c.config.Lifetime = 1 * time.Hour
	c.SetCacheDir(filepath.Join(tmpdir, "cache"))
	hosts, _ := c.Load(nil, mc)
	if len(hosts) != 1 {
		t.Errorf("First cache load did not succeed")
	}
	if c.source.(*fakeProvider).loaded != 1 {
		t.Errorf("First cache load did not appear to happen")
	}
	if c.mustRefresh() {
		t.Errorf("We must immediately refresh the cache")
	}
	hosts, _ = c.Load(nil, mc)
	if len(hosts) != 1 {
		t.Errorf("Second cache load did not succeed")
	}
	if c.source.(*fakeProvider).loaded != 1 {
		t.Errorf("Second cache load went to the backend provider")
	}
}
