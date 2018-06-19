package main

import (
	"go_libs/service/logger"

	"github.com/alecthomas/log4go"
)

func main() {
	log4go.LoadConfiguration("log.xml")
	log4go.Debug("123", "wangdazhuang")
	log4go.Debug("121312123", "wangdazhuang")
	// log4go.Log(log4go.INFO, "position", "tag++errMessage")
	logger.I("QUIT", "terminated")
}
