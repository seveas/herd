package prometheus

import (
	"context"
	"github.com/spf13/viper"
	"testing"

	"github.com/seveas/herd"

	"github.com/jarcoal/httpmock"
)

func TestPrometheus(t *testing.T) {
	data := `{"status":"success","data":{"activeTargets":[{"discoveredLabels":{"__address__":"node1.herd.ci:9100","__metrics_path__":"/metrics","__scheme__":"http","job":"node"},"labels":{"instance":"node1.herd.ci:9100","job":"node"},"scrapePool":"node","scrapeUrl":"http://node1.herd.ci:9100/metrics","globalUrl":"http://node1.herd.ci:9100/metrics","lastError":"","lastScrape":"2020-12-23T19:31:54.709546474Z","lastScrapeDuration":0.019804372,"health":"up"}]}}`

	p := newPrometheusProvider("prometheus").(*prometheusProvider)
	p.config.Jobs = []string{"node"}
	v := viper.New()
	v.Set("url", "http://prometheus.herd.ci:9100/api/v1/targets")
	p.ParseViper(v)

	ctx := context.Background()
	mc := make(chan herd.CacheMessage)
	go func() {
		for {
			if _, ok := <-mc; !ok {
				break
			}
		}
	}()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://prometheus.herd.ci:9100/api/v1/targets",
		httpmock.NewStringResponder(200, data))
	hosts, err := p.Load(ctx, mc)
	if err != nil {
		t.Errorf("Failed to query mock consul: %s", err)
	}
	if len(hosts) != 1 {
		t.Errorf("Incorrect number of hosts returned (%d)", len(hosts))
	}
	if hosts[0].Attributes["job"] != "node" {
		t.Errorf("Job label not copied to host attributes")
	}
	if hosts[0].Attributes["health"] != "up" {
		t.Errorf("Health not copied to host attributes")
	}
}
