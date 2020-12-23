package herd

import (
	"net"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(dataDir("0"), cacheDir("0"))
	if len(r.providers) > 0 {
		t.Errorf("got %d providers, expected none", len(r.providers))
	}
}

func TestMagicProviders(t *testing.T) {
	defer os.Setenv("HOME", realUserHome)

	os.Setenv("HOME", homeDir("1"))
	os.Unsetenv("CONSUL_HTTP_ADDR")
	r := NewRegistry(dataDir("1"), cacheDir("1"))
	r.LoadMagicProviders()
	expect := 1
	_, err := net.LookupHost("consul.service.consul")
	if err == nil {
		expect++
	}

	if len(r.providers) != expect {
		t.Errorf("Got %d providers, expected %d", len(r.providers), expect)
	}
	if _, ok := r.providers[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(r.providers[0]))
	}

	os.Setenv("HOME", homeDir("2"))
	r = NewRegistry(dataDir("2"), cacheDir("2"))
	r.LoadMagicProviders()
	expect++
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
	c := NewCacheFromProvider(p)
	r.AddProvider(c)
	if c.(*Cache).config.File != filepath.Join(cacheDir("2"), "ittest.cache") {
		t.Errorf("Proper cache path not set, found %s", c.(*Cache).config.File)
	}
}
