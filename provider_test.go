package katyusha

import (
	"reflect"
	"testing"
)

func TestLoadProviders(t *testing.T) {
	ret := LoadProviders()
	if len(ret) != 2 {
		t.Errorf("got %d providers, expected 2", len(ret))
		return
	}
	if _, ok := ret[0].(*KnownHostsProvider); !ok {
		t.Errorf("expected the first provider to be the known hosts provider, not %s", reflect.TypeOf(ret[0]))
	}
	if _, ok := ret[1].(*CliProvider); !ok {
		t.Errorf("expected the first provider to be the cli provider, not %s", reflect.TypeOf(ret[1]))
	}
}

type FakeProvider struct {
}

func (p *FakeProvider) GetHosts(glob string, attrs MatchAttributes) Hosts {
	return Hosts{NewHost(glob, HostAttributes{})}
}

func TestGetHosts(t *testing.T) {
	p := Providers{&FakeProvider{}, &FakeProvider{}}
	hosts := p.GetHosts("hostname.domainname", MatchAttributes{})
	if len(hosts) != 1 {
		t.Error("Hosts returned by multiple providers are not merged")
	}
}
