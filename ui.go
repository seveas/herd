package herd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mgutz/ansi"
)

type HerdUI interface {
	Println(str string)
	Printf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Progress(total, todo, queued, doneOk, doneFaile, doneError int)
	PrintHistoryItem(hi HistoryItem)
}

type SimpleUI struct {
	AtStart      bool
	LastProgress string
	Pchan        chan string
	Config       *UIConfig
}

func NewSimpleUI(config *UIConfig) *SimpleUI {
	c := make(chan string)
	ui := &SimpleUI{
		AtStart:      true,
		LastProgress: "",
		Pchan:        c,
		Config:       config,
	}
	go ui.Printer()
	return ui
}

func (ui *SimpleUI) Printer() {
	for {
		msg := <-ui.Pchan
		// If we're getting a normal message in the middle of printing progress, wipe the progress message
		if !ui.AtStart && msg[0] != '\r' && msg[0] != '\n' {
			os.Stderr.WriteString("\r\033[2K")
			os.Stderr.WriteString(msg)
			// After printing the real message, re-write the progress message
			msg = ui.LastProgress
		}
		os.Stderr.WriteString(msg)
		if msg[len(msg)-1] == '\n' {
			ui.AtStart = true
		} else {
			ui.AtStart = false
			ui.LastProgress = msg
			os.Stderr.Sync()
		}
	}
}

func (ui *SimpleUI) PrintHistoryItem(hi HistoryItem) {
	buf := strings.Builder{}
	ui.Config.Formatter.FormatHistoryItem(hi, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) Println(str string) {
	if ui.Config.LogLevel < INFO {
		return
	}
	ui.Pchan <- str + "\n"
}

func (ui *SimpleUI) Printf(format string, v ...interface{}) {
	if ui.Config.LogLevel < INFO {
		return
	}
	ui.Pchan <- fmt.Sprintf(format, v...)
}

func (ui *SimpleUI) Errorf(format string, v ...interface{}) {
	if ui.Config.LogLevel < ERROR {
		return
	}
	format = ansi.Color(format, "red+b")
	ui.Printf(format+"\n", v...)
}

func (ui *SimpleUI) Warnf(format string, v ...interface{}) {
	if ui.Config.LogLevel < WARNING {
		return
	}
	format = ansi.Color(format, "yellow")
	ui.Printf(format+"\n", v...)
}

func (ui *SimpleUI) Debugf(format string, v ...interface{}) {
	if ui.Config.LogLevel < DEBUG {
		return
	}
	format = ansi.Color(format, "black+h")
	ui.Printf(format+"\n", v...)
}

func (ui *SimpleUI) Progress(total, todo, queued, doneOk, doneFail, doneError int) {
	if ui.Config.LogLevel < INFO {
		return
	}
	if queued >= 0 {
		ui.Printf("\r\033[2kWaiting... %d/%d done, %d queued, %d ok, %d fail, %d error", total-todo, total, queued, doneOk, doneFail, doneError)
	} else {
		ui.Printf("\r\033[2KWaiting... %d/%d done, %d ok, %d fail, %d error", total-todo, total, doneOk, doneFail, doneError)
	}
	if todo == 0 {
		ui.Printf("\n")
	}
}

var UI HerdUI
