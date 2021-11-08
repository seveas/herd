package herd

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

type formatter interface {
	formatCommand(c string) string
	formatSummary(ok, fail, err int) string
	formatResult(r *Result, l int) string
	formatStatus(r *Result, l int) string
	formatOutput(r *Result, l int) string
	Format(e *logrus.Entry) ([]byte, error)
}

type prettyFormatter struct {
	colors map[logrus.Level]string
}

func (f prettyFormatter) formatCommand(command string) string {
	return ansi.Color(command, "cyan") + "\n"
}

func (f prettyFormatter) formatSummary(ok, fail, err int) string {
	return ansi.Color(fmt.Sprintf("%d ok, %d fail, %d error", ok, fail, err), "black+h") + "\n"
}

func (f prettyFormatter) formatResult(r *Result, l int) string {
	out := f.formatStatus(r, l)
	if len(r.Stdout) > 0 {
		out += f.indent(string(r.Stdout), "    ", "    ")
	}
	if len(r.Stderr) != 0 {
		out += ansi.Color("----", "black+h") + "\n" + f.indent(string(r.Stderr), "    ", "    ")
	}
	return out
}

func (f prettyFormatter) formatOutput(r *Result, l int) string {
	prefix := fmt.Sprintf("%-*s  ", l, r.Host.Name)
	indent := fmt.Sprintf("%-*s  ", l, "")
	out := ""
	if len(r.Stdout) > 0 {
		out += f.indent(string(r.Stdout), prefix, indent)
	}
	if len(r.Stderr) > 0 {
		out += f.indent(string(r.Stderr), ansi.Color(prefix, "red"), indent)
	}
	if out == "" || r.Err != nil {
		out += f.formatStatus(r, l)
	}
	return out
}

func (f prettyFormatter) formatStatus(r *Result, l int) string {
	if r.Err != nil {
		return ansi.Color(fmt.Sprintf("%-*s  %s after %s", l, r.Host.Name, r.Err, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), "red") + "\n"
	} else {
		return ansi.Color(fmt.Sprintf("%-*s  completed successfully after %s", l, r.Host.Name, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), "green") + "\n"
	}
}

func (f prettyFormatter) indent(msg, prefix, indent string) string {
	return prefix + strings.ReplaceAll(strings.TrimSuffix(msg, "\n"), "\n", "\n"+indent) + "\n"
}

func (f prettyFormatter) Format(e *logrus.Entry) ([]byte, error) {
	msg := e.Message
	if color, ok := f.colors[e.Level]; ok {
		msg = ansi.Color(msg, color)
	}
	return []byte(msg + "\n"), nil
}
