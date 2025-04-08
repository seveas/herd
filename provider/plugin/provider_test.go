package plugin

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	p := newPlugin("ci").(*pluginProvider)
	v := viper.New()
	v.Set("Mode", "normal")

	err := p.ParseViper(v)
	if err != nil {
		t.Errorf("Unable to configure plugin: %s", err)
	}
	hosts, err := p.Load(t.Context(), func(name string, done bool, err error) {})
	if err != nil {
		t.Errorf("Unable to load hosts: %s", err)
	}
	if hosts.Len() != 5 {
		t.Errorf("Received %d hosts, expecting %d", hosts.Len(), 5)
	}
}

func TestPluginConfigError(t *testing.T) {
	p := newPlugin("ci").(*pluginProvider)
	v := viper.New()
	v.Set("Mode", "config-error")

	err := p.ParseViper(v)
	if err.Error() != "Simulated configuration error" {
		t.Errorf("Unexpected configuration error: %s", err)
	}
}

func TestPluginConfigPanic(t *testing.T) {
	p := newPlugin("ci").(*pluginProvider)
	v := viper.New()
	v.Set("Mode", "config-panic")

	err := p.ParseViper(v)
	if status.Code(err) == codes.Unknown {
		t.Errorf("Unexpected rpc/configuration error: %s (%d)", err, status.Code(err))
	}
}

func TestPluginEmpty(t *testing.T) {
	p := newPlugin("ci").(*pluginProvider)
	v := viper.New()
	v.Set("Mode", "empty")

	err := p.ParseViper(v)
	if err != nil {
		t.Errorf("Unable to configure plugin: %s", err)
	}

	hosts, err := p.Load(t.Context(), func(name string, done bool, err error) {})
	if err != nil {
		t.Errorf("Unable to load hosts: %s", err)
	}
	if hosts != nil {
		t.Errorf("Received hosts when expecting nil")
	}
}

func TestPluginError(t *testing.T) {
	p := newPlugin("ci").(*pluginProvider)
	v := viper.New()
	v.Set("Mode", "error")

	err := p.ParseViper(v)
	if err != nil {
		t.Errorf("Unable to configure plugin: %s", err)
	}

	msg := struct {
		name string
		done bool
		err  error
	}{}
	hook := &logrusHook{seen: make(map[logrus.Level]bool)}
	logrus.AddHook(hook)
	_, err = p.Load(t.Context(), func(name string, done bool, err error) { msg.name = name })
	if err.Error() != "Simulated load error" {
		t.Errorf("Unexpected load error: %v. Expected: %v", err, "Simulated load error")
	}

	if msg.name != "ci" {
		t.Errorf("No loading message was received")
	}
	for _, level := range logrus.AllLevels {
		if _, ok := hook.seen[level]; !ok && level > logrus.FatalLevel {
			t.Errorf("No %s message was received", level)
		}
	}
}

func TestPluginPanic(t *testing.T) {
	p := newPlugin("ci").(*pluginProvider)
	v := viper.New()
	v.Set("Mode", "panic")

	err := p.ParseViper(v)
	if err != nil {
		t.Errorf("Unable to configure plugin: %s", err)
	}

	_, err = p.Load(t.Context(), func(name string, done bool, err error) {})
	if status.Code(err) == codes.Unknown {
		t.Errorf("Unexpected rpc/configuration error: %s (%d)", err, status.Code(err))
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
	hosts, err := p.Load(t.Context(), func(name string, done bool, err error) {})
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
	hosts, err := p.Load(t.Context(), func(name string, done bool, err error) {})
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
