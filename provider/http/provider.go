package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

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
	}
}

func NewProvider(name string) herd.HostProvider {
	return &HttpProvider{name: name, client: http.DefaultClient}
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
	req, err := http.NewRequest(http.MethodGet, p.config.Url, nil)
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
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", fmt.Sprintf("herd/%s", herd.Version()))
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("http response code %d: %s", resp.StatusCode, body)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (p *HttpProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	hosts := new(herd.HostSet)
	lm(p.name, false, nil)
	data, err := p.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}
