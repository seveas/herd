package ssh

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"syscall"
	"unsafe"

	"github.com/miekg/dns"
	"golang.org/x/sys/windows"
)

const (
	dnsConfigDnsServerList int32 = 6
)

func getDnsConfig() (*dns.ClientConfig, error) {
	dnsapi := windows.NewLazyDLL("Dnsapi.dll")
	dnsQueryConfig := dnsapi.NewProc("DnsQueryConfig")
	buffer := make([]byte, 60)

	blen := len(buffer)
	r1, _, _ := dnsQueryConfig.Call(uintptr(dnsConfigDnsServerList), uintptr(0), uintptr(0), uintptr(0), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&blen)))
	if syscall.Errno(r1) == windows.ERROR_MORE_DATA {
		buffer = make([]byte, blen)
		r1, _, _ = dnsQueryConfig.Call(uintptr(dnsConfigDnsServerList), uintptr(0), uintptr(0), uintptr(0), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&blen)))
	}

	if syscall.Errno(r1) != windows.NO_ERROR {
		return nil, fmt.Errorf("Unable to get dns configuration")
	}

	reader := bytes.NewReader(buffer)
	var l uint32
	binary.Read(reader, binary.LittleEndian, &l)
	ips := make([]string, l)
	for i, _ := range ips {
		ips[i] = net.IP(buffer[4*(i+1) : 4*(i+2)]).String()
	}
	config := fmt.Sprintf("nameserver %s\n", strings.Join(ips, " "))
	return dns.ClientConfigFromReader(strings.NewReader(config))
}
