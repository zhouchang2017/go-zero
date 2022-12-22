package conf

import (
	"bytes"
	"github.com/creasty/defaults"
	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
)

var validate = validator.New()

// Load loads config into v from file, .json, .yaml and .yml are acceptable.
func Load(file string, v interface{}, opts ...Option) (err error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	ext := strings.TrimLeft(strings.ToLower(path.Ext(file)), ".")

	var opt options
	for _, o := range opts {
		o(&opt)
	}

	if opt.env {
		content = []byte(os.ExpandEnv(string(content)))
	}
	viper.AutomaticEnv()
	viper.SetConfigType(ext)
	viper.SetConfigFile(file)

	if err := viper.ReadConfig(bytes.NewReader(content)); err != nil {
		return err
	}
	err = unmarshal(v)
	if err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("%s config change...\n", file)
		_ = unmarshal(v)
	})
	return nil
}

func unmarshal(v interface{}) error {
	err := viper.Unmarshal(v)
	if err != nil {
		return err
	}
	_ = defaults.Set(v)
	return validate.Struct(v)
}

// LoadConfig loads config into v from file, .json, .yaml and .yml are acceptable.
// Deprecated: use Load instead.
func LoadConfig(file string, v interface{}, opts ...Option) error {
	return Load(file, v, opts...)
}

// LoadFromJsonBytes loads config into v from content json bytes.
func LoadFromJsonBytes(content []byte, v interface{}) error {
	viper.SetConfigType("json")
	err := viper.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	return unmarshal(v)
}

// LoadConfigFromJsonBytes loads config into v from content json bytes.
// Deprecated: use LoadFromJsonBytes instead.
func LoadConfigFromJsonBytes(content []byte, v interface{}) error {
	return LoadFromJsonBytes(content, v)
}

// LoadFromTomlBytes loads config into v from content toml bytes.
func LoadFromTomlBytes(content []byte, v interface{}) error {
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	return unmarshal(v)
}

// LoadFromYamlBytes loads config into v from content yaml bytes.
func LoadFromYamlBytes(content []byte, v interface{}) error {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	return unmarshal(v)
}

// LoadConfigFromYamlBytes loads config into v from content yaml bytes.
// Deprecated: use LoadFromYamlBytes instead.
func LoadConfigFromYamlBytes(content []byte, v interface{}) error {
	return LoadFromYamlBytes(content, v)
}

// MustLoad loads config into v from path, exits on error.
func MustLoad(path string, v interface{}, opts ...Option) {
	if err := Load(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}
