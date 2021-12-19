package herd

import (
	"errors"
	"github.com/miekg/dns"
)

func getDnsResolvers() (*dns.ClientConfig, error) {
	return nil, errors.New("SSHFP lookups are not yet supported under windows")
}
