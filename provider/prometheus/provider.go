package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("prometheus", newPrometheusProvider, nil)
}

type prometheusProvider struct {
	name   string
	hp     *herd.HttpProvider
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

func newPrometheusProvider(name string) herd.HostProvider {
	return &prometheusProvider{name: name, hp: herd.NewHttpProvider(name).(*herd.HttpProvider)}
}

func (p *prometheusProvider) Name() string {
	return p.name
}

func (p *prometheusProvider) Prefix() string {
	return p.hp.Prefix()
}

func (p *prometheusProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*prometheusProvider)
	return p.hp.Equivalent(op.hp) &&
		reflect.DeepEqual(p.config.Jobs, op.config.Jobs)
}

func (p *prometheusProvider) ParseViper(v *viper.Viper) error {
	if err := p.hp.ParseViper(v); err != nil {
		return err
	}
	return v.Unmarshal(&p.config)
}

func (p *prometheusProvider) Load(ctx context.Context, mc chan herd.CacheMessage) (herd.Hosts, error) {
	data, err := p.hp.Fetch(ctx, mc)
	if err != nil {
		return herd.Hosts{}, err
	}
	var targets PrometheusTargets
	err = json.Unmarshal(data, &targets)
	if err != nil {
		return herd.Hosts{}, err
	}
	if targets.Status != "success" {
		return herd.Hosts{}, fmt.Errorf("Prometheus API returned: %s", targets.Status)
	}

	ret := make(herd.Hosts, 0)
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
		attributes := herd.HostAttributes{
			"scrape_pool":          target.ScrapePool,
			"scrape_url":           target.ScrapeUrl,
			"last_scrape":          target.LastScrape,
			"last_scrape_duration": target.LastScrapeDuration,
			"health":               target.Health,
		}
		for k, v := range target.Labels {
			attributes[k] = v
		}
		ret = append(ret, herd.NewHost(name, attributes))
	}

	return ret, nil
}
