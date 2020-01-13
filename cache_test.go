package herd

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

func (p *fakeProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	p.loaded++
	return Hosts{NewHost("test-host", HostAttributes{})}, nil
}

func (p *fakeProvider) String() string {
	return "fake"
}

func (p *fakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestCache(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "herd-test-cache-")
	mc := make(chan CacheMessage, 20)
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)
	c := Cache{
		Lifetime: 1 * time.Hour,
		File:     filepath.Join(tmpdir, "cache", "cache-test.cache"),
		Source:   &fakeProvider{},
	}
	hosts, _ := c.Load(nil, mc)
	if len(hosts) != 1 {
		t.Errorf("First cache load did not succeed")
	}
	if c.Source.(*fakeProvider).loaded != 1 {
		t.Errorf("First cache load did not appear to happen")
	}
	if c.mustRefresh() {
		t.Errorf("We must immediately refresh the cache")
	}
	hosts, _ = c.Load(nil, mc)
	if len(hosts) != 1 {
		t.Errorf("Second cache load did not succeed")
	}
	if c.Source.(*fakeProvider).loaded != 1 {
		t.Errorf("Second cache load went to the backend provider")
	}
}
