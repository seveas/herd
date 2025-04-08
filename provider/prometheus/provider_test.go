package prometheus

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

func TestPrometheus(t *testing.T) {
	data := `{"status":"success","data":{"activeTargets":[{"discoveredLabels":{"__address__":"node1.herd.ci:9100","__metrics_path__":"/metrics","__scheme__":"http","job":"node"},"labels":{"instance":"node1.herd.ci:9100","job":"node"},"scrapePool":"node","scrapeUrl":"http://node1.herd.ci:9100/metrics","globalUrl":"http://node1.herd.ci:9100/metrics","lastError":"","lastScrape":"2020-12-23T19:31:54.709546474Z","lastScrapeDuration":0.019804372,"health":"up"}]}}`

	p := newProvider("prometheus").(*prometheusProvider)
	p.config.Jobs = []string{"node"}
	v := viper.New()
	v.Set("url", "http://prometheus.herd.ci:9100/api/v1/targets")
	if err := p.ParseViper(v); err != nil {
		t.Errorf("ParseViper failed: %s", err)
	}

	ctx := t.Context()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://prometheus.herd.ci:9100/api/v1/targets",
		httpmock.NewStringResponder(200, data))
	hosts, err := p.Load(ctx, func(string, bool, error) {})
	if err != nil {
		t.Errorf("Failed to query mock consul: %s", err)
	}
	if hosts.Len() != 1 {
		t.Errorf("Incorrect number of hosts returned (%d)", hosts.Len())
	}
	h := hosts.Get(0)
	if h.Attributes["job"] != "node" {
		t.Errorf("Job label not copied to host attributes")
	}
	if h.Attributes["health"] != "up" {
		t.Errorf("Health not copied to host attributes")
	}
}
