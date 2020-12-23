package katyusha

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Cache struct {
	name   string
	source HostProvider
	cache  *JsonProvider
	config struct {
		Lifetime time.Duration
		File     string
		Prefix   string
	}
}

func NewCache(name string) HostProvider {
	c := &Cache{name: name, cache: NewJsonProvider(name).(*JsonProvider)}
	c.config.File = name + ".cache"
	c.config.Lifetime = 1 * time.Hour
	return c
}

func NewCacheFromProvider(p HostProvider) HostProvider {
	c := NewCache(p.Name()).(*Cache)
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

func (c *Cache) Equivalent(p HostProvider) bool {
	panic("This should never be called for a cache")
}

func (c *Cache) ParseViper(v *viper.Viper) error {
	sv := v.Sub("Source")
	if sv == nil {
		return fmt.Errorf("No source specified")
	}
	s, err := NewProvider(sv.GetString("provider"), c.name)
	if err != nil {
		return err
	}
	s.ParseViper(sv)
	c.source = s
	return v.Unmarshal(&c.config)
}

func (c *Cache) mustRefresh() bool {
	info, err := os.Stat(c.config.File)
	return err != nil || time.Since(info.ModTime()) > c.config.Lifetime
}

func (c *Cache) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	if !c.mustRefresh() {
		logrus.Debugf("Loading cached data from %s for %s", c.config.File, c.source.Name())
		c.cache.config.File = c.config.File
		return c.cache.Load(ctx, mc)
	}
	mc <- CacheMessage{Name: c.name, Finished: false, Err: nil}
	hosts, err := c.source.Load(ctx, mc)
	mc <- CacheMessage{Name: c.name, Finished: true, Err: err}
	if len(hosts) > 0 {
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
