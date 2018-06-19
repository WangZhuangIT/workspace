package flagutil

import (
	"flag"
	"go_libs/service/logger"
)

var signal = flag.String("s", "", "start or stop")
var confpath = flag.String("c", "", "config path")
var foreground = flag.Bool("f", false, "foreground")

func init() {
	flag.Parse()
}
func GetSignal() *string {
	return signal
}

func GetConfig() *string {
	logger.I("CONF", *confpath)
	return confpath
}

func GetForeground() *bool {
	return foreground
}
