package katyusha

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var PuttyNameMap = map[string]string{}

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
	mapName := "KatyushaPageantRequest"
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
	l := binary.BigEndian.Uint32(sharedMemoryArray[:4]) + 4
	p.inbuf.Write(sharedMemoryArray[:l])
	return nil
}

func puttyConfig(host string) map[string]string {
	ret := make(map[string]string)
	if v, ok := PuttyNameMap[host]; ok {
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
