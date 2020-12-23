package katyusha

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var testDataRoot string
var realUserHome string

func homeDir(h string) string {
	return filepath.Join(testDataRoot, "homes", h)
}

func dataDir(h string) string {
	return filepath.Join(homeDir(h), ".local", "share", "katyusha")
}

func cacheDir(h string) string {
	return filepath.Join(homeDir(h), ".cache", "katyusha")
}

func TestMain(m *testing.M) {
	_, me, _, _ := runtime.Caller(0)
	testDataRoot = filepath.Join(filepath.Dir(me), "testdata")
	// Clean the environment
	if home, ok := os.LookupEnv("HOME"); ok {
		realUserHome = home
	} else {
		u, err := user.Current()
		if err == nil && u.HomeDir != "" {
			realUserHome = u.HomeDir
		}
	}
	for _, envvar := range os.Environ() {
		if strings.HasPrefix(envvar, "KATYUSHA") {
			os.Unsetenv(envvar[:strings.IndexRune(envvar, '=')])
		}
		if strings.HasPrefix(envvar, "AWS") {
			os.Unsetenv(envvar[:strings.IndexRune(envvar, '=')])
		}
	}
	os.Exit(m.Run())
}
