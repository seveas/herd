package katyusha

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
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

func TestCache(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "katyusha-test-cache-")
	mc := make(chan CacheMessage, 20)
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)
	c := Cache{
		Lifetime: 1 * time.Hour,
		File:     path.Join(tmpdir, "cache", "cache-test.cache"),
		Provider: &fakeProvider{},
	}
	hosts, _ := c.Load(nil, mc)
	if len(hosts) != 1 {
		t.Errorf("First cache load did not succeed")
	}
	if c.Provider.(*fakeProvider).loaded != 1 {
		t.Errorf("First cache load did not appear to happen")
	}
	if c.mustRefresh() {
		t.Errorf("We must immediately refresh the cache")
	}
	hosts, _ = c.Load(nil, mc)
	if len(hosts) != 1 {
		t.Errorf("Second cache load did not succeed")
	}
	if c.Provider.(*fakeProvider).loaded != 1 {
		t.Errorf("Second cache load went to the backend provider")
	}
}
