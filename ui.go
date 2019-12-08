package katyusha

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/seveas/readline"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OutputMode int

const (
	OutputTail OutputMode = iota
	OutputPerhost
	OutputInline
	OutputAll
)

type UI interface {
	PrintHistoryItem(hi *HistoryItem)
	PrintHostList(hosts Hosts, oneline, csvOutput, allAttributes bool, attributes []string)
	SetOutputMode(OutputMode)
	SetPagerEnabled(bool)
	Write([]byte) (int, error)
	Wait()
	CacheUpdateChannel() chan CacheMessage
	OutputChannel(r *Runner) chan OutputLine
	ProgressChannel(r *Runner) chan ProgressMessage
}

type SimpleUI struct {
	Output       *os.File
	AtStart      bool
	LastProgress string
	Pchan        chan string
	Dchan        chan interface{}
	Formatter    Formatter
	OutputMode   OutputMode
	PagerEnabled bool
	Width        int
	Height       int
	LineBuf      string
}

func NewSimpleUI(f Formatter) *SimpleUI {
	ui := &SimpleUI{
		Output:       os.Stdout,
		OutputMode:   OutputAll,
		AtStart:      true,
		LastProgress: "",
		Pchan:        make(chan string),
		Dchan:        make(chan interface{}),
		Formatter:    f,
	}
	ui.GetSize()
	readline.DefaultOnWidthChanged(func() {
		ui.GetSize()
	})
	go ui.Printer()
	return ui
}

func (ui *SimpleUI) GetSize() {
	w, h, err := readline.GetSize(int(ui.Output.Fd()))
	if err == nil {
		ui.Width, ui.Height = w, h
		if w < 40 {
			ui.Width = 40
		}
	} else {
		ui.PagerEnabled = false
		ui.Width = 80
	}
}

func (ui *SimpleUI) SetOutputMode(o OutputMode) {
	ui.OutputMode = o
}

func (ui *SimpleUI) SetPagerEnabled(e bool) {
	ui.PagerEnabled = e
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
	ui.Pchan = make(chan string)
	<-ui.Dchan
	ui.Dchan = make(chan interface{})
	go ui.Printer()
}

func (ui *SimpleUI) PrintHistoryItem(hi *HistoryItem) {
	if ui.OutputMode != OutputAll && ui.OutputMode != OutputInline {
		return
	}
	usePager := ui.PagerEnabled
	hlen := hi.Hosts.MaxLen()
	linecount := 0
	buffer := ""
	var pager *Pager
	if usePager {
		buffer = ui.Formatter.FormatCommand(hi.Command)
		linecount = 1
	} else {
		ui.Pchan <- ui.Formatter.FormatCommand(hi.Command)
	}

	for _, h := range hi.Hosts {
		var txt string
		if ui.OutputMode == OutputAll {
			txt = ui.Formatter.FormatResult(hi.Results[h.Name])
		} else {
			txt = ui.Formatter.FormatOutput(hi.Results[h.Name], hlen)
		}
		if !usePager {
			ui.Pchan <- txt
		} else if pager != nil {
			pager.WriteString(txt)
		} else {
			buffer += txt
			linecount += strings.Count(txt, "\n")
			if linecount > ui.Height {
				pager = &Pager{}
				if err := pager.Start(); err != nil {
					logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
					ui.Pchan <- buffer
					usePager = false
				} else {
					pager.WriteString(buffer)
				}
				buffer = ""
			}
		}
	}
	if buffer != "" {
		ui.Pchan <- buffer
	}
	if pager != nil {
		pager.Wait()
	}
}

func (ui *SimpleUI) PrintHostList(hosts Hosts, oneline, csvOutput, allAttributes bool, attributes []string) {
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

func (ui *SimpleUI) OutputChannel(r *Runner) chan OutputLine {
	if ui.OutputMode != OutputTail {
		return nil
	}
	oc := make(chan OutputLine)
	hlen := r.Hosts.MaxLen()
	go func() {
		for msg := range oc {
			name := fmt.Sprintf("%-*s", hlen, msg.Host.Name)
			if msg.Stderr {
				name = ansi.Color(name, "red")
			}
			ui.Pchan <- fmt.Sprintf("%s  %s", name, msg.Data)
		}
	}()
	return oc
}

func (ui *SimpleUI) ProgressChannel(r *Runner) chan ProgressMessage {
	pc := make(chan ProgressMessage)
	go func() {
		start := time.Now()
		ticker := time.NewTicker(time.Second / 2)
		defer ticker.Stop()
		total := len(r.Hosts)
		queued, todo := total, total
		nok, nfail, nerr := 0, 0, 0
		hlen := r.Hosts.MaxLen()
		for {
			select {
			case <-ticker.C:
			case msg, ok := <-pc:
				if !ok {
					return
				}
				if msg.Result == nil {
					queued--
					continue
				}
				if msg.Result.ExitStatus == -1 {
					nerr++
				} else if msg.Result.ExitStatus == 0 {
					nok++
				} else {
					nfail++
				}
				if ui.OutputMode == OutputPerhost {
					ui.Pchan <- ui.Formatter.FormatResult(msg.Result)
				} else if ui.OutputMode == OutputTail {
					ui.Pchan <- ui.Formatter.FormatStatus(msg.Result, hlen)
				}
				todo--
			}
			since := time.Since(start).Truncate(time.Second)
			togo := viper.GetDuration("Timeout") - since
			if todo == 0 {
				ui.Pchan <- fmt.Sprintf("\r\033[2K%d done, %d ok, %d fail, %d error in %s\n", total, nok, nfail, nerr, since)
			} else if queued >= 0 {
				ui.Pchan <- fmt.Sprintf("\r\033[2KWaiting (%s/%s)... %d/%d done, %d queued, %d ok, %d fail, %d error", since, togo, total-todo, total, queued, nok, nfail, nerr)
			} else {
				ui.Pchan <- fmt.Sprintf("\r\033[2KWaiting (%s/%s)... %d/%d done, %d ok, %d fail, %d error", since, togo, total-todo, total, nok, nfail, nerr)
			}
		}
	}()
	return pc
}
