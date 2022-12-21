package logx

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

var _ Logger = &ZapLogger{}

const namedSep = "$#&)"
const splitNamedSep = "." + namedSep

// ZapLogger zap logger
type ZapLogger struct {
	base  *zap.Logger
	level *zap.AtomicLevel
}

func (z *ZapLogger) clone() *ZapLogger {
	copy := *z
	return &copy
}

// V checks if meet required log level.
func (z *ZapLogger) V(v int) bool {
	return v >= 2
}

func (z *ZapLogger) Slow(i ...interface{}) {
	z.log(zap.ErrorLevel, "", i)
}

func (z *ZapLogger) Slowf(s string, i ...interface{}) {
	z.log(zap.ErrorLevel, s, i)
}

func (z *ZapLogger) Slowv(i interface{}) {
	z.log(zap.ErrorLevel, "", []interface{}{i})
}

func (z *ZapLogger) Sloww(s string, fields map[string]interface{}) {
	var zFields []zap.Field
	for k, v := range fields {
		zFields = append(zFields, zap.Any(k, v))
	}
	z.log(zap.ErrorLevel, "", []interface{}{zFields})
}

func (z *ZapLogger) WithDuration(d time.Duration) Logger {
	return &ZapLogger{
		base:  z.base.With(zap.Duration("duration", d)),
		level: z.level,
	}
}

func (z *ZapLogger) buildRequestId(id string) string {
	buf := _pool.Get()
	defer buf.Free()
	buf.AppendByte('.')
	buf.AppendString(namedSep)
	buf.AppendString(id)
	return buf.String()
}

func (z *ZapLogger) WithRequestId(id string) Logger {
	return &ZapLogger{
		base:  z.base.Named(z.buildRequestId(id)),
		level: z.level,
	}
}

func (z *ZapLogger) WithField(key string, value interface{}) Logger {
	logger := &ZapLogger{
		level: z.level,
	}
	logger.base = z.base.With(zap.Any(key, value))
	return logger
}

func (z *ZapLogger) WithFields(fields map[string]interface{}) Logger {
	i := _fPool.Get().([]zap.Field)
	defer _fPool.Put(i)
	logger := &ZapLogger{
		level: z.level,
	}

	for key, v := range fields {
		i = append(i, zap.Any(key, v))
	}
	logger.base = z.base.With(i...)
	return logger
}

func (z *ZapLogger) WithError(err error) Logger {
	return &ZapLogger{
		base:  z.base.With(zap.Error(err)),
		level: z.level,
	}
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	z.log(zap.DebugLevel, format, args)
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
	z.log(zap.InfoLevel, format, args)
}

func (z *ZapLogger) Printf(format string, args ...interface{}) {
	z.log(zap.InfoLevel, format, args)
}

func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	z.log(zap.WarnLevel, format, args)
}

func (z *ZapLogger) Warningf(format string, args ...interface{}) {
	z.log(zap.WarnLevel, format, args)
}

func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	z.log(zap.ErrorLevel, format, args)
}

func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	z.log(zap.FatalLevel, format, args)
}

func (z *ZapLogger) Panicf(format string, args ...interface{}) {
	z.log(zap.PanicLevel, format, args)
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.log(zap.DebugLevel, "", args)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.log(zap.InfoLevel, "", args)
}

func (z *ZapLogger) Print(args ...interface{}) {
	z.log(zap.InfoLevel, "", args)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.log(zap.WarnLevel, "", args)
}

func (z *ZapLogger) Warning(args ...interface{}) {
	z.log(zap.WarnLevel, "", args)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.log(zap.ErrorLevel, "", args)
}

func (z *ZapLogger) Fatal(args ...interface{}) {
	z.log(zap.FatalLevel, "", args)
}

func (z *ZapLogger) Panic(args ...interface{}) {
	z.log(zap.PanicLevel, "", args)
}

func (z *ZapLogger) Debugln(args ...interface{}) {
	z.log(zap.DebugLevel, "", args)
}

func (z *ZapLogger) Infoln(args ...interface{}) {
	z.log(zap.InfoLevel, "", args)
}

func (z *ZapLogger) Println(args ...interface{}) {
	z.log(zap.InfoLevel, "", args)
}

func (z *ZapLogger) Warnln(args ...interface{}) {
	z.log(zap.WarnLevel, "", args)
}

func (z *ZapLogger) Warningln(args ...interface{}) {
	z.log(zap.WarnLevel, "", args)
}

func (z *ZapLogger) Errorln(args ...interface{}) {
	z.log(zap.ErrorLevel, "", args)
}

func (z *ZapLogger) Fatalln(args ...interface{}) {
	z.log(zap.FatalLevel, "", args)
}

func (z *ZapLogger) Panicln(args ...interface{}) {
	z.log(zap.PanicLevel, "", args)
}

func (z *ZapLogger) Tracef(format string, args ...interface{}) {
	z.log(zap.DebugLevel, format, args)
}

func (z *ZapLogger) Trace(args ...interface{}) {
	z.log(zap.DebugLevel, "", args)
}

func (z *ZapLogger) Traceln(args ...interface{}) {
	z.log(zap.DebugLevel, "", args)
}

func (z *ZapLogger) SetLevel(level string) error {
	z.level.SetLevel(strToLevel(level))
	return nil
}

// log message with Sprint, Sprintf, or neither.
func (z *ZapLogger) log(lvl zapcore.Level, template string, fmtArgs []interface{}) {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if lvl < zap.DPanicLevel && !z.base.Core().Enabled(lvl) {
		return
	}

	msg, fields := getMessage(template, fmtArgs)
	if ce := z.base.Check(lvl, msg); ce != nil {
		ce.Write(fields...)
	}
}

var _fPool = sync.Pool{New: func() interface{} { return make([]zap.Field, 0, 10) }}

// getMessage format with Sprint, Sprintf, or neither.
func getMessage(template string, fmtArgs []interface{}) (string, []zap.Field) {
	if len(fmtArgs) == 0 {
		return template, nil
	}

	if template == "" && len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str, nil
		}
	}

	fields := _fPool.Get().([]zap.Field)
	defer _fPool.Put(fields)
	args := make([]interface{}, 0, len(fmtArgs))
	for index, arg := range fmtArgs {
		switch arg.(type) {
		case string:
			if index == 0 && template == "" {
				template = arg.(string)
			} else {
				args = append(args, arg)
			}
		case zap.Field:
			fields = append(fields, arg.(zap.Field))
		default:
			if template == "" {
				template = "%v"
			}
			args = append(args, arg)
		}
	}

	if len(args) == 0 {
		return template, fields
	}

	return fmt.Sprintf(template, args...), fields
}
