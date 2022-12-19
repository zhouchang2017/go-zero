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
	Formatter      string
	TimeFormatter  string `mapstructure:"time_formatter"`
	Level          string
	EnableFileLine bool `mapstructure:"enable_file_line"`
	Path           string
	//MaxAge 日志保存时间 默认15天 单位天
	MaxAge int

	//单位为m 默认200m 分割
	MaxFileSize int
	//文件保存的个数 默认200
	MaxBackup int
	Default   bool
}

func (c *LogConf) fillDefaultValue() {
	if c.MaxAge <= 0 {
		c.MaxAge = 15
	}
	if c.MaxBackup <= 0 {
		c.MaxBackup = 200
	}
	if c.MaxFileSize <= 0 {
		c.MaxFileSize = 200
	}
}

func (c *LogConf) buildWriter() (writers []io.Writer, err error) {
	if c.Path == "" {
		return nil, errors.New("log path is empty")
	}
	writers = make([]io.Writer, 0, 2)
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
