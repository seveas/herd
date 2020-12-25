package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("http", NewProvider, nil)
}

type HttpProvider struct {
	name   string
	client *http.Client
	config struct {
		Prefix   string
		Url      string
		Username string
		Password string
		Headers  map[string]string
		Timeout  time.Duration
	}
}

func NewProvider(name string) herd.HostProvider {
	p := &HttpProvider{name: name, client: http.DefaultClient}
	p.config.Timeout = 5 * time.Second
	return p
}

func (p *HttpProvider) Name() string {
	return p.name
}

func (p *HttpProvider) Prefix() string {
	return p.config.Prefix
}

func (p *HttpProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*HttpProvider)
	return p.config.Url == op.config.Url &&
		p.config.Username == op.config.Username &&
		p.config.Password == op.config.Password &&
		reflect.DeepEqual(p.config.Headers, op.config.Headers)
}

func (p *HttpProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *HttpProvider) Fetch(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()
	req, err := http.NewRequest("GET", p.config.Url, nil)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	if p.config.Username != "" {
		req.SetBasicAuth(p.config.Username, p.config.Password)
	}
	if p.config.Headers != nil {
		for key, value := range p.config.Headers {
			req.Header.Set(key, value)
		}
	}
	resp, err := p.client.Do(req)
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

func (p *HttpProvider) Load(ctx context.Context, lm herd.LoadingMessage) (herd.Hosts, error) {
	hosts := herd.Hosts{}
	lm(p.name, false, nil)
	data, err := p.Fetch(ctx)
	if err != nil {
		return hosts, err
	}
	if err = json.Unmarshal(data, &hosts); err != nil {
		err = fmt.Errorf("Could not parse %s data from %s: %s", p.name, p.config.Url, err)
	}
	return hosts, err
}
