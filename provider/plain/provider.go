package plain

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/seveas/katyusha"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("plain", newPlainTextProvider, plainTextMagic)
}

type plainTextProvider struct {
	name   string
	config struct {
		File   string
		Prefix string
	}
}

func newPlainTextProvider(name string) katyusha.HostProvider {
	return &plainTextProvider{name: name}
}

func plainTextMagic(r *katyusha.Registry) {
	p := &plainTextProvider{name: "inventory"}
	p.config.File = "inventory"
	r.AddMagicProvider(p)
}

func (p *plainTextProvider) Name() string {
	return p.name
}

func (p *plainTextProvider) Prefix() string {
	return p.config.Prefix
}

func (p *plainTextProvider) Equivalent(o katyusha.HostProvider) bool {
	return p.config.File == o.(*plainTextProvider).config.File
}

func (p *plainTextProvider) SetDataDir(dir string) error {
	if !filepath.IsAbs(p.config.File) {
		p.config.File = filepath.Join(dir, p.config.File)
		_, err := os.Stat(p.config.File)
		return err
	}
	return nil
}
func (p *plainTextProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *plainTextProvider) Load(ctx context.Context, mc chan katyusha.CacheMessage) (katyusha.Hosts, error) {
	hosts := make(katyusha.Hosts, 0)
	data, err := ioutil.ReadFile(p.config.File)
	if err != nil {
		logrus.Errorf("Could not load %s data in %s: %s", p.name, p.config.File, err)
		return hosts, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		host := katyusha.NewHost(line, katyusha.HostAttributes{})
		hosts = append(hosts, host)
	}
	return hosts, nil
}

var _ katyusha.DataLoader = &plainTextProvider{}
