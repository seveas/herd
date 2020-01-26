package herd

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"io"
	"strings"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

var puttyNameMap = map[string]string{}

const agentMsgMax = 8192

// Pageant 0.73 and lower only work via WM_COPYDATA and shared memory, but
// they do use the agent protocol. So we pretend to be a socket, but
// communicate with pageant via WM_COPYDATA
type pageantWrapper struct {
	hwnd   win.HWND
	outbuf *bytes.Buffer
	inbuf  *bytes.Buffer
	mlen   int
}

type copyDataStruct struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

func findPageant() io.ReadWriter {
	name, _ := windows.UTF16PtrFromString("Pageant")
	hwnd := win.FindWindow(name, name)
	if hwnd == 0 {
		return nil
	}
	return &pageantWrapper{hwnd: hwnd, outbuf: new(bytes.Buffer), inbuf: new(bytes.Buffer)}
}

func (p *pageantWrapper) Read(b []byte) (int, error) {
	return p.inbuf.Read(b)
}

func (p *pageantWrapper) Write(b []byte) (int, error) {
	l, _ := p.outbuf.Write(b)
	var lb []byte
	if p.mlen == 0 && l > 4 {
		lb = make([]byte, 4)
		io.ReadFull(p.outbuf, lb)
		p.mlen = int(binary.BigEndian.Uint32(lb))
	}
	// If we have a complete message, query pageant
	if p.outbuf.Len() >= p.mlen {
		mb := make([]byte, p.mlen)
		io.ReadFull(p.outbuf, mb)
		p.mlen = 0
		return len(b), p.queryPageant(append(lb, mb...))
	}
	return len(b), nil
}

func (p *pageantWrapper) queryPageant(b []byte) error {
	if len(b) > agentMsgMax {
		return fmt.Errorf("Message too long")
	}

	// We need this message in shared memory for pageant
	mapName := "HerdPageantRequest"
	mapNamePtr, _ := windows.UTF16PtrFromString(mapName)
	fileMap, err := windows.CreateFileMapping(^windows.Handle(0), nil, 0x4, 0, agentMsgMax, mapNamePtr)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(fileMap)
	sharedMemory, err := windows.MapViewOfFile(fileMap, 0x2, 0, 0, 0)
	if err != nil {
		return err
	}
	defer windows.UnmapViewOfFile(sharedMemory)
	sharedMemoryArray := (*[agentMsgMax]byte)(unsafe.Pointer(sharedMemory))
	copy(sharedMemoryArray[:], b)

	// And now we talk to pageant
	mapNameZ := []byte(mapName + "\000")
	data := copyDataStruct{
		dwData: 0x804e50ba,
		cbData: uint32(len(mapNameZ)),
		lpData: uintptr(unsafe.Pointer(&mapNameZ[0])),
	}

	if win.SendMessage(p.hwnd, win.WM_COPYDATA, 0, uintptr(unsafe.Pointer(&data))) == 0 {
		return fmt.Errorf("WM_COPYDATA failed")
	}

	// Pageants results go in the buffer for the reader to read
	l := binary.BigEndian.Uint32(sharedMemoryArray[:4])+4
	p.inbuf.Write(sharedMemoryArray[:l])
	return nil
}

func puttyConfig(host string) map[string]string {
	ret := make(map[string]string)
	if v, ok := puttyNameMap[host]; ok {
		host = v
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, fmt.Sprintf(`Software\SimonTatham\PuTTY\Sessions\%s`, host), registry.QUERY_VALUE)
	if err != nil {
		k, err = registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\Sessions\Default%20Settings`, registry.QUERY_VALUE)
	}
	if err != nil {
		return ret
	}
	iv, _, err := k.GetIntegerValue("PortNumber")
	if err == nil {
		ret["port"] = fmt.Sprintf("%d", iv)
	}
	sv, _, err := k.GetStringValue("UserName")
	if err == nil {
		ret["user"] = sv
	}
	return ret
}

type PuttyProvider struct {
	Name string
}

func (p *PuttyProvider) String() string {
	return p.Name
}

func (p *PuttyProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func (p *PuttyProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	keys := p.allKeys()
	ret := Hosts{}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\Sessions`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return ret, err
	}
	defer k.Close()
	names, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return ret, err
	}
	for _, name := range names {
		k, err := registry.OpenKey(k, name, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		defer k.Close()
		hn, _, err := k.GetStringValue("HostName")
		puttyNameMap[hn] = name
		if err != nil || hn == "" {
			continue
		}
		h := NewHost(hn, HostAttributes{})
		for _, k := range keys[hn] {
			h.AddPublicKey(k)
		}
		delete(keys, hn)
		ret = append(ret, h)
	}
	for hn, hkeys := range keys {
		h := NewHost(hn, HostAttributes{})
		for _, k := range hkeys {
			h.AddPublicKey(k)
		}
		ret = append(ret, h )
	}
	return ret, nil
}

func (p *PuttyProvider) allKeys() map[string][]ssh.PublicKey {
	ret := make(map[string][]ssh.PublicKey)
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\SshHostKeys`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Err 1")
		return ret
	}
	defer k.Close()
	keys, err := k.ReadValueNames(-1)
	if err != nil {
		fmt.Println("Err 2", err)
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