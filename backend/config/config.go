package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Name string `yaml:"name"`
	} `yaml:"db"`
}

func New() *Config {
	config := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	data, err := os.ReadFile(*config)
	if err != nil {
		log.Fatalf("error on read config file: %s", err)
	}

	var confifStruct Config
	err = yaml.Unmarshal(data, &confifStruct)
	if err != nil {
		log.Fatalf("error on parse config file: %s", err)
	}

	return &confifStruct
}
