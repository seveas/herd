package known_hosts

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var testdata string

func init() {
	_, me, _, _ := runtime.Caller(0)
	testdata = filepath.Join(filepath.Dir(me), "testdata")
}

func TestParser(t *testing.T) {
	if oldHome, ok := os.LookupEnv("HOME"); ok {
		defer os.Setenv("HOME", oldHome)
	}
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
			os.Setenv("HOME", filepath.Join(testdata, test.path))
			p := magicProvider().(*knownHostsProvider)
			p.config.Files = p.config.Files[1:]
			hosts, err := p.Load(nil, nil)
			if err != nil {
				t.Errorf("Error parsing known_hosts: %s", err)
			}
			if len(hosts) != test.hosts {
				t.Errorf("Incorrect number of hosts (%d) returned, expected %d", len(hosts), test.hosts)
				return
			}
			if test.hosts == 0 {
				return
			}
			expect := 3
			if len(hosts[0].PublicKeys()) != expect {
				t.Errorf("Incorrect number of keys (%d) returned for the first host, expected %d", len(hosts[0].PublicKeys()), expect)
			}
		})
	}
}
