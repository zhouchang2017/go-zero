package logx

import (
	"time"
)

// A Logger represents a logger.
type Logger interface {
	// WithRequestId 附加请求id
	WithRequestId(id string) Logger
	// WithError 附加错误
	WithError(err error) Logger
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Traceln(args ...interface{})
	// Slow logs a message at slow level.
	Slow(...interface{})
	// Slowf logs a message at slow level.
	Slowf(string, ...interface{})
	// Slowv logs a message at slow level.
	Slowv(interface{})
	// Sloww logs a message at slow level.
	Sloww(string, map[string]interface{})
	// WithDuration returns a new logger with the given duration.
	WithDuration(d time.Duration) Logger
	// WithFields returns a new logger with the given fields.
	WithFields(fields map[string]interface{}) Logger
	// SetLevel 设置日志级别
	SetLevel(level string) error

	V(v int) bool
}
