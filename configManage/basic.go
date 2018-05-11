package main

import (
	"gopkg.in/yaml.v2"
)

type BasicConfig struct {
	LogLevel int    `yaml:"logLevel"`
	HttpPort int    `yaml:"httpPort"`
	Appid    int    `yaml:"appid"`
	GinMode  string `yaml:"ginMode"`
}

var basicConfig BasicConfig

func readBasicConfig(b []byte) error {
	err := yaml.Unmarshal(b, &basicConfig)
	return err
}

func GetBasicConfig() BasicConfig {
	return basicConfig
}
