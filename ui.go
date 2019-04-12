package herd

import (
	"fmt"
	"log"
	"os"

	"github.com/mgutz/ansi"
)

type HerdUI interface {
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Progress(total, doneOk, doneFaile, doneError, todo int)
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
func (ui SimpleUI) Progress(total, doneOk, doneFail, doneError, todo int) {
	fmt.Fprintf(os.Stderr, "\033[2k\rWaiting... %d/%d done, %d ok, %d fail, %d error", total-todo, total, doneOk, doneFail, doneError)
	if todo == 0 {
		fmt.Fprintf(os.Stderr, "\n")
	}
	os.Stderr.Sync()
}

var UI HerdUI
