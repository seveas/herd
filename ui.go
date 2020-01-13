package herd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
	"github.com/seveas/readline"
	"github.com/sirupsen/logrus"
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
	End()
	CacheUpdateChannel() chan CacheMessage
	OutputChannel(r *Runner) chan OutputLine
	ProgressChannel(r *Runner) chan ProgressMessage
	BindLogrus()
}

type SimpleUI struct {
	output       *os.File
	atStart      bool
	lastProgress string
	pchan        chan string
	dchan        chan interface{}
	formatter    formatter
	outputMode   OutputMode
	pagerEnabled bool
	width        int
	height       int
	lineBuf      string
	isTerminal   bool
}

func NewSimpleUI() *SimpleUI {
	f := prettyFormatter{
		colors: map[logrus.Level]string{
			logrus.WarnLevel:  "yellow",
			logrus.ErrorLevel: "red+b",
			logrus.DebugLevel: "black+h",
		},
	}
	ui := &SimpleUI{
		output:       os.Stdout,
		outputMode:   OutputAll,
		atStart:      true,
		lastProgress: "",
		pchan:        make(chan string),
		dchan:        make(chan interface{}),
		formatter:    f,
		isTerminal:   isatty.IsTerminal(os.Stdout.Fd()),
	}
	if ui.isTerminal {
		ui.getSize()
		readline.DefaultOnWidthChanged(func() {
			ui.getSize()
		})
	} else {
		ansi.DisableColors(true)
		ui.pagerEnabled = false
		ui.width = 80
	}
	go ui.printer()
	return ui
}

func (ui *SimpleUI) getSize() {
	w, h, err := readline.GetSize(int(ui.output.Fd()))
	if err == nil {
		ui.width, ui.height = w, h
		if w < 40 {
			ui.width = 40
		}
	} else {
		ui.pagerEnabled = false
		ui.width = 80
	}
}

func (ui *SimpleUI) SetOutputMode(o OutputMode) {
	ui.outputMode = o
}

func (ui *SimpleUI) SetPagerEnabled(e bool) {
	ui.pagerEnabled = e
}

func (ui *SimpleUI) printer() {
	for msg := range ui.pchan {
		// If we're getting a normal message in the middle of printing
		// progress, wipe the progress message and reprint it after this
		// message
		if !ui.atStart && msg[0] != '\r' && msg[0] != '\n' {
			ui.output.WriteString("\r\033[2K" + msg + ui.lastProgress)
		} else {
			ui.output.WriteString(msg)
			if msg[len(msg)-1] == '\n' || msg == "\r\033[2K" {
				ui.atStart = true
			} else {
				ui.atStart = false
				ui.lastProgress = msg
			}
		}
		ui.output.Sync()
	}
	close(ui.dchan)
}

func (ui *SimpleUI) BindLogrus() {
	logrus.SetFormatter(ui.formatter)
	logrus.SetOutput(ui)
}

func (ui *SimpleUI) Write(msg []byte) (int, error) {
	ui.lineBuf += string(msg)
	if strings.HasSuffix(ui.lineBuf, "\n") {
		ui.pchan <- ui.lineBuf
		ui.lineBuf = ""
	}
	return len(msg), nil
}

func (ui *SimpleUI) End() {
	close(ui.pchan)
	<-ui.dchan
}

func (ui *SimpleUI) PrintHistoryItem(hi *HistoryItem) {
	if ui.outputMode != OutputAll && ui.outputMode != OutputInline {
		return
	}
	usePager := ui.pagerEnabled
	hlen := hi.Hosts.maxLen()
	linecount := 0
	buffer := ""
	var pgr *pager
	if usePager {
		buffer = ui.formatter.formatCommand(hi.Command)
		linecount = 1
	} else {
		ui.pchan <- ui.formatter.formatCommand(hi.Command)
	}

	for _, h := range hi.Hosts {
		var txt string
		if ui.outputMode == OutputAll {
			txt = ui.formatter.formatResult(hi.Results[h.Name])
		} else {
			txt = ui.formatter.formatOutput(hi.Results[h.Name], hlen)
		}
		if !usePager {
			ui.pchan <- txt
		} else if pgr != nil {
			pgr.WriteString(txt)
		} else {
			buffer += txt
			linecount += strings.Count(txt, "\n")
			if linecount > ui.height {
				pgr = &pager{}
				if err := pgr.start(); err != nil {
					logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
					ui.pchan <- buffer
					usePager = false
				} else {
					defer pgr.Wait()
					pgr.WriteString(buffer)
				}
				buffer = ""
			}
		}
	}
	if buffer != "" {
		ui.pchan <- buffer
	}
}

func (ui *SimpleUI) PrintHostList(hosts Hosts, oneline, csvOutput, allAttributes bool, attributes []string) {
	if oneline {
		names := make([]string, len(hosts))
		for i, host := range hosts {
			names[i] = host.Name
		}
		ui.pchan <- strings.Join(names, ",")
		return
	}
	var pgr *pager
	if ui.pagerEnabled && len(hosts) > ui.height {
		pgr = &pager{}
		if err := pgr.start(); err != nil {
			logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
			pgr = nil
		} else {
			defer pgr.Wait()
		}
	}
	if allAttributes || len(attributes) > 0 {
		var writer datawriter
		var out io.Writer = ui
		if pgr != nil {
			out = pgr
		}
		if csvOutput {
			writer = csv.NewWriter(out)
		} else {
			writer = newColumnizer(out, "   ")
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
			if pgr != nil {
				pgr.WriteString(host.Name + "\n")
			} else {
				ui.pchan <- host.Name + "\n"
			}
		}
	}
}

func (ui *SimpleUI) CacheUpdateChannel() chan CacheMessage {
	mc := make(chan CacheMessage)
	go func() {
		start := time.Now()
		cached := false
		ticker := time.NewTicker(time.Second / 2)
		defer ticker.Stop()
		caches := make([]string, 0)
		for {
			select {
			case msg, ok := <-mc:
				// Cache message channel closed, we're done caching
				if !ok {
					if cached {
						ui.pchan <- fmt.Sprintf("\r\033[2KAll caches updated\n")
					}
					return
				}
				if msg.Err != nil {
					logrus.Errorf("Error contacting %s: %s", msg.Name, msg.Err)
				}
				if msg.Finished {
					logrus.Debugf("Cache updated for %s", msg.Name)
					for i, v := range caches {
						if v == msg.Name {
							caches = append(caches[:i], caches[i+1:]...)
							break
						}
					}
				} else {
					caches = append(caches, msg.Name)
				}
			case <-ticker.C:
			}
			if len(caches) > 0 {
				since := time.Since(start).Truncate(time.Second)
				cs := strings.Join(caches, ", ")
				if len(cs) > ui.width-25 {
					cs = cs[:ui.width-30] + "..."
				}
				ui.pchan <- fmt.Sprintf("\r\033[2K%s Refreshing caches %s", since, ansi.Color(cs, "green"))
				cached = true
			}
		}
	}()
	return mc
}

func (ui *SimpleUI) OutputChannel(r *Runner) chan OutputLine {
	if ui.outputMode != OutputTail {
		return nil
	}
	oc := make(chan OutputLine)
	hlen := r.hosts.maxLen()
	lastcolor := []byte{}
	reset := []byte("\033[0m")
	cr := regexp.MustCompile("\033\\[[0-9;]+m")
	go func() {
		for msg := range oc {
			name := fmt.Sprintf("%-*s", hlen, msg.Host.Name)
			if msg.Stderr {
				name = ansi.Color(name, "red")
			}
			line := msg.Data
			suffix := []byte{}
			colors := cr.FindAll(line, -1)
			if colors != nil && !bytes.Equal(colors[len(colors)-1], reset) {
				lastcolor = colors[len(colors)-1]
				suffix = reset
			}
			if colors == nil && len(lastcolor) != 0 {
				suffix = reset
			}
			ui.pchan <- fmt.Sprintf("%s  %s%s%s", name, lastcolor, line, suffix)
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
		total := len(r.hosts)
		queued, todo := total, total
		nok, nfail, nerr := 0, 0, 0
		hlen := r.hosts.maxLen()
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
				if ui.outputMode == OutputPerhost {
					ui.pchan <- ui.formatter.formatResult(msg.Result)
				} else if ui.outputMode == OutputTail {
					ui.pchan <- ui.formatter.formatStatus(msg.Result, hlen)
				}
				todo--
			}
			since := time.Since(start).Truncate(time.Second)
			togo := r.timeout - since
			if todo == 0 {
				ui.pchan <- fmt.Sprintf("\r\033[2K%d done, %d ok, %d fail, %d error in %s\n", total, nok, nfail, nerr, since)
			} else if queued >= 0 {
				ui.pchan <- fmt.Sprintf("\r\033[2KWaiting (%s/%s)... %d/%d done, %d queued, %d ok, %d fail, %d error", since, togo, total-todo, total, queued, nok, nfail, nerr)
			} else {
				ui.pchan <- fmt.Sprintf("\r\033[2KWaiting (%s/%s)... %d/%d done, %d ok, %d fail, %d error", since, togo, total-todo, total, nok, nfail, nerr)
			}
		}
	}()
	return pc
}
