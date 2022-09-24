package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("cache", newCache, nil)
}

type Cache struct {
	name   string
	source herd.HostProvider
	config struct {
		Lifetime time.Duration
		File     string
		Prefix   string
	}
}

func newCache(name string) herd.HostProvider {
	c := &Cache{name: name}
	c.config.File = name + ".cache"
	c.config.Lifetime = 1 * time.Hour
	return c
}

func NewFromProvider(p herd.HostProvider) herd.HostProvider {
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

func (c *Cache) Source() herd.HostProvider {
	return c.source
}

func (c *Cache) Invalidate() {
	c.config.Lifetime = -1
}

func (c *Cache) Equivalent(p herd.HostProvider) bool {
	panic("This should never be called for a cache")
}

func (c *Cache) ParseViper(v *viper.Viper) error {
	sv := v.Sub("Source")
	if sv == nil {
		return fmt.Errorf("No source specified")
	}
	s, err := herd.NewProvider(sv.GetString("provider"), c.name)
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

func (c *Cache) Load(ctx context.Context, lm herd.LoadingMessage) (herd.Hosts, error) {
	if !c.mustRefresh() {
		logrus.Debugf("Loading cached data from %s for %s", c.config.File, c.source.Name())
		hosts := make(herd.Hosts, 0)
		data, err := os.ReadFile(c.config.File)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &hosts); err != nil {
			return nil, err
		}
		return hosts, nil
	}
	hosts, err := c.source.Load(ctx, lm)
	if err == nil && len(hosts) > 0 {
		var data []byte
		dir := filepath.Dir(c.config.File)
		if err = os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("Unable to create cache directory %s: %s", dir, err.Error())
		}
		if data, err = json.Marshal(hosts); err != nil {
			return nil, err
		}
		//#nosec G306 -- Cache file may be shared among users
		if err := os.WriteFile(c.config.File, data, 0o644); err != nil {
			return nil, err
		}
	}
	return hosts, err
}

// Make sure we actually are a cache
var _ herd.Cache = &Cache{}
