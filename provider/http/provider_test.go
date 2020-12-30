package http

import (
	"context"
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
	p.ParseViper(v)
	seenLoadingMessage := false
	hosts, err := p.Load(context.Background(), func(string, bool, error) { seenLoadingMessage = true })

	if err != nil {
		t.Errorf("HTTP fetch produced an error: %s", err)
	} else if !seenLoadingMessage {
		t.Errorf("HTTP fetch did not produce a loading message")
	} else if len(mhosts) != len(hosts) {
		t.Errorf("HTTP fetch returned the wrong number of hosts")
	} else {
		for i := 0; i < len(hosts); i++ {
			if hosts[i].Name != mhosts[i].Name {
				t.Errorf("Hostname mismatch: %s != %s", hosts[i].Name, mhosts[i].Name)
			}
			if hosts[i].Attributes["number"] != mhosts[i].Attributes["number"] {
				t.Errorf("Number mismatch: %d != %d", hosts[i].Attributes["number"], mhosts[i].Attributes["number"])
			}
		}
	}

}

func mockData() (herd.Hosts, string) {
	nhosts := 10
	hosts := make(herd.Hosts, nhosts)
	for i := 0; i < nhosts; i++ {
		h := herd.NewHost(fmt.Sprintf("host-%d.example.com", i), herd.HostAttributes{"number": int64(i)})
		hosts[i] = h
	}
	j, _ := json.Marshal(hosts)
	return hosts, string(j)
}
