package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

func TestProviderEquivalence(t *testing.T) {
	p1 := newProvider("test").(*consulProvider)
	p1.config.Address = "http://consul:8080"

	p2 := newProvider("test 2").(*consulProvider)
	p2.config.Address = "http://consul:8080"
	p2.config.Prefix = "consul:"

	if !p1.Equivalent(p2) {
		t.Errorf("Equivalence not properly detected")
	}

	p2.config.Address = "http://consul:8081"
	if p1.Equivalent(p2) {
		t.Errorf("Non-equivalence not properly detected")
	}
}

func TestConsulMock(t *testing.T) {
	p := newProvider("test").(*consulProvider)
	p.consulConfig.HttpClient = http.DefaultClient
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://consul.ci:8080/v1/catalog/datacenters",
		httpmock.NewStringResponder(200, `["site1","site2"]`))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/nodes.*site1",
		httpmock.NewStringResponder(200, mockHosts("site1")))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/nodes.*site2",
		httpmock.NewStringResponder(200, mockHosts("site2")))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/services",
		httpmock.NewStringResponder(200, `{"service1":["tag1","tag2"],"service2":["tag1","tag2"]}`))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/service/service1.*site1",
		httpmock.NewStringResponder(200, mockServices("site1", "service1", 1)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/service/service1.*site2",
		httpmock.NewStringResponder(200, mockServices("site2", "service1", 1)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/service/service2.*site1",
		httpmock.NewStringResponder(200, mockServices("site1", "service2", 2)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/catalog/service/service2.*site2",
		httpmock.NewStringResponder(200, mockServices("site2", "service2", 2)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/health/checks/service1.*site1",
		httpmock.NewStringResponder(200, mockHealthchecks("site1", "service1", "passing", 1)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/health/checks/service1.*site2",
		httpmock.NewStringResponder(200, mockHealthchecks("site2", "service1", "passing", 1)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/health/checks/service2.*site1",
		httpmock.NewStringResponder(200, mockHealthchecks("site1", "service2", "failing", 2)))
	httpmock.RegisterResponder("GET", "=~http://consul.ci:8080/v1/health/checks/service2.*site2",
		httpmock.NewStringResponder(200, mockHealthchecks("site2", "service2", "failing", 2)))

	p.config.Address = "http://consul.ci:8080"
	ctx := context.Background()
	hosts, err := p.Load(ctx, func(string, bool, error) {})
	if err != nil {
		t.Errorf("Failed to query mock consul: %s", err)
	}
	if len(hosts) != 20 {
		t.Errorf("Incorrect number of hosts returned")
	}
	for i, h := range hosts {
		if dc, ok := h.Attributes["datacenter"]; !ok {
			t.Errorf("datacenter attribute not set for %s", h.Name)
		} else if h.Attributes["domainname"] != dc.(string)+".consul.ci" {
			t.Errorf("incorrect datacenter %s set for %s", dc, h.Name)
		}
		if svc, ok := h.Attributes["service"]; !ok {
			t.Errorf("service attribute not set for %s", h.Name)
		} else {
			s := svc.([]string)
			if len(s) != 1+((i+1)%2) {
				t.Errorf("incorrect services %v set for %s", s, h.Name)
			}
			if s[0] != "service1" || len(s) == 2 && s[1] != "service2" {
				t.Errorf("incorrect services %v set for %s", s, h.Name)
			}
			for _, ss := range s {
				if sts, ok := h.Attributes["service:"+ss]; !ok {
					t.Errorf("service tags not set for %s/%s", h.Name, ss)
				} else {
					st := sts.([]string)
					if st[0] != ss || st[1] != ss+"X" {
						t.Errorf("incorrect service tags %v set for %s/%s", st, h.Name, ss)
					}
				}
			}
			if si, ok := h.Attributes["service_healthy"]; !ok {
				t.Errorf("service_healthy attribute not set for %s", h.Name)
			} else {
				s := si.([]string)
				if len(s) != 1 || s[0] != "service1" {
					if len(s) != 1+((i+1)%2) {
						t.Errorf("incorrect services_healthy %v set for %s", s, h.Name)
					}
				}
			}
			if i%2 == 0 {
				if si, ok := h.Attributes["service_unhealthy"]; !ok {
					t.Errorf("service_unhealthy attribute not set for %s", h.Name)
				} else {
					s := si.([]string)
					if len(s) != 1 || s[0] != "service2" {
						if len(s) != 1+((i+1)%2) {
							t.Errorf("incorrect services_unhealthy %v set for %s", s, h.Name)
						}
					}
				}
			}
		}
		_ = i
	}
}

func mockHosts(site string) string {
	nodes := make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		nodes[i] = map[string]interface{}{
			"ID":         uuid.New().String(),
			"Node":       fmt.Sprintf("node-%d.%s.consul.ci", i, site),
			"Address":    "127.0.0.1",
			"Datacenter": site,
			"TaggedAddresses": map[string]string{
				"lan":      "127.0.0.1",
				"lan_ipv4": "127.0.0.1",
				"wan":      "127.0.0.1",
				"wan_ipv4": "127.0.0.1",
			},
			"Meta": map[string]string{
				"consul-network-segment": "",
			},
			"NodeMeta": map[string]string{
				"consul-network-segment": "",
			},
			"CreateIndex": i,
			"ModifyIndex": i,
		}
	}
	data, _ := json.Marshal(nodes)
	return string(data)
}

func mockServices(site, service string, skip int) string {
	services := make([]map[string]interface{}, 0)
	for i := 0; i < 10; i += skip {
		svc := map[string]interface{}{
			"ID":         uuid.New().String(),
			"Node":       fmt.Sprintf("node-%d.%s.consul.ci", i, site),
			"Address":    "127.0.0.1",
			"Datacenter": site,
			"TaggedAddresses": map[string]string{
				"lan":      "127.0.0.1",
				"lan_ipv4": "127.0.0.1",
				"wan":      "127.0.0.1",
				"wan_ipv4": "127.0.0.1",
			},
			"NodeMeta": map[string]string{
				"consul-network-segment": "",
			},
			"ServiceKind":    "",
			"ServiceID":      service,
			"ServiceName":    service,
			"ServiceTags":    []string{service, service + "X"},
			"ServiceAddress": "127.0.0.2",
			"ServiceTaggedAddresses": map[string]map[string]interface{}{
				"lan_ipv4": {
					"Address": "127.0.0.2",
					"Port":    12345,
				},
				"wan_ipv4": {
					"Address": "127.0.0.2",
					"Port":    12345,
				},
			},
			"ServiceWeights": map[string]int{
				"Passing": 1,
				"Warning": 1,
			},
			"ServiceMeta":              map[string]string{},
			"ServicePort":              12345,
			"ServiceEnableTagOverride": false,
			"ServiceProxy": map[string]map[string]string{
				"MeshGateway": {},
				"Expose":      {},
			},
			"ServiceConnect": map[string]string{},
			"CreateIndex":    i,
			"ModifyIndex":    i,
		}
		services = append(services, svc)
	}
	data, _ := json.Marshal(services)
	return string(data)
}

func mockHealthchecks(site, service, status string, skip int) string {
	checks := make([]map[string]interface{}, 0)
	for i := 0; i < 10; i += skip {
		check := map[string]interface{}{
			"Node":        fmt.Sprintf("node-%d.%s.consul.ci", i, site),
			"CheckID":     fmt.Sprintf("service:%s:1", service),
			"Name":        "Service check 1",
			"Status":      "passing",
			"Notes":       "CI check",
			"Output":      "",
			"ServiceID":   service,
			"ServiceName": service,
			"ServiceTags": []string{},
			"Type":        "script",
			"Definition":  map[string]interface{}{},
			"CreateIndex": 1807514456,
			"ModifyIndex": 1807514456,
		}
		checks = append(checks, check)
		check = map[string]interface{}{
			"Node":        fmt.Sprintf("node-%d.%s.consul.ci", i, site),
			"CheckID":     fmt.Sprintf("service:%s:2", service),
			"Name":        "Service check 2",
			"Status":      status,
			"Notes":       "CI check",
			"Output":      "",
			"ServiceID":   service,
			"ServiceName": service,
			"ServiceTags": []string{},
			"Type":        "script",
			"Definition":  map[string]interface{}{},
			"CreateIndex": 1807514456,
			"ModifyIndex": 1807514456,
		}
		checks = append(checks, check)
	}
	data, _ := json.Marshal(checks)
	return string(data)
}
