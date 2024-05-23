package ssh

import (
	"crypto/sha1" // #nosec G505 -- We want to support sha1 fingerprints for now
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/miekg/dns"
	"golang.org/x/crypto/ssh"
)

func init() {
	var err error
	config, err := getDnsConfig()
	if err == nil {
		sshfpResolver = &dnsResolver{client: new(dns.Client), config: config}
	}
}

type dnsResolver struct {
	client *dns.Client
	config *dns.ClientConfig
}

var sshfpResolver *dnsResolver

func (r *dnsResolver) resolve(hostname string, qtype uint16) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), qtype)
	m.RecursionDesired = true

	resp, _, err := r.client.Exchange(m, net.JoinHostPort(r.config.Servers[0], r.config.Port))
	if err != nil {
		return nil, err
	}

	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS error (Rcode %d)", resp.Rcode)
	}

	return resp.Answer, nil
}

var sshfpAlgorithms = map[string]uint8{
	ssh.KeyAlgoRSA:      1,
	ssh.KeyAlgoDSA:      2,
	ssh.KeyAlgoECDSA256: 3,
	ssh.KeyAlgoECDSA384: 3,
	ssh.KeyAlgoECDSA521: 3,
	ssh.KeyAlgoED25519:  4,
}

func verifyHostKeyDns(hostname string, key ssh.PublicKey) bool {
	if sshfpResolver == nil {
		return false
	}
	algo, ok := sshfpAlgorithms[key.Type()]
	if !ok {
		return false
	}
	blob := key.Marshal()
	sha1sum := fmt.Sprintf("%x", sha1.Sum(blob)) // #nosec:G401 -- We want to support sha1 fingerprints for now
	sha256sum := fmt.Sprintf("%x", sha256.Sum256(blob))

	rrset, err := sshfpResolver.resolve(hostname+".", dns.TypeSSHFP)
	if err != nil {
		return false
	}
	for _, rr := range rrset {
		if srr, ok := rr.(*dns.SSHFP); ok {
			if srr.Algorithm == algo {
				if srr.Type == dns.SHA1 && srr.FingerPrint == sha1sum {
					return true
				}
				if srr.Type == dns.SHA256 && srr.FingerPrint == sha256sum {
					return true
				}
			}
		}
	}
	return false
}
