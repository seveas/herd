package katyusha

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
	"github.com/seveas/readline"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
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
	PrintHostList(hosts Hosts, opts HostListOptions)
	PrintKnownHosts(hosts Hosts)
	SetOutputMode(OutputMode)
	SetOutputTimestamp(bool)
	SetPagerEnabled(bool)
	Write([]byte) (int, error)
	Sync()
	End()
	CacheUpdateChannel() chan CacheMessage
	OutputChannel(r *Runner) chan OutputLine
	ProgressChannel(r *Runner) chan ProgressMessage
	BindLogrus()
}

type HostListOptions struct {
	OneLine       bool
	Separator     string
	Csv           bool
	Attributes    []string
	AllAttributes bool
	Align         bool
	Header        bool
}

type SimpleUI struct {
	output          *os.File
	atStart         bool
	lastProgress    string
	pchan           chan string
	dchan           chan interface{}
	formatter       formatter
	outputMode      OutputMode
	outputTimestamp bool
	pagerEnabled    bool
	width           int
	height          int
	lineBuf         string
	isTerminal      bool
	wg              *sync.WaitGroup
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
		wg:           &sync.WaitGroup{},
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

func (ui *SimpleUI) SetOutputTimestamp(e bool) {
	ui.outputTimestamp = e
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

func (ui *SimpleUI) Sync() {
	ui.wg.Wait()
}

func (ui *SimpleUI) End() {
	ui.wg.Wait()
	ui.wg = &sync.WaitGroup{}
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

func (ui *SimpleUI) PrintHostList(hosts Hosts, opts HostListOptions) {
	if len(hosts) == 0 {
		logrus.Error("No hosts to list")
		return
	}
	if opts.OneLine {
		names := make([]string, len(hosts))
		for i, host := range hosts {
			names[i] = host.Name
		}
		ui.pchan <- strings.Join(names, opts.Separator) + "\n"
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
	if opts.AllAttributes || len(opts.Attributes) > 0 {
		var writer datawriter
		var out io.Writer = ui
		if pgr != nil {
			out = pgr
		}
		if opts.Csv {
			writer = csv.NewWriter(out)
		} else if opts.Align {
			writer = newColumnizer(out, "   ")
		} else {
			writer = newPassthrough(out)
		}
		if opts.AllAttributes {
			attrs := make(map[string]bool)
			for _, host := range hosts {
				for key, _ := range host.Attributes {
					attrs[key] = true
				}
			}
			for attr, _ := range attrs {
				opts.Attributes = append(opts.Attributes, attr)
			}
			sort.Strings(opts.Attributes)
		}
		if opts.Header {
			attrline := make([]string, len(opts.Attributes)+1)
			attrline[0] = "name"
			copy(attrline[1:], opts.Attributes)
			writer.Write(attrline)
		}
		for _, host := range hosts {
			line := make([]string, len(opts.Attributes)+1)
			line[0] = host.Name
			for i, attr := range opts.Attributes {
				val, ok := host.GetAttribute(attr)
				value := ""
				if ok {
					if k, ok := val.(ssh.PublicKey); ok {
						val = fmt.Sprintf("%s %s", k.Type(), base64.StdEncoding.EncodeToString(k.Marshal()))
					}
					value = fmt.Sprintf("%v", val)
				}
				line[i+1] = value
			}
			writer.Write(line)
		}
		// Start the pager after all if we are getting too wide
		if w, ok := writer.(*columnizer); ok && ui.pagerEnabled && pgr == nil {
			sum := 0
			for i := 0; i < len(w.lengths); i++ {
				sum += w.lengths[i]
			}
			if sum > ui.width {
				pgr = &pager{}
				if err := pgr.start(); err != nil {
					logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
					pgr = nil
				} else {
					w.output = pgr
					defer pgr.Wait()
				}
			}
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

func (ui *SimpleUI) PrintKnownHosts(hosts Hosts) {
	for _, host := range hosts {
		if host.sshKey != nil {
			ui.pchan <- fmt.Sprintf("%s %s %s\n", host.Name, host.sshKey.Type(), base64.StdEncoding.EncodeToString(host.sshKey.Marshal()))
		}
	}
}

func (ui *SimpleUI) CacheUpdateChannel() chan CacheMessage {
	mc := make(chan CacheMessage)
	ui.wg.Add(1)
	go func() {
		defer ui.wg.Done()
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
	ui.wg.Add(1)
	hlen := r.hosts.maxLen()
	lastcolor := []byte{}
	reset := []byte("\033[0m")
	cr := regexp.MustCompile("\033\\[[0-9;]+m")
	ts := ""
	go func() {
		defer ui.wg.Done()
		for msg := range oc {
			if ui.outputTimestamp {
				ts = time.Now().Format("15:04:05.000 ")
			}
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
			ui.pchan <- fmt.Sprintf("%s%s  %s%s%s", ts, name, lastcolor, line, suffix)
		}
	}()
	return oc
}

func (ui *SimpleUI) ProgressChannel(r *Runner) chan ProgressMessage {
	pc := make(chan ProgressMessage)
	ui.wg.Add(1)
	go func() {
		defer ui.wg.Done()
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
					status := ui.formatter.formatStatus(msg.Result, hlen)
					if ui.outputTimestamp {
						status = msg.Result.EndTime.Format("15:04:05.000 ") + status
					}
					ui.pchan <- status
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
