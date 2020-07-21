package katyusha

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-test/deep"
	"github.com/spf13/viper"
)

type HttpProvider struct {
	BaseProvider `mapstructure:",squash"`
	Url          string
	Username     string
	Password     string
	Headers      map[string]string
}

func NewHttpProvider(name string) HostProvider {
	return &HttpProvider{BaseProvider: BaseProvider{Name: name, Timeout: 30 * time.Second}}
}

func (p *HttpProvider) Equals(o HostProvider) bool {
	if c, ok := o.(*Cache); ok {
		o = c.Source
	}
	op, ok := o.(*HttpProvider)
	return ok &&
		p.BaseProvider.Equals(&op.BaseProvider) &&
		p.Url == op.Url &&
		p.Username == op.Username &&
		p.Password == op.Password &&
		deep.Equal(p.Headers, op.Headers) == nil
}

func (p *HttpProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(p)
}

func (p *HttpProvider) Fetch(ctx context.Context, mc chan CacheMessage) ([]byte, error) {
	req, err := http.NewRequest("GET", p.Url, nil)
	if err != nil {
		return []byte{}, err
	}
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
