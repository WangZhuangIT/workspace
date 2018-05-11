package main

import (
	"fmt"
	"time"

	config "git.xesv5.com/senior/lib/go/configManager"
)

func main() {
	testUse()
}

func testUse() {
	// 1 ETCD
	// 2 File
	flag := 1
	initConfig(flag)
	go func() {
		time.AfterFunc(time.Duration(5*time.Second), func() {
			if flag == 1 {
				loadEtcd()
			} else {
				loadFile()
			}
		})

		time.AfterFunc(time.Duration(10*time.Second), func() {
			if flag == 1 {
				loadEtcd()
			} else {
				loadFile()
			}
		})

	}()
	//循环
	select {}
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

func loadEtcd() {
	for k, v := range config.GetConfigMap() {
		fmt.Printf("%v : %v", k, string(v))
	}
}

func loadFile() {

	readBasicConfig(config.GetConfigMap()["./dev/basic.yml"])
	basicConfig := GetBasicConfig()
	fmt.Println(basicConfig.Appid)
	fmt.Println(basicConfig.GinMode)
	fmt.Println(basicConfig.HttpPort)
	fmt.Println(basicConfig.LogLevel)
}
