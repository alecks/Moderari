package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type configModel struct {
	Token     string `yaml:"token"`
	Prefix    string `yaml:"prefix"`
	SentryDSN string `yaml:"sentry-dsn"`
	DB        struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	} `yaml:"db"`
	Colors struct {
		Info  int `yaml:"info"`
		Error int `yaml:"error"`
	} `yaml:"colours"`
	
	Version string `yaml:"version"`
}

// C is the main configuration used by the program.
var C configModel

func init() {
	filename := "config.yml"
	if envFilename := os.Getenv("CONFIG_FILE"); envFilename != "" {
		filename = envFilename
	}

	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(configFile, &C); err != nil {
		panic(err)
	}
}
