package confdao

import (
	"strings"

	"go_libs/utils/confutil"
	"go_libs/utils/flagutil"
	"github.com/Unknwon/goconfig"
	"github.com/spf13/cast"
)

var conf_ini_cache *goconfig.ConfigFile

func init() {
}

func get_conf_config() (conf *goconfig.ConfigFile, err error) {
	if conf_ini_cache != nil {
		return conf_ini_cache, nil
	}

	configPath := flagutil.GetConfig()
	if len(configPath) <= 0 {
		configPath = "../conf/conf.ini"
	}
	conf, err = confutil.Load(configPath)
	conf_ini_cache = conf
	return
}

