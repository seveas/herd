package katyusha

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Cache struct {
	Name     string
	Lifetime time.Duration
	File     string
	Source   HostProvider
}

func NewCache(name string) HostProvider {
	return &Cache{Lifetime: 1 * time.Hour, Name: name}
}

func (c *Cache) ParseViper(v *viper.Viper) error {
	sv := v.Sub("Source")
	if sv == nil {
		return fmt.Errorf("No source specified")
	}
	s, err := NewProvider(sv.GetString("provider"), c.Name)
	if err != nil {
		return err
	}
	s.ParseViper(sv)
	v.Set("Source", s)
	return v.Unmarshal(c)
}

func (c *Cache) mustRefresh() bool {
	info, err := os.Stat(c.File)
	return err != nil || time.Since(info.ModTime()) > c.Lifetime
}

func (c *Cache) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	if !c.mustRefresh() {
		jp := &JsonProvider{Name: c.Source.String(), File: c.File}
		hosts, err := jp.Load(ctx, mc)
		if err != nil {
			return hosts, err
		}
		return hosts, err
	}
	mc <- CacheMessage{Name: c.Source.String(), Finished: false, Err: nil}
	hosts, err := c.Source.Load(ctx, mc)
	mc <- CacheMessage{Name: c.Source.String(), Finished: true, Err: err}
	if len(hosts) > 0 {
		var data []byte
		dir := filepath.Dir(c.File)
		if err = os.MkdirAll(dir, 0700); err != nil {
			return hosts, fmt.Errorf("Unable to create cache directory %s: %s", dir, err.Error())
		}
		data, err = json.Marshal(hosts)
		if err == nil {
			err = ioutil.WriteFile(c.File, data, 0644)
		}
	}
	return hosts, err
}

func (c *Cache) String() string {
	return fmt.Sprintf("%s (cached)", c.Source.String())
}
