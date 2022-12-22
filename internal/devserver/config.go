package devserver

// Config is config for inner http server.
type Config struct {
	Enabled       bool
	Host          string
	Port          int    `default:"6470"`
	MetricsPath   string `default:"/metrics"`
	HealthPath    string `default:"/healthz"`
	EnableMetrics bool
	EnablePprof   bool
}
