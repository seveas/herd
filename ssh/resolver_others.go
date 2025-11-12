//go:build !windows

package ssh

import (
	"github.com/miekg/dns"
)

func getDnsConfig() (*dns.ClientConfig, error) {
	cc, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	return cc, nil
}
