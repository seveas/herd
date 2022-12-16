package plugin

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var testdata string

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(io.Discard)
	testdata = filepath.Join(".", "provider", "plugin")
	if _, me, _, ok := runtime.Caller(0); ok {
		testdata = filepath.Join(filepath.Dir(me), "testdata")
	}
	os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), filepath.Join(testdata, "bin")}, ":"))
}

func TestFindingProvider(t *testing.T) {
	p := newPlugin("fake").(*pluginProvider)
	if p.config.Command != filepath.Join(testdata, "bin", "herd-provider-fake") {
		t.Errorf("Provider was not found by name automatically")
	}
	p = newPlugin("fake2").(*pluginProvider)
	if p.config.Command != filepath.Join(testdata, "bin", "herd-provider-fake2.exe") {
		t.Errorf("Provider with .exe suffix was not found by name automatically")
	}
}

func TestPluginConnection(t *testing.T) {
	testcases := []struct {
		mode  string
		err   string
		count int
		log   bool
	}{
		{"normal", "<nil>", 5, true},
		{"config-error", "Simulated configuration error", 0, false},
		{"config-panic", "rpc error: code = Unavailable desc = error reading from server: EOF", 0, false},
		{"empty", "<nil>", 0, true},
		{"error", "Simulated load error", 0, true},
		{"panic", "rpc error: code = Unavailable desc = error reading from server: EOF", 0, true},
	}
	for _, test := range testcases {
		t.Run(test.mode, func(t *testing.T) {
			msg := struct {
				name string
				done bool
				err  error
			}{}
			p := newPlugin("ci").(*pluginProvider)
			hook := &logrusHook{seen: make(map[logrus.Level]bool)}
			logrus.AddHook(hook)
			if p.config.Command != filepath.Join(testdata, "bin", "herd-provider-ci") {
				t.Errorf("Provider was not found by name automatically")
			}
			v := viper.New()
			v.Set("Mode", test.mode)
			err := p.ParseViper(v)
			if err != nil {
				if test.mode == "config-error" || test.mode == "config-panic" {
					if err.Error() != test.err {
						t.Errorf("Unexpected configuration error: %s", err)
					}
				} else {
					t.Errorf("Unable to configure plugin: %s", err)
				}
				return
			}
			hosts, err := p.Load(context.Background(), func(name string, done bool, err error) { msg.name = name })
			if fmt.Sprintf("%v", err) != test.err {
				t.Errorf("Unexpected load error: %v. Expected: %v", err, test.err)
				return
			}
			if hosts != nil && hosts.Len() != test.count {
				t.Errorf("Received %d hosts, expecting %d", hosts.Len(), test.count)
			}
			if test.log {
				if msg.name != "ci" {
					t.Errorf("No loading message was received")
				}
				for _, level := range logrus.AllLevels {
					if _, ok := hook.seen[level]; !ok && level > logrus.FatalLevel {
						t.Errorf("No %s message was received", level)
					}
				}
			}
		})
	}
}

type logrusHook struct {
	seen map[logrus.Level]bool
}

func (h *logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logrusHook) Fire(entry *logrus.Entry) error {
	h.seen[entry.Level] = true
	if entry.Level <= logrus.FatalLevel {
		entry.Level = logrus.ErrorLevel
	}
	return nil
}

func TestDataDirProvider(t *testing.T) {
	p := newPlugin("ci_dataloader").(*pluginProvider)
	v := viper.New()
	if err := p.ParseViper(v); err != nil {
		t.Errorf("Unable to configure plugin: %s", err)
	}
	if err := p.SetDataDir("testdata"); err != nil {
		t.Errorf("Unable to set data dir: %s", err)
		return
	}
	hosts, err := p.Load(context.Background(), func(name string, done bool, err error) {})
	if err != nil {
		t.Errorf("Unable to load hosts: %s", err)
		return
	}
	if hosts.Len() != 1 {
		t.Errorf("Received %d hosts, expecting %d", hosts.Len(), 1)
	}
	if dir, ok := hosts.Get(0).Attributes["datadir"]; !ok || dir != "testdata" {
		t.Errorf("Data dir was not set correctly, got %s", dir)
	}
}

func TestCache(t *testing.T) {
	p := newPlugin("ci_cache").(*pluginProvider)
	v := viper.New()
	if err := p.ParseViper(v); err != nil {
		t.Errorf("Unable to configure plugin: %s", err)
	}
	p.SetCacheDir("testcache")
	p.Keep()
	p.Invalidate()
	hosts, err := p.Load(context.Background(), func(name string, done bool, err error) {})
	if err != nil {
		t.Errorf("Unable to load hosts: %s", err)
		return
	}
	if hosts.Len() != 1 {
		t.Errorf("Received %d hosts, expecting %d", hosts.Len(), 1)
		return
	}
	attrs := hosts.Get(0).Attributes
	if dir, ok := attrs["cachedir"]; !ok || dir != "testcache" {
		t.Errorf("Cache dir was not set correctly, got %s", dir)
	}
	if keep, ok := attrs["keep"]; !ok || keep != true {
		t.Errorf("Keep was not called")
	}
	if invalidate, ok := attrs["invalidate"]; !ok || invalidate != true {
		t.Errorf("Invalidate was not called")
	}
}
