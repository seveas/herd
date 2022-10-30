package known_hosts

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
)

var testdata string

func init() {
	testdata = filepath.Join(".", "provider", "known_hosts")
	if _, me, _, ok := runtime.Caller(0); ok {
		testdata = filepath.Join(filepath.Dir(me), "testdata")
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		path  string
		hosts int
	}{
		{"normal", 2},
		{"hashed", 0},
		{"mixed", 2},
		{"malformed_start", 0},
		{"malformed_end", 2},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			t.Setenv("HOME", filepath.Join(testdata, test.path))
			p := magicProvider().(*knownHostsProvider)
			p.config.Files = p.config.Files[1:]
			hosts, err := p.Load(context.Background(), nil)
			if err != nil {
				t.Errorf("Error parsing known_hosts: %s", err)
			}
			if hosts.Len() != test.hosts {
				t.Errorf("Incorrect number of hosts (%d) returned, expected %d", hosts.Len(), test.hosts)
				return
			}
			if test.hosts == 0 {
				return
			}
			expect := 3
			if len(hosts.Get(0).PublicKeys()) != expect {
				t.Errorf("Incorrect number of keys (%d) returned for the first host, expected %d", len(hosts.Get(0).PublicKeys()), expect)
			}
		})
	}
}
