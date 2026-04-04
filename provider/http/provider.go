package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/seveas/herd"

	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("http", NewProvider, nil)
}

type HttpProvider struct {
	name     string
	client   *http.Client
	proxyUrl *url.URL
	config   struct {
		Prefix   string
		Url      string
		Username string
		Password string // #nosec G117 -- Credential field required for HTTP basic auth
		Headers  map[string]string
		Proxy    string
	}
}

func NewProvider(name string) herd.HostProvider {
	p := &HttpProvider{name: name, client: &http.Client{}}
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
	if err := v.Unmarshal(&p.config); err != nil {
		return err
	}
	if p.config.Proxy != "" {
		proxyUrl, err := url.Parse(p.config.Proxy)
		if err != nil {
			return fmt.Errorf("invalid proxy URL: %w", err)
		}
		p.proxyUrl = proxyUrl
		p.client.Transport = &http.Transport{
			Proxy: p.proxy,
		}
	}
	return nil
}

func (p *HttpProvider) proxy(req *http.Request) (*url.URL, error) {
	if p.proxyUrl == nil {
		return http.ProxyFromEnvironment(req)
	}
	return p.proxyUrl, nil
}

func (p *HttpProvider) Fetch(ctx context.Context) ([]byte, error) {
	return p.FetchUrl(ctx, p.config.Url)
}

func (p *HttpProvider) FetchUrl(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	resp, err := p.client.Do(req) // #nosec G704 -- URL is from configuration, not user input
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
