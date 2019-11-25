package herd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Cache struct {
	Lifetime time.Duration
	File     string
	Provider HostProvider
}

func (c *Cache) MustRefresh() bool {
	info, err := os.Stat(c.File)
	return err != nil || time.Since(info.ModTime()) > c.Lifetime
}

func (c *Cache) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	if !c.MustRefresh() {
		jp := &JsonProvider{Name: c.Provider.String(), File: c.File}
		hosts, err := jp.Load(ctx, mc)
		if err != nil {
			return hosts, err
		}
		if p, ok := c.Provider.(PostProcessor); ok {
			p.PostProcess(hosts)
		}
		return hosts, err
	}
	mc <- CacheMessage{name: c.Provider.String(), finished: false, err: nil}
	hosts, err := c.Provider.Load(ctx, mc)
	mc <- CacheMessage{name: c.Provider.String(), finished: true, err: err}
	if len(hosts) > 0 {
		var data []byte
		if err = os.MkdirAll(viper.GetString("CacheDir"), 0700); err != nil {
			return hosts, fmt.Errorf("Unable to create cache: %s", err.Error())
		}
		data, err = json.Marshal(hosts)
		if err == nil {
			err = ioutil.WriteFile(c.File, data, 0644)
		}
	}
	return hosts, err
}

func (c *Cache) String() string {
	return fmt.Sprintf("%s (cached)", c.Provider.String())
}
