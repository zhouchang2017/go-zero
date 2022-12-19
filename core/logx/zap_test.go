package logx

import (
	"errors"
	"testing"
)

func TestZapLogger_WithRequestId(t *testing.T) {
	conf := make(LogConfigMap)
	conf["access"] = LogConf{
		Stdout:    false,
		Formatter: "json",
		Level:     "info",
		Path:      "/tmp/access.log",
	}

	conf["app"] = LogConf{
		Stdout:         true,
		Formatter:      "plain",
		EnableFileLine: true,
		Path:           "/tmp/app.log",
		Default:        true,
	}
	MustSetup(conf)
	log2, err := Get("access")
	if err != nil {
		panic(err)
	}
	log2.Info("hello")
	appLogger, err := Get("app")
	if err != nil {
		panic(err)
	}
	appLogger.Info("no requestId")
	l := appLogger.WithRequestId("xxxxx1111.2222")
	l.Infof("hello")
	l.WithRequestId("xxx2222.2222").Info("bye")

	development := newJsonDevelopment()
	development.WithRequestId("requestxxxxx01").Info("ok")

	Must(errors.New("err"))
}
