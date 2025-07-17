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
	testdata = filepath.Join(".", "provider", "plugin")
	if _, me, _, ok := runtime.Caller(0); ok {
		testdata = filepath.Join(filepath.Dir(me), "testdata")
	}
	if err := os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), filepath.Join(testdata, "bin")}, ":")); err != nil {
		panic(err)
	}
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
