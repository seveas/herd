package herd

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
	"github.com/seveas/readline"
	"github.com/seveas/variance"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type OutputMode int

type LoadingMessage func(string, bool, error)

const clearLine string = "\r\033[2K"

// A 256 bit bitmap mapping the 256 possible background colors to light (0) or dark (1)
var darkBackground [4]uint64 = [4]uint64{
	0xfff0_0003_ffff_00ff,
	0xf000_3fff_ff00_003f,
	0x0003_ffff_0000_3fff,
	0x000f_ff00_003f_fff0,
}

func isBackgroundDark() bool {
	color := os.Getenv("COLORFGBG")
	if color == "" {
		return false
	}
	parts := strings.Split(color, ";")
	if len(parts) != 2 {
		return false
	}
	bg, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return false
	}
	if bg > 255 {
		return false
	}
	// The 256 bit bitmap is split into 4 64 bit integers, so we need to shift and mask
	return darkBackground[bg>>6]&(1<<(bg&0x3f)) != 0
}

const (
	OutputTail OutputMode = iota
	OutputPerhost
	OutputInline
	OutputAll
)

var outputModeString map[OutputMode]string = map[OutputMode]string{
	OutputTail:    "tail",
	OutputPerhost: "per-host",
	OutputInline:  "inline",
	OutputAll:     "all",
}

type ColorConfig struct {
	LogDebug   string
	LogInfo    string
	LogWarn    string
	LogError   string
	Command    string
	Summary    string
	Provider   string
	HostStdout string
	HostStderr string
	HostOK     string
	HostFail   string
	HostError  string
	HostCancel string
}

func (c ColorConfig) Merge(other ColorConfig) ColorConfig {
	if c.LogDebug == "" {
		c.LogDebug = other.LogDebug
	}
	if c.LogInfo == "" {
		c.LogInfo = other.LogInfo
	}
	if c.LogWarn == "" {
		c.LogWarn = other.LogWarn
	}
	if c.LogError == "" {
		c.LogError = other.LogError
	}
	if c.Command == "" {
		c.Command = other.Command
	}
	if c.Summary == "" {
		c.Summary = other.Summary
	}
	if c.Provider == "" {
		c.Provider = other.Provider
	}
	if c.HostStdout == "" {
		c.HostStdout = other.HostStdout
	}
	if c.HostStderr == "" {
		c.HostStderr = other.HostStderr
	}
	if c.HostOK == "" {
		c.HostOK = other.HostOK
	}
	if c.HostFail == "" {
		c.HostFail = other.HostFail
	}
	if c.HostError == "" {
		c.HostError = other.HostError
	}
	if c.HostCancel == "" {
		c.HostCancel = other.HostCancel
	}
	return c
}

// FIXME: this is just a copy of the dark theme
var defaultColorConfigLight = ColorConfig{
	LogDebug:   "black+h",
	LogInfo:    "default",
	LogWarn:    "yellow",
	LogError:   "red+b",
	Command:    "cyan",
	Summary:    "black+h",
	Provider:   "green",
	HostStdout: "default",
	HostStderr: "yellow",
	HostOK:     "green",
	HostFail:   "yellow",
	HostError:  "red",
	HostCancel: "black+h",
}

var defaultColorConfigDark = ColorConfig{
	LogDebug:   "black+h",
	LogInfo:    "default",
	LogWarn:    "yellow",
	LogError:   "red+b",
	Command:    "cyan",
	Summary:    "black+h",
	Provider:   "green",
	HostStdout: "default",
	HostStderr: "yellow",
	HostOK:     "green",
	HostFail:   "yellow",
	HostError:  "red",
	HostCancel: "black+h",
}

type SettingsFunc func() (string, map[string]interface{})

type UI interface {
	PrintHistoryItem(hi *HistoryItem)
	PrintHostList(opts HostListOptions)
	PrintSettings(...SettingsFunc)
	SetOutputMode(OutputMode)
	SetOutputTimestamp(bool)
	SetPagerEnabled(bool)
	Sync()
	End()
	LoadingMessage(what string, done bool, err error)
	OutputChannel() chan OutputLine
	ProgressChannel(deadline time.Time) chan ProgressMessage
	BindLogrus()
	Settings() (string, map[string]interface{})
}

type HostListOptions struct {
	OneLine       bool
	Separator     string
	Csv           bool
	Attributes    []string
	AllAttributes bool
	Align         bool
	Header        bool
	Template      string
	Count         []string
	SortByCount   bool
	Group         string
}

type SimpleUI struct {
	hosts           *HostSet
	colors          ColorConfig
	output          *os.File
	altOutput       *os.File
	lastProgress    string
	pchan           chan outputMessage
	formatter       formatter
	outputMode      OutputMode
	outputTimestamp bool
	pagerEnabled    bool
	width           int
	height          int
	syncCond        *sync.Cond
	loading         []string
	loadStart       time.Time
	loadOnce        sync.Once
	loadLock        sync.Mutex
	loadTicker      *time.Ticker
}

type outputMessageType int

const (
	outputMessageFlush outputMessageType = iota
	outputMessageLog
	outputMessageCommandOutput
	outputMessageResult
	outputMessageProgress
	outputMessageHostlist
)

type outputMessage struct {
	messageType outputMessageType
	message     string
}

var templateFuncs = template.FuncMap{
	"yaml": func(data interface{}) (string, error) {
		b, err := yaml.Marshal(data)
		return "---\n" + string(b), err
	},
	"sshkey": func(data interface{}) (string, error) {
		key, ok := data.(ssh.PublicKey)
		if !ok {
			return "", fmt.Errorf("sshkey only knows how to show ssh keys")
		}
		k := key.Marshal()
		b := make([]byte, base64.StdEncoding.EncodedLen(len(k)))
		base64.StdEncoding.Encode(b, k)
		return fmt.Sprintf("%s %s", key.Type(), string(b)), nil
	},
}

func NewSimpleUI(colors ColorConfig, hosts *HostSet) *SimpleUI {
	defaultColorConfig := defaultColorConfigLight
	if isBackgroundDark() {
		defaultColorConfig = defaultColorConfigDark
	}
	colors = colors.Merge(defaultColorConfig)

	f := newPrettyFormatter(colors)
	ui := &SimpleUI{
		colors:       colors,
		hosts:        hosts,
		output:       os.Stdout,
		altOutput:    os.Stderr,
		outputMode:   OutputAll,
		lastProgress: "",
		pchan:        make(chan outputMessage, 100),
		syncCond:     &sync.Cond{L: new(sync.Mutex)},
		formatter:    f,
	}
	if isatty.IsTerminal(ui.output.Fd()) {
		ui.altOutput = ui.output
	}
	ui.getSize()
	readline.DefaultOnWidthChanged(func() {
		ui.getSize()
	})
	go ui.printer()
	return ui
}

func (ui *SimpleUI) getSize() {
	w, h, err := readline.GetSize(int(ui.output.Fd())) // nolint:gosec // FD's shouldn'tt reach higher than 2^31-1
	if err != nil {
		ui.pagerEnabled = false
		ansi.DisableColors(true)
		w, h, err = readline.GetSize(int(ui.altOutput.Fd())) // nolint:gosec // FD's shouldn'tt reach higher than 2^31-1
	}
	if err == nil {
		ui.width, ui.height = w, h
		if w < 40 {
			ui.width = 40
		}
	} else {
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
		if msg.messageType == outputMessageFlush {
			ui.syncCond.L.Lock()
			ui.syncCond.Broadcast()
			ui.syncCond.L.Unlock()
			continue
		}

		out := ui.output
		if ui.output != ui.altOutput && (msg.messageType == outputMessageLog || msg.messageType == outputMessageProgress) {
			out = ui.altOutput
		}

		// If we're getting a normal message in the middle of printing
		// progress, wipe the progress message and reprint it after this
		// message
		if ui.lastProgress != "" && msg.messageType != outputMessageProgress && out == ui.altOutput {
			out.WriteString(clearLine + msg.message + ui.lastProgress)
			out.Sync()
			continue
		}

		if msg.messageType == outputMessageProgress {
			out.WriteString(clearLine)
			ui.lastProgress = msg.message
		}

		out.WriteString(msg.message)
		out.Sync()
	}
}

type simpleUIWriter struct {
	messageType outputMessageType
	lineBuf     string
	ui          *SimpleUI
}

func (w *simpleUIWriter) Write(msg []byte) (int, error) {
	w.lineBuf += string(msg)
	for i := strings.IndexRune(w.lineBuf, '\n'); i != -1; i = strings.IndexRune(w.lineBuf, '\n') {
		w.ui.pchan <- outputMessage{w.messageType, w.lineBuf[:i+1]}
		w.lineBuf = w.lineBuf[i+1:]
	}
	return len(msg), nil
}

func (ui *SimpleUI) BindLogrus() {
	logrus.SetFormatter(ui.formatter)
	logrus.SetOutput(&simpleUIWriter{outputMessageLog, "", ui})
}

func (ui *SimpleUI) Sync() {
	ui.syncCond.L.Lock()
	ui.pchan <- outputMessage{outputMessageFlush, ""}
	defer ui.syncCond.L.Unlock()
	ui.syncCond.Wait()
}

func (ui *SimpleUI) End() {
	ui.Sync()
	close(ui.pchan)
}

func (ui *SimpleUI) PrintHistoryItem(hi *HistoryItem) {
	if ui.outputMode != OutputAll && ui.outputMode != OutputInline {
		return
	}
	usePager := ui.pagerEnabled
	hlen := hi.maxHostNameLength
	linecount := 0
	buffer := ""
	var pgr *pager
	if usePager {
		linecount = 2
	}

	for _, result := range hi.Results {
		var txt string
		if ui.outputMode == OutputAll {
			txt = ui.formatter.formatResult(result, hlen)
		} else {
			txt = ui.formatter.formatOutput(result, hlen)
		}
		if !usePager {
			ui.pchan <- outputMessage{outputMessageResult, txt}
		} else if pgr != nil {
			pgr.WriteString(txt)
		} else {
			buffer += txt
			linecount += strings.Count(txt, "\n")
			if linecount > ui.height {
				pgr = &pager{}
				if err := pgr.start(); err != nil {
					logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
					ui.pchan <- outputMessage{outputMessageResult, buffer}
					usePager = false
				} else {
					pgr.WriteString(ui.formatter.formatCommand(hi.Command))
					pgr.WriteString(ui.formatter.formatSummary(hi.Summary.Ok, hi.Summary.Fail, hi.Summary.Err))
					pgr.WriteString(buffer)
				}
				buffer = ""
			}
		}
	}
	if buffer != "" {
		ui.pchan <- outputMessage{outputMessageResult, buffer}
	}
	if pgr != nil {
		if err := pgr.Wait(); err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				logrus.Warnf("Pager error: %s", err)
			}
		}
	}
}

func startPager(p *pager, o *io.Writer) {
	if p == nil {
		return
	}
	if err := p.start(); err != nil {
		logrus.Warnf("Unable to start pager, displaying on stdout: %s", err)
	} else {
		*o = p
	}
}

func (ui *SimpleUI) PrintHostList(opts HostListOptions) {
	hosts := ui.hosts.hosts
	if len(opts.Count) == 1 && opts.Count[0] == "*" {
		ui.pchan <- outputMessage{outputMessageHostlist, fmt.Sprintf("%d\n", len(hosts))}
		return
	}
	if len(hosts) == 0 {
		logrus.Error("No hosts to list")
		return
	}
	if opts.OneLine {
		names := make([]string, len(hosts))
		for i, host := range hosts {
			names[i] = host.Name
		}
		ui.pchan <- outputMessage{outputMessageHostlist, strings.Join(names, opts.Separator) + "\n"}
		return
	}
	var writer datawriter
	var out io.Writer = &simpleUIWriter{outputMessageHostlist, "", ui}
	pgr := &pager{}
	if !ui.pagerEnabled {
		pgr = nil
	}
	if len(opts.Attributes) > 0 {
		if len(hosts) > ui.height {
			startPager(pgr, &out)
		}
		if opts.Csv {
			writer = csv.NewWriter(out)
		} else if opts.Align {
			writer = newColumnizer(out, "   ")
		} else {
			writer = newPassthrough(out)
		}
		allAttrs := make(map[string]bool)
		attrs := make([]string, 0, len(opts.Attributes))
		for _, attr := range opts.Attributes {
			if !strings.ContainsRune(attr, '*') {
				attrs = append(attrs, attr)
				continue
			}
			if len(allAttrs) == 0 {
				for _, host := range hosts {
					for key := range host.Attributes {
						allAttrs[key] = true
					}
				}
			}
			myattrs := make([]string, 0)
			for key := range allAttrs {
				if match, _ := filepath.Match(attr, key); match {
					myattrs = append(myattrs, key)
				}
			}
			sort.Strings(myattrs)
			attrs = append(attrs, myattrs...)
		}
		opts.Attributes = attrs
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
		if w, ok := writer.(*columnizer); ok && w.width > ui.width {
			startPager(pgr, &w.output)
		}
		writer.Flush()
	} else if opts.Template != "" {
		if len(hosts) > ui.height {
			startPager(pgr, &out)
		}
		tmpl, err := template.New("host").Funcs(templateFuncs).Parse(opts.Template + "\n")
		if err != nil {
			logrus.Errorf("Unable to parse template '%s': %s", opts.Template, err)
			return
		}
		for _, host := range hosts {
			err := tmpl.Execute(out, host)
			if err != nil {
				logrus.Errorf("Error executing template: %s", err)
			}
		}
	} else if len(opts.Count) != 0 {
		// First we generate the counts
		valueKeys := make([]string, 0)
		values := make(map[string][]string)
		counts := make(map[string]int)
		groupCounts := make(map[string]map[any]int)
		groupValues := make(map[string]bool)
		totals := make(map[string]int)
		total := 0

		for _, host := range hosts {
			total++
			v := make([]string, len(opts.Count))
			for i, attr := range opts.Count {
				v[i] = fmt.Sprintf("%v", host.Attributes[attr])
			}
			vs := strings.Join(v, "\000")
			if _, ok := counts[vs]; ok {
				counts[vs]++
			} else {
				values[vs] = v
				counts[vs] = 1
				valueKeys = append(valueKeys, vs)
			}
			if opts.Group != "" {
				if _, ok := groupCounts[vs]; !ok {
					groupCounts[vs] = make(map[any]int)
				}
				attr := fmt.Sprintf("%v", host.Attributes[opts.Group])
				groupCounts[vs][attr]++
				groupValues[attr] = true
				totals[attr]++
			}
		}
		if opts.SortByCount {
			// We sort by count, keeping the order of entries with the same count intact
			positions := make(map[string]int)
			for i, k := range valueKeys {
				positions[k] = i
			}
			sort.Slice(valueKeys, func(i, j int) bool {
				if counts[valueKeys[i]] == counts[valueKeys[j]] {
					return positions[valueKeys[i]] < positions[valueKeys[j]]
				}
				return counts[valueKeys[i]] > counts[valueKeys[j]]
			})
		}
		groups := make([]string, 0, len(groupValues))
		for k := range groupValues {
			groups = append(groups, k)
		}
		sort.Strings(groups)

		// And now we write
		if len(valueKeys) > ui.height {
			startPager(pgr, &out)
		}
		if opts.Csv {
			writer = csv.NewWriter(out)
		} else if opts.Align {
			writer = newColumnizer(out, "   ")
		} else {
			writer = newPassthrough(out)
		}
		if opts.Header {
			attrline := make([]string, len(opts.Count)+len(groups), len(opts.Count)+len(groups)+3)
			copy(attrline, opts.Count)
			copy(attrline[len(opts.Count):], groups)
			if len(groups) == 0 {
				attrline = append(attrline, "count", "percentage")
			} else {
				attrline = append(attrline, "total", "percentage", "average", "stddev")
			}
			writer.Write(attrline)
		}
		for _, k := range valueKeys {
			v := values[k]
			stats := variance.New()
			for _, group := range groups {
				v = append(v, fmt.Sprintf("%v", groupCounts[k][group]))
				stats.Add(float64(groupCounts[k][group]))
			}
			v = append(v, fmt.Sprintf("%v", counts[k]), fmt.Sprintf("%v%%", f2s(float64(counts[k])/float64(total)*100)))
			if stats.NumDataValues() != 0 {
				v = append(v, f2s(stats.Mean()), f2s(stats.StandardDeviationPopulation()))
			}

			writer.Write(v)
		}
		if opts.Header {
			v := []string{"total"}
			for range len(opts.Count) - 1 {
				v = append(v, "")
			}
			if opts.Group != "" {
				stats := variance.New()
				for _, group := range groups {
					v = append(v, fmt.Sprintf("%v", totals[group]))
					stats.Add(float64(totals[group]))
				}
				v = append(v, fmt.Sprintf("%v", total), "", f2s(stats.Mean()), f2s(stats.StandardDeviationPopulation()))
			} else {
				v = append(v, fmt.Sprintf("%v", total), "")
			}
			writer.Write(v)
		}
		// Start the pager after all if we are getting too wide
		if w, ok := writer.(*columnizer); ok && w.width > ui.width {
			startPager(pgr, &w.output)
		}
		writer.Flush()
	} else {
		if len(hosts) > ui.height {
			startPager(pgr, &out)
		}
		for _, host := range hosts {
			fmt.Fprintln(out, host.Name)
		}
	}
	if err := pgr.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			logrus.Warnf("Pager error: %s", err)
		}
	}
}

// Would love to use %g, but I always want %f style ouptut, just without trailing zeroes
// strconv.FormatFloat wit precision -1 comes close too, but can have more than 2 digits of precision
func f2s(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", f), "0"), ".")
}

func (ui *SimpleUI) LoadingMessage(what string, done bool, err error) {
	if !logrus.IsLevelEnabled(logrus.InfoLevel) {
		return
	}

	ui.loadLock.Lock()
	defer ui.loadLock.Unlock()
	if what == "" && done {
		if ui.loadTicker != nil {
			ui.pchan <- outputMessage{outputMessageProgress, ""}
			ui.loadTicker.Stop()
		}
		return
	}
	ui.loadOnce.Do(func() {
		ui.loadStart = time.Now()
		ui.loading = []string{}
		ui.loadTicker = time.NewTicker(time.Second / 2)
		go func() {
			for {
				<-ui.loadTicker.C
				ui.LoadingMessage("", false, nil)
			}
		}()
	})
	if err != nil {
		logrus.Errorf("Error loading data from %s: %s", what, err)
	}
	if done {
		logrus.Debugf("Done loading %s", what)
		for i, v := range ui.loading {
			if v == what {
				ui.loading = append(ui.loading[:i], ui.loading[i+1:]...)
				break
			}
		}
		if len(ui.loading) == 0 {
			return
		}
	} else if what != "" {
		ui.loading = append(ui.loading, what)
	}

	since := time.Since(ui.loadStart).Truncate(time.Second)
	cs := strings.Join(ui.loading, ", ")
	if len(cs) > ui.width-25 {
		cs = cs[:ui.width-30] + "..."
	}
	ui.pchan <- outputMessage{outputMessageProgress, fmt.Sprintf("%s Loading data %s", since, ansi.Color(cs, ui.colors.Provider))}
}

func (ui *SimpleUI) OutputChannel() chan OutputLine {
	if ui.outputMode != OutputTail {
		return nil
	}
	oc := make(chan OutputLine)
	hlen := ui.hosts.maxNameLength
	lastcolor := []byte{}
	reset := []byte("\033[0m")
	cr := regexp.MustCompile("\033\\[[0-9;]+m")
	ts := ""
	go func() {
		for msg := range oc {
			if ui.outputTimestamp {
				ts = time.Now().Format("15:04:05.000 ")
			}
			name := fmt.Sprintf("%-*s", hlen, msg.Host.Name)
			if msg.Stderr {
				name = ansi.Color(name, ui.colors.HostStderr)
			}
			line := msg.Data
			// Strip all text up to the last embedded \r, unless that is part of a \r\n
			for bytes.HasSuffix(line, []byte("\r\n")) {
				line = line[:len(line)-1]
				line[len(line)-1] = '\n'
			}
			if idx := bytes.LastIndex(line, []byte("\r")); idx != -1 {
				line = line[idx+1:]
			}
			// Make sure we don't pollute hostnames with colors
			suffix := []byte{}
			colors := cr.FindAll(line, -1)
			if len(colors) > 0 && !bytes.Equal(colors[len(colors)-1], reset) {
				lastcolor = colors[len(colors)-1]
				suffix = reset
			}
			if colors == nil && len(lastcolor) != 0 {
				suffix = reset
			}
			ui.pchan <- outputMessage{outputMessageCommandOutput, fmt.Sprintf("%s%s  %s%s%s", ts, name, lastcolor, line, suffix)}
		}
	}()
	return oc
}

func (ui *SimpleUI) ProgressChannel(deadline time.Time) chan ProgressMessage {
	pc := make(chan ProgressMessage)
	go func() {
		start := time.Now()
		ticker := time.NewTicker(time.Second / 2)
		defer ticker.Stop()
		defer func() {
			ui.pchan <- outputMessage{outputMessageProgress, ""}
		}()
		total := len(ui.hosts.hosts)
		queued, todo, waiting, running, done := total, total, 0, 0, 0
		nok, nfail, nerr := 0, 0, 0
		hlen := ui.hosts.maxNameLength
		show_waiting := false
		for {
			select {
			case <-ticker.C:
			case msg, ok := <-pc:
				if !ok {
					return
				}
				switch msg.State {
				case Waiting:
					show_waiting = true
					queued--
					waiting++
				case Running:
					if show_waiting {
						waiting--
					} else {
						queued--
					}
					running++
				case Finished:
					running--
					todo--
					done++
					switch msg.Result.ExitStatus {
					case -1:
						nerr++
					case 0:
						nok++
					default:
						nfail++
					}
					if ui.outputMode == OutputPerhost {
						ui.pchan <- outputMessage{outputMessageResult, ui.formatter.formatResult(msg.Result, hlen)}
					} else if ui.outputMode == OutputTail {
						status := ui.formatter.formatStatus(msg.Result, hlen)
						if ui.outputTimestamp {
							status = msg.Result.EndTime.Format("15:04:05.000 ") + status
						}
						ui.pchan <- outputMessage{outputMessageResult, status}
					}
				}
			}
			since := (time.Since(start) + time.Second/2).Truncate(time.Second)
			if todo == 0 {
				ui.pchan <- outputMessage{outputMessageProgress, fmt.Sprintf("%d done, %d ok, %d fail, %d error in %s\n", total, nok, nfail, nerr, since)}
			} else {
				togo := (time.Until(deadline) + time.Second/2).Truncate(time.Second)
				msg := fmt.Sprintf("Waiting (%s/%s)... %d/%d done", since, togo, done, total)
				if queued > 0 {
					msg += fmt.Sprintf(", %d queued", queued)
				}
				if show_waiting {
					msg += fmt.Sprintf(", %d waiting", waiting)
				}
				msg += fmt.Sprintf(", %d in progress, %d ok, %d fail, %d error", running, nok, nfail, nerr)
				ui.pchan <- outputMessage{outputMessageProgress, msg}
			}
		}
	}()
	return pc
}

func (ui *SimpleUI) PrintSettings(funcs ...SettingsFunc) {
	for _, f := range funcs {
		name, settings := f()
		ui.pchan <- outputMessage{outputMessageLog, name + "\n"}
		l := 0
		keys := make([]string, 0, len(settings))
		for k := range settings {
			keys = append(keys, k)
			if len(k) > l {
				l = len(k)
			}
		}
		l += 5
		sort.Strings(keys)
		for _, k := range keys {
			ui.pchan <- outputMessage{outputMessageLog, fmt.Sprintf("    %-*s %v\n", l, k+":", settings[k])}
		}
	}
}

func (ui *SimpleUI) Settings() (string, map[string]interface{}) {
	return "User Interface", map[string]interface{}{
		"Type":      "Simple",
		"Output":    outputModeString[ui.outputMode],
		"Timestamp": ui.outputTimestamp,
		"NoPager":   !ui.pagerEnabled,
		"NoColor":   ansi.Black == "",
	}
}
