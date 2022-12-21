package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string
	Port int    `default:"9101"`
	Path string `default:"/metrics"`
}
