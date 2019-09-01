package katyusha

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

func (p *HttpProvider) String() string {
	return p.Name
}

func (p *HttpProvider) Cache(ctx context.Context) error {
	if info, err := os.Stat(p.File); err == nil && time.Since(info.ModTime()) < p.CacheLifetime {
		return nil
	}
	UI.Infof("Refreshing %s cache", p.Name)

	req, err := http.NewRequest("GET", p.Url, nil)
	if err != nil {
		return err
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
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("http response code %d: %s", resp.StatusCode, body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(p.File+".new", body, 0600); err != nil {
		return err
	}
	if err := os.Rename(p.File+".new", p.File); err != nil {
		return err
	}
	return nil
}

func (p *HttpProvider) GetHosts(hostnameGlob string) Hosts {
	jp := &JsonProvider{Name: p.Name, File: p.File}
	return jp.GetHosts(hostnameGlob)
}
