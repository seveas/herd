package plugin

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/seveas/herd"
)

func init() {
	_, me, _, _ := runtime.Caller(0)
	testdata = filepath.Join(filepath.Dir(me), "testdata")
	os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), filepath.Join(testdata, "bin")}, ":"))
}

func TestPluginDiscovery(t *testing.T) {
	p, err := herd.NewProvider("fake", "fake")
	if err != nil {
		t.Errorf("Did not automatically find a plugin-based provider")
	}
	if pl, ok := p.(*pluginProvider); !ok {
		t.Errorf("Did not find a plugin, but another provider")
	} else if pl.config.Command != filepath.Join(testdata, "bin", "herd-provider-fake") {
		t.Errorf("Provider was not found by name automatically")
	}
}
