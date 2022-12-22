package logx

import (
	"github.com/zeromicro/go-zero/core/iox"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func NewTestLogger(w io.Writer) *ZapLogger {
	return New("plain", "", "debug", true, 3, iox.NopCloser(w))
}

// NewDevelopment 开发Logger
func NewDevelopment() *ZapLogger {
	return New("plain", "", "debug", true, 0)
}

func newMustLogger() *ZapLogger {
	return New("plain", "", "debug", true, 3)
}

func newJsonDevelopment() *ZapLogger {
	return New("json", "", "debug", false, 0)
}

// NewWithConf 通过配置实例化 Logger
func NewWithConf(conf *LogConf) (logger *ZapLogger, err error) {
	if conf == nil {
		return NewDevelopment(), nil
	}
	writers, err := conf.buildWriter()
	if err != nil {
		return nil, err
	}
	return New(conf.Formatter, conf.TimeFormatter, conf.Level, conf.EnableFileLine, 0, writers...), nil
}

// New 构造函数
func New(format string, timeFormat string, level string, enableFileLine bool, callerSkip int,
	writers ...io.WriteCloser) *ZapLogger {
	if callerSkip == 0 {
		callerSkip = 2
	}
	logger, le := newZapLogger(format, timeFormat, level, enableFileLine, writers...)
	if callerSkip != 0 {
		logger = logger.WithOptions(zap.AddCallerSkip(callerSkip))
	}
	return &ZapLogger{
		base:  logger,
		level: le,
	}
}

func newZapLogger(format string, timeFormat string, level string, enableFileLine bool,
	writers ...io.WriteCloser) (*zap.Logger, *zap.AtomicLevel) {
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(strToLevel(level))

	if len(writers) == 0 {
		// use stdout
		writers = append(writers, os.Stdout)
	}

	ws := make([]zapcore.WriteSyncer, 0, len(writers))
	for _, w := range writers {
		ws = append(ws, zapcore.AddSync(w))
	}
	syncers := zap.CombineWriteSyncers(ws...)

	core := zapcore.NewCore(
		newEncoder(format, timeFormat),
		syncers,
		atomicLevel,
	)

	return zap.New(core, zap.WithCaller(enableFileLine)), &atomicLevel
}

func newEncoder(format string, timeFormat string) zapcore.Encoder {
	switch format {
	case "plain":
		return newPlainEncoder(timeFormat)
	case "json":
		return newJsonEncoder(timeFormat)
	default:
		log.Printf("format: %s not allow, use default json format", format)
		return newJsonEncoder(timeFormat)
	}
}

func newJsonEncoder(timeFormat string) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeCaller = encodeCaller
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	config.TimeKey = "time"
	config.EncodeTime = newEncodeTime(timeFormat)
	config.EncodeName = encodeName
	config.NameKey = requestIdKey
	return zapcore.NewJSONEncoder(config)
}

func newPlainEncoder(timeFormat string) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeCaller = encodeCaller
	if timeFormat == "" {
		timeFormat = "2006-01-02 15:04:05.000000"
	}

	config.ConsoleSeparator = " "
	config.EncodeTime = newEncodeTime(timeFormat)
	config.EncodeLevel = encodeLevel
	config.EncodeName = encodeName
	return NewPlainEncoder(config)
}

var _pool = buffer.NewPool()

func strToLevel(str string) zapcore.Level {
	switch str {
	case "debug", "trace":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.DebugLevel
	}
}

func newEncodeTime(timeFormat string) func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(timeFormat))
	}
}

func encodeLevel(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(level.CapitalString()[:4])
}

func encodeName(s string, encoder zapcore.PrimitiveArrayEncoder) {
	split := strings.Split(s, splitNamedSep)
	encoder.AppendString(split[len(split)-1])
}

func encodeCaller(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
	pool := _pool.Get()
	defer pool.Free()
	idx := strings.LastIndexByte(caller.File, '/')
	if idx == -1 {
		pool.AppendString(caller.FullPath())
		pool.AppendByte(':')
		pool.AppendInt(int64(caller.Line))
	} else {
		pool.AppendString(caller.File[idx+1:])
		pool.AppendByte(':')
		pool.AppendInt(int64(caller.Line))
	}
	encoder.AppendString(pool.String())
}
