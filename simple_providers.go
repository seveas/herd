package katyusha

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	providerMagic["plain"] = func(dataDir string) []HostProvider {
		fn := path.Join(dataDir, "inventory")
		if _, err := os.Stat(fn); err != nil {
			return []HostProvider{}
		}
		return []HostProvider{&PlainTextProvider{Name: "inventory", File: fn}}
	}
	providerMakers["plain"] = func(dataDir, name string, v *viper.Viper) (HostProvider, error) {
		return &PlainTextProvider{Name: name, File: v.GetString("File")}, nil
	}
}

type PlainTextProvider struct {
	Name string
	File string
}

func (p *PlainTextProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := make(Hosts, 0)
	data, err := ioutil.ReadFile(p.File)
	if err != nil {
		logrus.Errorf("Could not load %s data in %s: %s", p.Name, p.File, err)
		return hosts, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		host := NewHost(line, HostAttributes{})
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func (p *PlainTextProvider) String() string {
	return p.Name
}
