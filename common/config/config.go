package config

import (
	"io/ioutil"
	"log"

	"github.com/yangb8/webservice/common/logger"
	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Log struct {
		Location logger.LogLocation `yaml:"location"`
	} `yaml:"log"`
	Sentry struct {
		Enabled bool   `yaml: "enabled"`
		Dsn     string `yaml:"dsn"`
	} `yaml:"sentry"`
	Statsd struct {
		Enabled bool   `yaml: "enabled"`
		Address string `yaml:"address"`
	} `yaml:"statsd"`
	Storage struct {
		Kind  string `yaml:"kind"`
		Local struct {
			Path string `yaml:"path"`
		}
		S3 struct {
			URL    string `yaml:"url"`
			ID     string `yaml:"id"`
			Key    string `yaml:"key"`
			Bucket string `yaml: "bucket"`
		} `yaml:"s3"`
	} `yaml:"storage"`
}

// GetConfig ...
func GetConfig() *Config {
	yf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("failed to open config.yaml")
	}
	res := &Config{}
	if err = yaml.Unmarshal(yf, res); err != nil {
		log.Fatal("failed to parse config.yaml")
	}
	return res
}
