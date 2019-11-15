// +build !windows

package herd

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

func GetTerminalWidth() int {
	var termDim [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(0), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&termDim)), 0, 0, 0); err != 0 {
		return 80
	}
	return int(termDim[1])
}

func ListenForWindowChange(ui *SimpleUI) {
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGWINCH)
	for {
		<-sc
		ui.Width = GetTerminalWidth()
	}
}
