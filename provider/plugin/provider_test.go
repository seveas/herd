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
	_, me, _, _ := runtime.Caller(0)
	testdata = filepath.Join(filepath.Dir(me), "testdata")
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
				t.Errorf("Unable to configure plugin: %s", err)
			}
			hosts, err := p.Load(context.Background(), func(name string, done bool, err error) { msg.name = name })
			if fmt.Sprintf("%v", err) != test.err {
				t.Errorf("Unexpected load error: %v. Expected: %v", err, test.err)
			}
			if len(hosts) != test.count {
				t.Errorf("Received %d hosts, expecting %d", len(hosts), test.count)
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
		fmt.Println("NOOO")
		entry.Level = logrus.ErrorLevel
	}
	return nil
}
