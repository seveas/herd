package katyusha

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/spf13/viper"
)

const (
	ERROR = iota
	WARNING
	NORMAL
	INFO
	DEBUG
)

type KatyushaUI interface {
	Println(str string)
	Printf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Progress(start time.Time, total, todo, queued, doneOk, doneFaile, doneError int)
	PrintHistoryItem(hi HistoryItem)
	PrintHistoryItemWithPager(hi HistoryItem)
	PrintCommand(command string)
	PrintResult(r Result)
	Wait()
	NewLineWriterBuffer(host *Host, prefix string, isError bool) *LineWriterBuffer
}

type SimpleUI struct {
	Output       *os.File
	AtStart      bool
	LastProgress string
	Pchan        chan string
	Dchan        chan interface{}
	Formatter    Formatter
}

func NewSimpleUI() *SimpleUI {
	ui := &SimpleUI{
		Output:       os.Stdout,
		AtStart:      true,
		LastProgress: "",
		Pchan:        make(chan string),
		Dchan:        make(chan interface{}),
		Formatter:    Formatters[viper.GetString("Formatter")],
	}
	go ui.Printer()
	return ui
}


func (ui *SimpleUI) Printer() {
	for msg := range ui.Pchan {
		// If we're getting a normal message in the middle of printing
		// progress, wipe the progress message and reprint it after this
		// message
		if !ui.AtStart && msg[0] != '\r' && msg[0] != '\n' {
			ui.Output.WriteString("\r\033[2K" + msg + ui.LastProgress)
		} else {
			ui.Output.WriteString(msg)
			if msg[len(msg)-1] == '\n' {
				ui.AtStart = true
			} else {
				ui.AtStart = false
				ui.LastProgress = msg
			}
		}
		ui.Output.Sync()
	}
	close(ui.Dchan)
}

func (ui *SimpleUI) Wait() {
	close(ui.Pchan)
	<-ui.Dchan
}

func (ui *SimpleUI) PrintCommand(command string) {
	if viper.GetInt("LogLevel") < NORMAL {
		return
	}
	buf := strings.Builder{}
	ui.Formatter.FormatCommand(command, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) PrintHistoryItem(hi HistoryItem) {
	if viper.GetInt("LogLevel") < NORMAL {
		return
	}
	buf := strings.Builder{}
	ui.printHistoryItem(hi, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) PrintHistoryItemWithPager(hi HistoryItem) {
	p := Pager{}
	if err := p.Start(); err != nil {
		ui.Errorf("Unable to start pager, displaying on stdout: %s", err)
		ui.PrintHistoryItem(hi)
	} else {
		ui.printHistoryItem(hi, &p)
		p.Wait()
	}
}

func (ui *SimpleUI) printHistoryItem(hi HistoryItem, w io.Writer) {
	ui.Formatter.FormatCommand(hi.Command, w)
	for _, h := range hi.Hosts {
		if ui.Filtered(h) {
			ui.Formatter.FormatResult(hi.Results[h.Name], w)
		}
	}
}

func (ui *SimpleUI) PrintResult(r Result) {
	if viper.GetInt("LogLevel") < NORMAL {
		return
	}
	buf := strings.Builder{}
	ui.Formatter.FormatResult(r, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) Println(str string) {
	if viper.GetInt("LogLevel") < NORMAL {
		return
	}
	ui.Pchan <- str + "\n"
}

func (ui *SimpleUI) Printf(format string, v ...interface{}) {
	if viper.GetInt("LogLevel") < NORMAL {
		return
	}
	ui.Pchan <- fmt.Sprintf(format, v...)
}

func (ui *SimpleUI) Infof(format string, v ...interface{}) {
	if viper.GetInt("LogLevel") < INFO {
		return
	}
	ui.Pchan <- fmt.Sprintf(format+"\n", v...)
}

func (ui *SimpleUI) Errorf(format string, v ...interface{}) {
	if viper.GetInt("LogLevel") < ERROR {
		return
	}
	format = ansi.Color(format, "red+b")
	ui.Pchan <- fmt.Sprintf(format+"\n", v...)
}

func (ui *SimpleUI) Warnf(format string, v ...interface{}) {
	if viper.GetInt("LogLevel") < WARNING {
		return
	}
	format = ansi.Color(format, "yellow")
	ui.Pchan <- fmt.Sprintf(format+"\n", v...)
}

func (ui *SimpleUI) Debugf(format string, v ...interface{}) {
	if viper.GetInt("LogLevel") < DEBUG {
		return
	}
	format = ansi.Color(format, "black+h")
	ui.Pchan <- fmt.Sprintf(format+"\n", v...)
}

func (ui *SimpleUI) Progress(start time.Time, total, todo, queued, doneOk, doneFail, doneError int) {
	if viper.GetInt("LogLevel") < INFO {
		return
	}
	since := time.Since(start).Truncate(time.Second)
	if todo == 0 {
		ui.Pchan <- fmt.Sprintf("\r\033[2K%d done, %d ok, %d fail, %d error in %s\n", total, doneOk, doneFail, doneError, since)
	} else if queued >= 0 {
		ui.Pchan <- fmt.Sprintf("\r\033[2KWaiting (%s)... %d/%d done, %d queued, %d ok, %d fail, %d error", since, total-todo, total, queued, doneOk, doneFail, doneError)
	} else {
		ui.Pchan <- fmt.Sprintf("\r\033[2KWaiting (%s)... %d/%d done, %d ok, %d fail, %d error", since, total-todo, total, doneOk, doneFail, doneError)
	}
}

type ByteWriter interface {
	Write([]byte) (int, error)
	Bytes() []byte
}

type LineWriterBuffer struct {
	buf     *bytes.Buffer
	lineBuf []byte
	host    *Host
	prefix  string
	pos     int
	ui      *SimpleUI
}

func (ui *SimpleUI) NewLineWriterBuffer(host *Host, prefix string, isError bool) *LineWriterBuffer {
	if isError {
		prefix = ansi.Color(prefix, "red")
	}
	return &LineWriterBuffer{buf: bytes.NewBuffer([]byte{}), lineBuf: []byte{}, host: host, prefix: prefix, pos: 0, ui: ui}
}

func (buf *LineWriterBuffer) Write(p []byte) (int, error) {
	n, err := buf.buf.Write(p)
	buf.lineBuf = bytes.Join([][]byte{buf.lineBuf, p}, []byte{})
	for {
		idx := bytes.Index(buf.lineBuf, []byte("\n"))
		if idx == -1 {
			break
		}
		buf.ui.Printf("%s %s", buf.prefix, buf.lineBuf[:idx+1])
		buf.lineBuf = buf.lineBuf[idx+1:]
	}
	return n, err
}

func (buf *LineWriterBuffer) WriteStatus(r Result) {
	if !buf.ui.Filtered(buf.host) {
		return
	}
	sb := strings.Builder{}
	buf.ui.Formatter.FormatStatus(r, &sb)
	buf.Write([]byte(sb.String()))
}

func (buf *LineWriterBuffer) Bytes() []byte {
	return buf.buf.Bytes()
}

var UI KatyushaUI
