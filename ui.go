package katyusha

import (
	"log"
	"os"

	"github.com/mgutz/ansi"
)

type KatyushaUI interface {
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

type SimpleUI struct {
	Logger *log.Logger
}

func NewSimpleUI() SimpleUI {
	return SimpleUI{
		Logger: log.New(os.Stderr, "", 0),
	}
}

func (ui SimpleUI) Errorf(format string, v ...interface{}) {
	format = ansi.Color(format, "red+b")
	ui.Logger.Printf(format, v...)
}

func (ui SimpleUI) Warnf(format string, v ...interface{}) {
	format = ansi.Color(format, "yellow")
	ui.Logger.Printf(format, v...)
}

func (ui SimpleUI) Debugf(format string, v ...interface{}) {
	format = ansi.Color(format, "black+h")
	ui.Logger.Printf(format, v...)
}

var UI KatyushaUI
