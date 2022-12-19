package logx

import (
	"fmt"
	"sync"
)

// LogConfigMap 日志初始化配置集合
type LogConfigMap map[string]LogConf

var loggers sync.Map

// Get 通过name获取Logger
func Get(name string) (Logger, error) {
	res, ok := loggers.Load(name)
	if !ok {
		return nil, fmt.Errorf("%s logger not found", name)
	}
	return res.(Logger), nil
}

// Init viper初始化函数
func (config *LogConfigMap) Init() (err error) {
	for subSection, item := range *config {
		globalLogger.Info("log [%s] loading", subSection)
		logger, err := NewWithConf(&item)
		if err != nil {
			globalLogger.Errorf("log [%s] init err: %s", subSection, err.Error())
			return err
		}
		if item.Default {
			globalLogger.Info("log [%s] set as default log", subSection)
			SetGlobalLogger(logger)
		}
		loggers.Store(subSection, logger)
		globalLogger.Info("log [%s] load success", subSection)
	}
	return nil
}
