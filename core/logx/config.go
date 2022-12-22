package logx

import (
	"errors"
	rotatelogs "github.com/natefinch/lumberjack"
	rotatelogsbytime "github.com/patch-mirrors/file-rotatelogs"
	"io"
	"os"
	"strings"
	"time"
)

// A LogConf is a logging config.
type LogConf struct {
	// 是否控制台同时输出
	Stdout bool
	// 格式化输出，json|plain
	Formatter string `default:"plain"`
	// 时间格式化
	// example for Time.Format demonstrates
	// Formatter=json  TimeFormatter=2006-01-02T15:04:05Z07:00
	// Formatter=plain TimeFormatter=2006-01-02 15:04:05.000000
	TimeFormatter string
	// 日志级别，debug|info|warn|error|panic|fatal
	Level string `default:"debug"`
	// 是否显示行号
	EnableFileLine bool
	Path           string
	//MaxAge 日志保存时间 默认15天 单位天
	MaxAge int `default:"15"`
	//单位为m 默认200m 分割
	MaxFileSize int `default:"200"`
	//文件保存的个数 默认200
	MaxBackup int `default:"200"`
	Default   bool
}

func (c *LogConf) buildWriter() (writers []io.WriteCloser, err error) {
	if c.Path == "" {
		return nil, errors.New("log path is empty")
	}
	writers = make([]io.WriteCloser, 0, 2)
	if strings.Index(c.Path, "%Y%m%d") >= 0 {
		writer, err := rotatelogsbytime.New(
			c.Path,
			rotatelogsbytime.WithRotationCount(uint(c.MaxBackup)),
			rotatelogsbytime.WithRotationTime(time.Hour),
		)
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	} else {
		writers = append(writers, &rotatelogs.Logger{
			Filename:   c.Path,
			MaxSize:    c.MaxFileSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackup,
			LocalTime:  true,
		})
	}
	if c.Stdout {
		writers = append(writers, os.Stdout)
	}
	return writers, nil
}

func (c LogConf) New() (Logger, error) {
	logger, err := NewWithConf(&c)
	if err != nil {
		return nil, err
	}
	if c.Default {
		SetGlobalLogger(logger)
	}
	return logger, nil
}
