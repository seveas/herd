package katyusha

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/readline"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type KatyushaUI interface {
	Printf(format string, v ...interface{})
	Progress(start time.Time, total, todo, queued, doneOk, doneFaile, doneError int)
	PrintHistoryItem(hi *HistoryItem)
	PrintHistoryItemWithPager(hi *HistoryItem)
	PrintCommand(command string)
	PrintResult(r Result)
	PrintHostList(hosts Hosts, oneline, allAttributes bool, attributes []string, csvOutput bool)
	Write([]byte) (int, error)
	Wait()
	NewLineWriterBuffer(host *Host, prefix string, isError bool) *LineWriterBuffer
	SetOutputFilter([]MatchAttributes)
	CacheUpdateChannel() chan CacheMessage
}

type SimpleUI struct {
	Output       *os.File
	AtStart      bool
	LastProgress string
	Pchan        chan string
	Dchan        chan interface{}
	Formatter    Formatter
	OutputFilter []MatchAttributes
	Width        int
	LineBuf      string
}

func NewSimpleUI(f Formatter) *SimpleUI {
	ui := &SimpleUI{
		Output:       os.Stdout,
		AtStart:      true,
		LastProgress: "",
		Pchan:        make(chan string),
		Dchan:        make(chan interface{}),
		Formatter:    f,
		OutputFilter: []MatchAttributes{},
		Width:        readline.GetScreenWidth(),
	}
	if ui.Width == -1 {
		ui.Width = 80
	}
	if ui.Width < 40 {
		ui.Width = 40
	}
	readline.DefaultOnWidthChanged(func() {
		w := readline.GetScreenWidth()
		if w >= 40 {
			ui.Width = w
		}
	})
	go ui.Printer()
	return ui
}

func (ui *SimpleUI) SetOutputFilter(m []MatchAttributes) {
	ui.OutputFilter = m
}

func (ui *SimpleUI) Filtered(h *Host) bool {
	if len(ui.OutputFilter) == 0 {
		return true
	}
	for _, f := range ui.OutputFilter {
		if h.Match("", f) {
			return true
		}
	}
	return false
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
			if msg[len(msg)-1] == '\n' || msg == "\r\033[2K" {
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

func (ui *SimpleUI) Write(msg []byte) (int, error) {
	ui.LineBuf += string(msg)
	if strings.HasSuffix(ui.LineBuf, "\n") {
		ui.Pchan <- ui.LineBuf
		ui.LineBuf = ""
	}
	return len(msg), nil
}

func (ui *SimpleUI) Wait() {
	close(ui.Pchan)
	<-ui.Dchan
}

func (ui *SimpleUI) PrintCommand(command string) {
	buf := strings.Builder{}
	ui.Formatter.FormatCommand(command, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) PrintHistoryItem(hi *HistoryItem) {
	buf := strings.Builder{}
	ui.printHistoryItem(hi, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) PrintHistoryItemWithPager(hi *HistoryItem) {
	p := Pager{}
	if err := p.Start(); err != nil {
		logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
		ui.PrintHistoryItem(hi)
	} else {
		ui.printHistoryItem(hi, &p)
		p.Wait()
	}
}

func (ui *SimpleUI) printHistoryItem(hi *HistoryItem, w io.Writer) {
	ui.Formatter.FormatCommand(hi.Command, w)
	for _, h := range hi.Hosts {
		if ui.Filtered(h) {
			ui.Formatter.FormatResult(hi.Results[h.Name], w)
		}
	}
}

func (ui *SimpleUI) PrintResult(r *Result) {
	buf := strings.Builder{}
	ui.Formatter.FormatResult(r, &buf)
	ui.Pchan <- buf.String()
}

func (ui *SimpleUI) Printf(format string, v ...interface{}) {
	ui.Pchan <- fmt.Sprintf(format, v...)
}

func (ui *SimpleUI) Progress(start time.Time, total, todo, queued, doneOk, doneFail, doneError int) {
	since := time.Since(start).Truncate(time.Second)
	togo := viper.GetDuration("Timeout") - since
	if todo == 0 {
		ui.Pchan <- fmt.Sprintf("\r\033[2K%d done, %d ok, %d fail, %d error in %s\n", total, doneOk, doneFail, doneError, since)
	} else if queued >= 0 {
		ui.Pchan <- fmt.Sprintf("\r\033[2KWaiting (%s/%s)... %d/%d done, %d queued, %d ok, %d fail, %d error", since, togo, total-todo, total, queued, doneOk, doneFail, doneError)
	} else {
		ui.Pchan <- fmt.Sprintf("\r\033[2KWaiting (%s/%s)... %d/%d done, %d ok, %d fail, %d error", since, togo, total-todo, total, doneOk, doneFail, doneError)
	}
}

func (ui *SimpleUI) PrintHostList(hosts Hosts, oneline, allAttributes bool, attributes []string, csvOutput bool) {
	if oneline {
		names := make([]string, len(hosts))
		for i, host := range hosts {
			names[i] = host.Name
		}
		ui.Pchan <- strings.Join(names, ",")
		return
	}
	if allAttributes || len(attributes) > 0 {
		var writer datawriter
		if csvOutput {
			writer = csv.NewWriter(ui)
		} else {
			writer = NewColumnizer(ui, "   ")
		}
		if allAttributes {
			attrs := make(map[string]bool)
			for _, host := range hosts {
				for key, _ := range host.Attributes {
					attrs[key] = true
				}
			}
			for attr, _ := range attrs {
				attributes = append(attributes, attr)
			}
			sort.Strings(attributes)
			attrline := make([]string, len(attributes)+1)
			attrline[0] = "name"
			copy(attrline[1:], attributes)
			writer.Write(attrline)
		}
		for _, host := range hosts {
			line := make([]string, len(attributes)+1)
			line[0] = host.Name
			for i, attr := range attributes {
				val, ok := host.Attributes[attr]
				if ok {
					line[i+1] = fmt.Sprintf("%v", val)
				} else {
					line[i+1] = ""
				}
			}
			writer.Write(line)
		}
		writer.Flush()
	} else {
		for _, host := range hosts {
			ui.Pchan <- host.Name + "\n"
		}
	}
}

func (ui *SimpleUI) CacheUpdateChannel() chan CacheMessage {
	mc := make(chan CacheMessage)
	go func() {
		start := time.Now()
		ticker := time.NewTicker(time.Second / 2)
		defer ticker.Stop()
		caches := make([]string, 0)
		for {
			select {
			case msg, ok := <-mc:
				// Cache message channel closed, we're done caching
				if !ok {
					ui.Pchan <- fmt.Sprintf("\r\033[2K")
					return
				}
				if msg.err != nil {
					logrus.Errorf("Error contacting %s: %s", msg.name, msg.err)
				}
				if msg.finished {
					logrus.Debugf("Cache updated for %s", msg.name)
					for i, v := range caches {
						if v == msg.name {
							caches = append(caches[:i], caches[i+1:]...)
							break
						}
					}
				} else {
					caches = append(caches, msg.name)
				}
			case <-ticker.C:
			}
			if len(caches) > 0 {
				since := time.Since(start).Truncate(time.Second)
				cs := strings.Join(caches, ", ")
				if len(cs) > ui.Width-25 {
					cs = cs[:ui.Width-30] + "..."
				}
				ui.Pchan <- fmt.Sprintf("\r\033[2K%s Refreshing caches %s", since, ansi.Color(cs, "green"))
			}
		}
	}()
	return mc
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
