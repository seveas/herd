package herd

import (
	"testing"
)

func TestHostSetSorting(t *testing.T) {
	h1 := NewHost("host-a.example.com", "", HostAttributes{"site": "site1", "role": "db"})
	h2 := NewHost("host-b.example.com", "", HostAttributes{"site": "site2", "role": "db"})
	h3 := NewHost("host-c.example.com", "", HostAttributes{"site": "site1", "role": "app"})
	h4 := NewHost("host-d.example.com", "", HostAttributes{"site": "site2", "role": "app"})
	hosts := HostSet{hosts: []*Host{h1, h2, h3, h4}}

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

	hosts.SetSortFields([]string{"site", "name"})
	hosts.Sort()
	if !eq(hosts.hosts, []*Host{h1, h3, h2, h4}) {
		t.Errorf("Sorting by site+name is failing, got %v", hosts)
	}

	hosts.SetSortFields([]string{"site", "role"})
	hosts.Sort()
	if !eq(hosts.hosts, []*Host{h3, h1, h4, h2}) {
		t.Errorf("Sorting by site+role is failing, got %v", hosts)
	}
}

func TestHostSetUniq(t *testing.T) {
	h1 := NewHost("host-a.example.com", "", HostAttributes{"site": "site1", "role": "db"})
	h2 := NewHost("host-a.example.com", "", HostAttributes{"site": "site1", "role": "db"})

	hosts := HostSet{hosts: []*Host{h1, h2}}
	hosts.Uniq()
	if len(hosts.hosts) != 1 {
		t.Errorf("Uniq did not deduplicate hosts")
	}

	hosts = HostSet{hosts: []*Host{}}
	hosts.Uniq()
}

func eq(a, b []*Host) bool {
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
