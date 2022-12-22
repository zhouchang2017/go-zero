package logx

import (
	"go.uber.org/zap"
	"os"
	"sync"
)

var (
	globalLogger Logger = NewDevelopment()
	lock         sync.Mutex

	_mustLogger = newMustLogger()
)

func SetGlobalLogger(l Logger) {
	lock.Lock()
	defer lock.Unlock()
	if l == nil {
		return
	}
	globalLogger = l
}

func GlobalLogger() Logger {
	return globalLogger
}

// CloneWithAddCallerSkip 返回一个新的Logger实例
func CloneWithAddCallerSkip(skip int) Logger {
	if logger, ok := globalLogger.(*ZapLogger); ok {
		return &ZapLogger{
			base:  logger.clone().base.WithOptions(zap.AddCallerSkip(skip)),
			level: logger.level,
		}
	}
	return globalLogger
}

func Disable() {
	// nothing
}

func Must(err error) {
	if err == nil {
		return
	}

	_mustLogger.Error(err.Error())
	os.Exit(1)
}
