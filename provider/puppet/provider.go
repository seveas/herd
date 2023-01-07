package puppet

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func init() {
	herd.RegisterProvider("puppet", newProvider, nil)
}

type fact struct {
	Certname string `json:"certname"`
	Name     string `json:"name"`
	Value    any    `json:"value"`
}

type puppetProvider struct {
	name   string
	hp     *http.HttpProvider
	config struct {
		Facts []string
	}
}

func newProvider(name string) herd.HostProvider {
	p := &puppetProvider{name: name, hp: http.NewProvider(name).(*http.HttpProvider)}
	p.config.Facts = []string{"ssh", "os", "dmi"}
	return p
}

func (p *puppetProvider) Name() string {
	return p.name
}

func (p *puppetProvider) Prefix() string {
	return p.hp.Prefix()
}

func (p *puppetProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*puppetProvider)
	return p.hp.Equivalent(op.hp) &&
		reflect.DeepEqual(p.config, op.config)
}

func (p *puppetProvider) ParseViper(v *viper.Viper) error {
	if err := v.Unmarshal(&p.config); err != nil {
		return err
	}
	u, err := url.Parse(v.GetString("Url"))
	if err != nil {
		return err
	}
	u.Path = "/pdb/query/v4"
	query := fmt.Sprintf("facts[certname,name,value]{name in [%s]}", strings.Join(quote(p.config.Facts), ","))
	u.RawQuery = url.Values{"query": []string{query}}.Encode()

	v.Set("Url", u.String())
	if err := p.hp.ParseViper(v); err != nil {
		return err
	}
	return nil
}

func (p *puppetProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	lm(p.name, false, nil)
	data, err := p.hp.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	var facts []fact
	err = json.Unmarshal(data, &facts)
	if err != nil {
		return nil, err
	}
	hosts := make(map[string]*herd.Host)
	for _, f := range facts {
		host, ok := hosts[f.Certname]
		if !ok {
			host = herd.NewHost(f.Certname, "", herd.HostAttributes{})
			hosts[host.Name] = host
		}
		if f.Name == "ssh" {
			addHostKeys(host, f.Value.(map[string]any))
		} else {
			addFact(host, f.Name, f.Value, "")
		}
	}

	ret := herd.NewHostSet()
	for _, host := range hosts {
		ret.AddHost(host)
	}
	return ret, nil
}

func addFact(host *herd.Host, name string, value any, prefix string) {
	switch v := value.(type) {
	case string:
		host.Attributes[prefix+name] = v
	case map[string]interface{}:
		for k, v := range v {
			addFact(host, k, v, prefix+name+":")
		}
	default:
		host.Attributes[prefix+name] = value
	}
}

func addHostKeys(host *herd.Host, keys map[string]any) {
	for _, k := range keys {
		kd := k.(map[string]any)
		b, err := base64.StdEncoding.DecodeString(kd["key"].(string))
		if err != nil {
			logrus.Warnf("Error parsing %s public keys: %s", kd["type"], err)
			continue
		}
		key, err := ssh.ParsePublicKey(b)
		if err != nil {
			logrus.Warnf("Error parsing %s public keys: %s", kd["type"], err)
			continue
		}
		host.AddPublicKey(key)
	}
}

func quote(in []string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = fmt.Sprintf("%q", s)
	}
	return out
}
