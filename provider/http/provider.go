package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/seveas/katyusha"

	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("http", NewProvider, nil)
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

func NewProvider(name string) katyusha.HostProvider {
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

func (p *HttpProvider) Equivalent(o katyusha.HostProvider) bool {
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
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("http response code %d: %s", resp.StatusCode, body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (p *HttpProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (katyusha.Hosts, error) {
	hosts := katyusha.Hosts{}
	lm(p.name, false, nil)
	data, err := p.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}
