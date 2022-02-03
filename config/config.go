package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Addr string `yaml:"address"`
	Crt  string `yaml:"certificate"`
	Key  string `yaml:"key"`
	Log  string `yaml:"log"`
}

func Make() config {
	content, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatal("[FATAL] cannot read config file", err)
	}

	var cfg config
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatal("[FATAL] cannot read file contents", err)
	}

	return cfg
}
