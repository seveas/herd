package putty

import (
	"bytes"
	"context"
	"encoding/binary"
	"math/big"
	"strings"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows/registry"
)

func init() {
	herd.RegisterProvider("putty", newProvider, magicProvider)
}

type puttyProvider struct {
	name   string
	config struct {
		Prefix string
	}
}

func newProvider(name string) herd.HostProvider {
	return &puttyProvider{name: name}
}

func magicProvider() herd.HostProvider {
	return newProvider("putty")
}

func (p *puttyProvider) Name() string {
	return p.name
}

func (p *puttyProvider) Prefix() string {
	return p.config.Prefix
}

func (p *puttyProvider) Equivalent(o herd.HostProvider) bool {
	return true
}

func (p *puttyProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *puttyProvider) Load(ctx context.Context, lm herd.LoadingMessage) (herd.Hosts, error) {
	keys := p.allKeys()
	ret := herd.Hosts{}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\Sessions`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		logrus.Debugf("No putty sessions found")
		return nil, nil
	}
	defer k.Close()
	names, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		k, err := registry.OpenKey(k, name, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		defer k.Close()
		hn, _, err := k.GetStringValue("HostName")
		herd.PuttyNameMap[hn] = name
		if err != nil || hn == "" {
			continue
		}
		h := herd.NewHost(hn, "", herd.HostAttributes{})
		for _, k := range keys[hn] {
			h.AddPublicKey(k)
		}
		delete(keys, hn)
		ret = append(ret, h)
	}
	for hn, hkeys := range keys {
		h := herd.NewHost(hn, "", herd.HostAttributes{})
		for _, k := range hkeys {
			h.AddPublicKey(k)
		}
		ret = append(ret, h)
	}
	return ret, nil
}

func (p *puttyProvider) allKeys() map[string][]ssh.PublicKey {
	ret := make(map[string][]ssh.PublicKey)
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\SshHostKeys`, registry.QUERY_VALUE)
	if err != nil {
		return ret
	}
	defer k.Close()
	keys, err := k.ReadValueNames(-1)
	if err != nil {
		return ret
	}
	for _, kn := range keys {
		key, _, err := k.GetStringValue(kn)
		var b sshBuffer
		// Format of the key name is: algorithm@port:host
		hn := kn[strings.IndexRune(kn, ':')+1:]
		algo := kn[:strings.IndexRune(kn, '@')]
		parts := strings.Split(key, ",")
		if strings.HasPrefix(algo, "ecdsa") {
			b.Write([]byte(algo))
			b.Write([]byte(parts[0]))
			x := big.NewInt(0)
			y := big.NewInt(0)
			x.SetString(parts[1], 0)
			y.SetString(parts[2], 0)
			var b2 bytes.Buffer
			// PuTTY doesn't do point compression
			b2.Write([]byte{4})
			b2.Write(x.Bytes())
			b2.Write(y.Bytes())
			b.Write(b2.Bytes())
		} else if algo == "rsa2" {
			b.Write([]byte("ssh-rsa"))
			e := big.NewInt(0)
			e.SetString(parts[0], 0)
			b.Write(e.Bytes())
			var b2 bytes.Buffer
			// Extra padding required in the protocol
			b2.Write([]byte{0})
			n := big.NewInt(0)
			n.SetString(parts[1], 0)
			b2.Write(n.Bytes())
			b.Write(b2.Bytes())
		} else if algo == "ssh-ed25519" {
			b.Write([]byte(algo))
			y := big.NewInt(0)
			y.SetString(parts[1], 0)
			yb := y.Bytes()
			l := len(yb)
			// Correct endianness
			for i := 0; i < l/2; i++ {
				yb[i], yb[l-1-i] = yb[l-1-i], yb[i]
			}
			b.Write(yb)
		} else {
			logrus.Warnf("Unsupported public key type in PuTTY: %s", algo)
			continue
		}
		pubkey, err := ssh.ParsePublicKey(b.Bytes())
		if err != nil {
			logrus.Warnf("Unable to parse PuTTY known hostkey for %s: %s", kn, err)
			continue
		}
		if _, ok := ret[hn]; ok {
			ret[hn] = append(ret[hn], pubkey)
		} else {
			ret[hn] = []ssh.PublicKey{pubkey}
		}
	}
	return ret
}

type sshBuffer struct {
	buf bytes.Buffer
}

func (buf *sshBuffer) Write(b []byte) {
	l := make([]byte, 4)
	binary.BigEndian.PutUint32(l, uint32(len(b)))
	buf.buf.Write(l)
	buf.buf.Write(b)
}

func (buf *sshBuffer) Bytes() []byte {
	return buf.buf.Bytes()
}
