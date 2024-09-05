package common

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

var levelMap = map[hclog.Level]logrus.Level{
	hclog.NoLevel: logrus.InfoLevel,
	hclog.Trace:   logrus.TraceLevel,
	hclog.Debug:   logrus.DebugLevel,
	hclog.Info:    logrus.InfoLevel,
	hclog.Warn:    logrus.DebugLevel,
	hclog.Error:   logrus.ErrorLevel,
	// hclog.Off:     logrus.FatalLevel,
}

var levelReverseMap = map[logrus.Level]hclog.Level{
	logrus.TraceLevel: hclog.Trace,
	logrus.DebugLevel: hclog.Debug,
	logrus.InfoLevel:  hclog.Info,
	logrus.WarnLevel:  hclog.Warn,
	logrus.ErrorLevel: hclog.Error,
	logrus.FatalLevel: hclog.Error,
	logrus.PanicLevel: hclog.Error,
}

type logrusLogger struct {
	logger *logrus.Logger
	name   string
	args   []interface{}
}

func NewLogrusLogger(l *logrus.Logger, name string) *logrusLogger {
	return &logrusLogger{logger: l, name: name}
}

func (l *logrusLogger) format(msg string, args ...interface{}) string {
	return fmt.Sprintf("%s: %s %v", l.name, msg, args)
}

func (l *logrusLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	l.logger.Log(levelMap[level], l.format(msg, args...))
}

func (l *logrusLogger) Trace(msg string, args ...interface{}) {
	l.logger.Trace(l.format(msg, args...))
}

func (l *logrusLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(l.format(msg, args...))
}

func (l *logrusLogger) Info(msg string, args ...interface{}) {
	// Downgrade the severity of this message, unlike what
	// https://github.com/hashicorp/go-plugin/pull/195/files says, these
	// messages are not valuable to us.
	if strings.HasPrefix(msg, "plugin process exited") {
		l.Debug(msg, args...)
		return
	}
	l.logger.Info(l.format(msg, args...))
}

func (l *logrusLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(l.format(msg, args...))
}

func (l *logrusLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(l.format(msg, args...))
}

func (l *logrusLogger) IsTrace() bool {
	return l.logger.IsLevelEnabled(logrus.TraceLevel)
}

func (l *logrusLogger) IsDebug() bool {
	return l.logger.IsLevelEnabled(logrus.DebugLevel)
}

func (l *logrusLogger) IsInfo() bool {
	return l.logger.IsLevelEnabled(logrus.InfoLevel)
}

func (l *logrusLogger) IsWarn() bool {
	return l.logger.IsLevelEnabled(logrus.WarnLevel)
}

func (l *logrusLogger) IsError() bool {
	return l.logger.IsLevelEnabled(logrus.ErrorLevel)
}

func (l *logrusLogger) ImpliedArgs() []interface{} {
	return l.args
}

func (l *logrusLogger) With(args ...interface{}) hclog.Logger {
	return &logrusLogger{logger: l.logger, name: l.name, args: args}
}

func (l *logrusLogger) Name() string {
	return l.name
}

func (l *logrusLogger) Named(name string) hclog.Logger {
	if l.name != "" {
		name = fmt.Sprintf("%s: %s", l.name, name)
	}
	return &logrusLogger{logger: l.logger, name: name, args: l.args}
}

func (l *logrusLogger) ResetNamed(name string) hclog.Logger {
	return &logrusLogger{logger: l.logger, name: name, args: l.args}
}

func (l *logrusLogger) SetLevel(level hclog.Level) {
	l.logger.SetLevel(levelMap[level])
}

func (l *logrusLogger) GetLevel() (level hclog.Level) {
	return levelReverseMap[l.logger.GetLevel()]
}

func (l *logrusLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	panic("we don't support StandardLogger")
}

func (l *logrusLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return l.logger.Writer()
}

var _ hclog.Logger = &logrusLogger{logger: logrus.StandardLogger()}
