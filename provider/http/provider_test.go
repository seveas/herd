package http

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/seveas/herd"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

func TestHttpProvider(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mhosts, data := mockData()
	httpmock.RegisterResponder("GET", "http://inventory.example.com/inventory",
		httpmock.NewStringResponder(200, data))

	p := NewProvider("http")
	v := viper.New()
	v.SetDefault("Url", "http://inventory.example.com/inventory")
	if err := p.ParseViper(v); err != nil {
		t.Errorf("ParseViper failed %s", err)
	}
	seenLoadingMessage := false
	hosts, err := p.Load(t.Context(), func(string, bool, error) { seenLoadingMessage = true })

	if err != nil {
		t.Errorf("HTTP fetch produced an error: %s", err)
	} else if !seenLoadingMessage {
		t.Errorf("HTTP fetch did not produce a loading message")
	} else if mhosts.Len() != hosts.Len() {
		t.Errorf("HTTP fetch returned the wrong number of hosts")
	} else {
		for i := 0; i < hosts.Len(); i++ {
			h1 := hosts.Get(i)
			h2 := mhosts.Get(i)
			if h1.Name != h2.Name {
				t.Errorf("Hostname mismatch: %s != %s", h1.Name, h2.Name)
			}
			if h1.Attributes["number"] != h2.Attributes["number"] {
				t.Errorf("Number mismatch: %d != %d", h1.Attributes["number"], h2.Attributes["number"])
			}
		}
	}
}

func mockData() (*herd.HostSet, string) {
	nhosts := 10
	hosts := new(herd.HostSet)
	for i := 0; i < nhosts; i++ {
		h := herd.NewHost(fmt.Sprintf("host-%d.example.com", i), "", herd.HostAttributes{"number": int64(i)})
		hosts.AddHost(h)
	}
	j, _ := json.Marshal(hosts)
	return hosts, string(j)
}
