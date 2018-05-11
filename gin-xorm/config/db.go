package config

import (
	"gopkg.in/yaml.v2"
)

type DbConfig struct {
	Master DbItem   `yaml:"master"`
	Slaves []DbItem `yaml:"slaves"`
}

type DbItem struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
	Pwd     string `yaml:"pwd"`
	Db      string `yaml:"db"`
	Charset string `yaml:"charset"`
}

var dbConfig DbConfig

func ReadDbConfig(b []byte) error {
	err := yaml.Unmarshal(b, &dbConfig)
	return err
}

func GetDbConfig() DbConfig {
	return dbConfig
}
