package service

import (
	"log"

	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/internal/devserver"
)

const (
	// DevMode means development mode.
	DevMode = "dev"
	// TestMode means test mode.
	TestMode = "test"
	// RtMode means regression test mode.
	RtMode = "rt"
	// PreMode means pre-release mode.
	PreMode = "pre"
	// ProMode means production mode.
	ProMode = "pro"
)

// A ServiceConf is a service config.
type ServiceConf struct {
	Name       string
	Mode       string `default:"pro"`
	MetricsUrl string
	// Deprecated: please use DevServer
	Prometheus prometheus.Config
	Telemetry  trace.Config
	DevServer  devserver.Config
}

// MustSetUp sets up the service, exits on error.
func (sc ServiceConf) MustSetUp() {
	if err := sc.SetUp(); err != nil {
		log.Fatal(err)
	}
}

// SetUp sets up the service.
func (sc ServiceConf) SetUp() error {
	sc.initMode()
	prometheus.StartAgent(sc.Prometheus)

	if len(sc.Telemetry.Name) == 0 {
		sc.Telemetry.Name = sc.Name
	}
	trace.SetTraceName(sc.Name)
	trace.StartAgent(sc.Telemetry)
	proc.AddShutdownListener(func() {
		trace.StopAgent()
	})

	if len(sc.MetricsUrl) > 0 {
		stat.SetReportWriter(stat.NewRemoteWriter(sc.MetricsUrl))
	}
	devserver.StartAgent(sc.DevServer)

	return nil
}

func (sc ServiceConf) initMode() {
	switch sc.Mode {
	case DevMode, TestMode, RtMode, PreMode:
		load.Disable()
		stat.SetReporter(nil)
	}
}
