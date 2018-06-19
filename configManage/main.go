package main

import (
	"fmt"

	config "git.xesv5.com/senior/lib/go/configManager"
	yaml "gopkg.in/yaml.v2"
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

func main() {
	//  1.ETCD  2.File
	flag := 2
	initConfig(flag)
	loadFile()
}

func initConfig(flag int) {
	if flag == 1 {
		dir := "config/"
		endPoints := []string{"10.99.1.151:2379"}
		config.Init(dir, config.NewEtcdLoader(endPoints, dir))
	} else {
		config.Init("./dev", nil)
	}

}

func loadFile() {
	readBasicConfig(config.GetConfigMap()["basic.yml"])
	fmt.Println(basicConfig)

	value := string(config.GetConfigByKey("basic.yml"))
	fmt.Println(value)
}
