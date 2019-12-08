package katyusha

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

var Formatters = map[string]Formatter{
	"pretty": PrettyFormatter{
		colors: map[logrus.Level]string{
			logrus.WarnLevel:  "yellow",
			logrus.ErrorLevel: "red+b",
			logrus.DebugLevel: "black+h",
		},
	},
}

type Formatter interface {
	FormatCommand(c string) string
	FormatResult(r *Result) string
	FormatStatus(r *Result, l int) string
	FormatOutput(r *Result, l int) string
	Format(e *logrus.Entry) ([]byte, error)
}

type PrettyFormatter struct {
	colors map[logrus.Level]string
}

func (f PrettyFormatter) FormatCommand(command string) string {
	return ansi.Color(command, "cyan") + "\n"
}

func (f PrettyFormatter) FormatResult(r *Result) string {
	out := f.FormatStatus(r, 0)
	if len(r.Stdout) > 0 {
		out += f.indent(string(r.Stdout), "    ", "    ")
	}
	if len(r.Stderr) != 0 {
		out += ansi.Color("----", "black+h") + "\n" + f.indent(string(r.Stderr), "    ", "    ")
	}
	return out
}

func (f PrettyFormatter) FormatOutput(r *Result, l int) string {
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
		out += f.FormatStatus(r, l)
	}
	return out
}

func (f PrettyFormatter) FormatStatus(r *Result, l int) string {
	if r.Err != nil {
		return ansi.Color(fmt.Sprintf("%-*s  %s after %s", l, r.Host.Name, r.Err, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), "red") + "\n"
	} else {
		return ansi.Color(fmt.Sprintf("%-*s  completed successfully after %s", l, r.Host.Name, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), "green") + "\n"
	}
}

func (f PrettyFormatter) indent(msg, prefix, indent string) string {
	return prefix + strings.ReplaceAll(strings.TrimSuffix(msg, "\n"), "\n", "\n"+indent) + "\n"
}

func (f PrettyFormatter) Format(e *logrus.Entry) ([]byte, error) {
	msg := e.Message
	if color, ok := f.colors[e.Level]; ok {
		msg = ansi.Color(msg, color)
	}
	return []byte(msg + "\n"), nil
}
