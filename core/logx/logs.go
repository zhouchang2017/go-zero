package logx

import (
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
)

var (
	globalLogger Logger = NewDevelopment()
	lock         sync.Mutex
	setupOnce    sync.Once

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

// MustSetup sets up logging with given config c. It exits on error.
func MustSetup(c LogConfigMap) {
	setupOnce.Do(func() {
		err := c.Init()
		if err == nil {
			return
		}

		msg := err.Error()
		log.Print(msg)
		globalLogger.Error(msg)
		os.Exit(1)
	})
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
