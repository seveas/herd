// +build !no_prometheus

package katyusha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	availableProviders["prometheus"] = NewPrometheusProvider
}

type PrometheusProvider struct {
	name   string
	hp     *HttpProvider
	config struct {
		Jobs []string
	}
}

type PrometheusTargets struct {
	Status string                        `json:"status"`
	Data   map[string][]PrometheusTarget `json:"data"`
}

type PrometheusTarget struct {
	Labels             map[string]string `json:"labels"`
	ScrapePool         string            `json:"scrapePool"`
	ScrapeUrl          string            `json:"scrapeUrl"`
	LastError          string            `json:"lastError"`
	LastScrape         time.Time         `json:"lastScrape"`
	LastScrapeDuration float64           `json:"lastScrapeDuration"`
	Health             string            `json:"health"`
}

func NewPrometheusProvider(name string) HostProvider {
	return &PrometheusProvider{name: name, hp: NewHttpProvider(name).(*HttpProvider)}
}

func (p *PrometheusProvider) Name() string {
	return p.name
}

func (p *PrometheusProvider) Prefix() string {
	return p.hp.Prefix()
}

func (p *PrometheusProvider) Equivalent(o HostProvider) bool {
	op := o.(*PrometheusProvider)
	return p.hp.Equivalent(op.hp) &&
		reflect.DeepEqual(p.config.Jobs, op.config.Jobs)
}

func (p *PrometheusProvider) ParseViper(v *viper.Viper) error {
	if err := v.Unmarshal(&p.hp.config); err != nil {
		return err
	}
	return v.Unmarshal(&p.config)
}

func (p *PrometheusProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	data, err := p.hp.Fetch(ctx, mc)
	if err != nil {
		return Hosts{}, err
	}
	var targets PrometheusTargets
	err = json.Unmarshal(data, &targets)
	if err != nil {
		return Hosts{}, err
	}
	if targets.Status != "success" {
		return Hosts{}, fmt.Errorf("Prometheus API returned: %s", targets.Status)
	}

	ret := make(Hosts, 0)
	for _, target := range targets.Data["activeTargets"] {
		job := target.Labels["job"]
		found := false
		for _, j := range p.config.Jobs {
			if j == job {
				found = true
			}
		}
		if !found {
			continue
		}
		u, err := url.Parse(target.ScrapeUrl)
		if err != nil {
			logrus.Warnf("Unable to parse scrape URL: %s", target.ScrapeUrl)
			continue
		}
		name := u.Host
		if instance, ok := target.Labels["instance"]; ok {
			if strings.Contains(instance, "://") {
				u, err := url.Parse(instance)
				if err != nil {
					logrus.Warnf("Unable to parse scrape instance: %s", instance)
					continue
				}
				name = u.Host
			} else {
				parts := strings.Split(instance, ":")
				name = parts[0]
			}
		}
		attributes := HostAttributes{
			"scrape_pool":          target.ScrapePool,
			"scrape_url":           target.ScrapeUrl,
			"last_scrape":          target.LastScrape,
			"last_scrape_duration": target.LastScrapeDuration,
			"health":               target.Health,
		}
		for k, v := range target.Labels {
			attributes[k] = v
		}
		ret = append(ret, NewHost(name, attributes))
	}

	return ret, nil
}
