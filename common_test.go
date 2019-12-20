package herd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

var testDataRoot string
var realUserHome string

func TestMain(m *testing.M) {
	_, me, _, _ := runtime.Caller(0)
	testDataRoot = filepath.Join(filepath.Dir(me), "testdata")
	// Clean the environment
	homedir.DisableCache = true
	realUserHome, _ = homedir.Dir()
	os.Unsetenv("CONSUL_HTTP_ADDR")
	for _, envvar := range os.Environ() {
		if strings.HasPrefix(envvar, "HERD") {
			os.Unsetenv(envvar[:strings.IndexRune(envvar, '=')])
		}
	}
	os.Exit(m.Run())
}
