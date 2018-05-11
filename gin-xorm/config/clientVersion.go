package config

import (
	"gopkg.in/yaml.v2"
)

var cleintVersion map[string]ClientVersionItem

type ClientVersionItem struct {
	Version     string `yaml:"version"`
	ForceUpdate bool   `yaml:"forceUpdate"`
}

func readClientVersion(b []byte) error {
	err := yaml.Unmarshal(b, &cleintVersion)
	return err
}

func GetClientVersion() map[string]ClientVersionItem {
	data := make(map[string]ClientVersionItem)
	for k, v := range cleintVersion {
		data[k] = v
	}
	return data
}
