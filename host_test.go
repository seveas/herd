package herd

import (
	"encoding/json"
	"testing"
)

func TestHostDeserialization(t *testing.T) {
	data := map[string]interface{}{
		"Name": "test-host",
	}
	bdata, _ := json.Marshal(data)
	var host Host
	err := json.Unmarshal(bdata, &host)
	if err != nil {
		t.Errorf("Unable to deserialize host data: %s", err)
	}
	host.init()
	if host.Attributes == nil {
		t.Errorf("Deserialized host does not get attributes")
	}
	if host.Port != 22 {
		t.Errorf("Deserialized host does not get default port")
	}
}

func TestHostSorting(t *testing.T) {
	h1 := &Host{Name: "host-a.example.com", Attributes: HostAttributes{"site": "site1", "role": "db"}}
	h2 := &Host{Name: "host-b.example.com", Attributes: HostAttributes{"site": "site2", "role": "db"}}
	h3 := &Host{Name: "host-c.example.com", Attributes: HostAttributes{"site": "site1", "role": "app"}}
	h4 := &Host{Name: "host-d.example.com", Attributes: HostAttributes{"site": "site2", "role": "app"}}
	hosts := Hosts{h1, h2, h3, h4}
	h1.init()
	h2.init()

	if !h1.less(h2, []string{}) {
		t.Errorf("Sorting hosts with no fields is failing")
	}
	if h2.less(h1, []string{}) {
		t.Errorf("Sorting hosts with no fields is failing")
	}
	if !h1.less(h2, []string{"name"}) {
		t.Errorf("Sorting hosts by name is failing")
	}
	if h2.less(h1, []string{"name"}) {
		t.Errorf("Sorting hosts by name fields is failing")
	}

	hosts.Sort([]string{"site", "name"})
	if !eq(hosts, Hosts{h1, h3, h2, h4}) {
		t.Errorf("Sorting by site+name is failing, got %v", hosts)
	}

	hosts.Sort([]string{"site", "role"})
	if !eq(hosts, Hosts{h3, h1, h4, h2}) {
		t.Errorf("Sorting by site+role is failing, got %v", hosts)
	}
}

func eq(a, b Hosts) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
