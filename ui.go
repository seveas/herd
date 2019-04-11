package herd

import (
	"log"
	"os"
)

type HerdUI interface {
	Warnf(format string, v ...interface{})
}

type SimpleUI struct {
	Logger *log.Logger
}

func NewSimpleUI() SimpleUI {
	return SimpleUI{
		Logger: log.New(os.Stderr, "", 0),
	}
}

func (ui SimpleUI) Warnf(fmt string, v ...interface{}) {
	ui.Logger.Printf(fmt, v...)
}

var UI HerdUI
