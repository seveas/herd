package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/seveas/katyusha"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("cache", newCache, nil)
}

type Cache struct {
	name   string
	source katyusha.HostProvider
	config struct {
		Lifetime time.Duration
		File     string
		Prefix   string
	}
}

func newCache(name string) katyusha.HostProvider {
	c := &Cache{name: name}
	c.config.File = name + ".cache"
	c.config.Lifetime = 1 * time.Hour
	return c
}

func NewFromProvider(p katyusha.HostProvider) katyusha.HostProvider {
	c := newCache(p.Name()).(*Cache)
	c.source = p
	return c
}

func (c *Cache) Name() string {
	return c.name
}

func (c *Cache) Prefix() string {
	return c.config.Prefix + c.source.Prefix()
}

func (c *Cache) SetCacheDir(dir string) {
	if !filepath.IsAbs(c.config.File) {
		c.config.File = filepath.Join(dir, c.config.File)
	}
}

func (c *Cache) Source() katyusha.HostProvider {
	return c.source
}

func (c *Cache) Invalidate() {
	c.config.Lifetime = -1
}

func (c *Cache) Equivalent(p katyusha.HostProvider) bool {
	panic("This should never be called for a cache")
}

func (c *Cache) ParseViper(v *viper.Viper) error {
	sv := v.Sub("Source")
	if sv == nil {
		return fmt.Errorf("No source specified")
	}
	s, err := katyusha.NewProvider(sv.GetString("provider"), c.name)
	if err != nil {
		return err
	}
	if err := s.ParseViper(sv); err != nil {
		return err
	}
	c.source = s
	return v.Unmarshal(&c.config)
}

func (c *Cache) mustRefresh() bool {
	info, err := os.Stat(c.config.File)
	return err != nil || time.Since(info.ModTime()) > c.config.Lifetime
}

func (c *Cache) Load(ctx context.Context, lm katyusha.LoadingMessage) (katyusha.Hosts, error) {
	if !c.mustRefresh() {
		logrus.Debugf("Loading cached data from %s for %s", c.config.File, c.source.Name())
		hosts := make(katyusha.Hosts, 0)
		data, err := ioutil.ReadFile(c.config.File)
		if err != nil {
			return hosts, err
		}

		err = json.Unmarshal(data, &hosts)
		return hosts, err
	}
	hosts, err := c.source.Load(ctx, lm)
	if err == nil && len(hosts) > 0 {
		var data []byte
		dir := filepath.Dir(c.config.File)
		if err = os.MkdirAll(dir, 0700); err != nil {
			return hosts, fmt.Errorf("Unable to create cache directory %s: %s", dir, err.Error())
		}
		data, err = json.Marshal(hosts)
		if err == nil {
			err = ioutil.WriteFile(c.config.File, data, 0644)
		}
	}
	return hosts, err
}

// Make sure we actually are a cache
var _ katyusha.Cache = &Cache{}
