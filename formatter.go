package herd

import (
	"context"
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
	colors       ColorConfig
	logrusColors map[logrus.Level]string
}

func newPrettyFormatter(colors ColorConfig) prettyFormatter {
	return prettyFormatter{
		colors: colors,
		logrusColors: map[logrus.Level]string{
			// logrus.PanicLevel: colors.LogPanic,
			// logrus.FatalLevel: colors.LogFatal,
			logrus.ErrorLevel: colors.LogError,
			logrus.WarnLevel:  colors.LogWarn,
			logrus.InfoLevel:  colors.LogInfo,
			logrus.DebugLevel: colors.LogDebug,
			// logrus.TraceLevel: colors.LogTrace,
		},
	}
}

func (f prettyFormatter) formatCommand(command string) string {
	return ansi.Color(command, f.colors.Command) + "\n"
}

func (f prettyFormatter) formatSummary(ok, fail, err int) string {
	return ansi.Color(fmt.Sprintf("%d ok, %d fail, %d error", ok, fail, err), f.colors.Summary) + "\n"
}

func (f prettyFormatter) formatResult(r *Result, l int) string {
	out := f.formatStatus(r, l)
	if len(r.Stdout) > 0 {
		out += f.indent(string(r.Stdout), "    ", "    ")
	}
	if len(r.Stderr) != 0 {
		out += ansi.Color("----", f.colors.Summary) + "\n" + f.indent(string(r.Stderr), "    ", "    ")
	}
	return out
}

func (f prettyFormatter) formatOutput(r *Result, l int) string {
	prefix := fmt.Sprintf("%-*s  ", l, r.Host)
	indent := fmt.Sprintf("%-*s  ", l, "")
	out := ""
	if len(r.Stdout) > 0 {
		prefix := prefix
		if r.Err == nil {
			prefix = ansi.Color(prefix, f.colors.HostOK)
		} else if r.ExitStatus != -1 {
			prefix = ansi.Color(prefix, f.colors.HostFail)
		} else if r.Err.Error() == context.Canceled.Error() {
			prefix = ansi.Color(prefix, f.colors.HostCancel)
		} else {
			prefix = ansi.Color(prefix, f.colors.HostError)
		}
		out += f.indent(string(r.Stdout), prefix, indent)
	}
	if len(r.Stderr) > 0 {
		out += f.indent(string(r.Stderr), ansi.Color(prefix, f.colors.HostStderr), indent)
	}
	if out == "" || r.Err != nil {
		out += f.formatStatus(r, l)
	}
	return out
}

func (f prettyFormatter) formatStatus(r *Result, l int) string {
	if r.Err == nil {
		return ansi.Color(fmt.Sprintf("%-*s  completed successfully after %s", l, r.Host, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), f.colors.HostOK) + "\n"
	} else if r.ExitStatus != -1 {
		return ansi.Color(fmt.Sprintf("%-*s  exited with status %d after %s", l, r.Host, r.ExitStatus, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), f.colors.HostFail) + "\n"
	} else if r.Err.Error() == context.Canceled.Error() {
		return ansi.Color(fmt.Sprintf("%-*s  skipped due to global timeout", l, r.Host), f.colors.HostCancel) + "\n"
	} else {
		return ansi.Color(fmt.Sprintf("%-*s  %s after %s", l, r.Host, r.Err, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), f.colors.HostError) + "\n"
	}
}

func (f prettyFormatter) indent(msg, prefix, indent string) string {
	return prefix + strings.ReplaceAll(strings.TrimSuffix(msg, "\n"), "\n", "\n"+indent) + "\n"
}

func (f prettyFormatter) Format(e *logrus.Entry) ([]byte, error) {
	msg := e.Message
	if color, ok := f.logrusColors[e.Level]; ok {
		msg = ansi.Color(msg, color)
	}
	return []byte(msg + "\n"), nil
}
