package config

import (
	"fmt"

	configManager "git.xesv5.com/senior/lib/go/configManager"
)

//初始化配置文件configManager
func InitConfigManager(dir string) {
	configManager.Init(dir, nil)
	fmt.Println("初始化配置文件")
}

func InitEtcdConfigManager() {
	dir := "config/"
	endPoints := []string{"10.99.1.151:2379"}
	configManager.Init(dir, configManager.NewEtcdLoader(endPoints, dir))
	fmt.Println("初始化ETCD中的配置")
}
