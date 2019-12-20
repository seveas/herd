package herd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/spf13/viper"
)

type HttpProvider struct {
	Name          string
	File          string
	Url           string
	Username      string
	Password      string
	Headers       map[string]string
	Timeout       time.Duration
	CacheLifetime time.Duration
}

func init() {
	providerMakers["http"] = func(dataDir, name string, v *viper.Viper) (HostProvider, error) {
		p := &HttpProvider{
			Name:          name,
			File:          path.Join(dataDir, "cache", name+".cache"),
			CacheLifetime: 1 * time.Hour,
			Timeout:       30 * time.Second,
		}
		err := v.Unmarshal(p)
		if err != nil {
			return nil, err
		}
		return &Cache{File: p.File, Lifetime: p.CacheLifetime, Provider: p}, nil
	}
}

func (p *HttpProvider) String() string {
	return p.Name
}

func (p *HttpProvider) Fetch(ctx context.Context, mc chan CacheMessage) ([]byte, error) {
	req, err := http.NewRequest("GET", p.Url, nil)
	if err != nil {
		return []byte{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()
	req = req.WithContext(ctx)
	if p.Username != "" {
		req.SetBasicAuth(p.Username, p.Password)
	}
	if p.Headers != nil {
		for key, value := range p.Headers {
			req.Header.Set(key, value)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return []byte{}, fmt.Errorf("http response code %d: %s", resp.StatusCode, body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, err
}

func (p *HttpProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	hosts := Hosts{}
	data, err := p.Fetch(ctx, mc)
	if err != nil {
		return hosts, err
	}
	if err = json.Unmarshal(data, &hosts); err != nil {
		err = fmt.Errorf("Could not parse %s data from %s: %s", p.Name, p.Url, err)
	}
	return hosts, err
}
