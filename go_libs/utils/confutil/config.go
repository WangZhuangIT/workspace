package confutil

import (
	"go_libs/service/logger"
	"go_libs/utils/flagutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/kardianos/osext"
)

var (
	g_cfg *goconfig.ConfigFile
)
var GConf *goconfig.ConfigFile
var config_cache = make(map[string]*goconfig.ConfigFile, 0)
var USER_CONF_PATH string

func InitConfig() {
	if g_cfg != nil {
		return
	}
	config_path := flagutil.GetConfig()
	var err error
	if len(*config_path) == 0 {
		*config_path = "../conf/conf.ini"
	}
	logger.I("CONF INIT", *config_path)
	if g_cfg, err = Load(*config_path); err != nil {
		logger.I("Conf", err)
	} else {
		GConf = g_cfg
	}
}
func home() string {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return "/home/work"
}
func Binhome() string {
	if len(USER_CONF_PATH) > 0 {
		return USER_CONF_PATH
	}
	if path, err := osext.ExecutableFolder(); err == nil {
		if strings.HasPrefix(path, "/tmp/go-build") {
			return home() + "/conf/"
		}
		return path
	} else {
		return "."
	}
}

func Load(path string) (cfg *goconfig.ConfigFile, err error) {
	if path[len(path)-4:] != ".ini" {
		path = path + ".ini"
	}
	var ok bool
	if cfg, ok = config_cache[path]; !ok {
		path = Binhome() + "/" + path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = home() + "/conf/" + filepath.Base(path)
		}
		if cfg, err = goconfig.LoadConfigFile(path); err == nil {
			config_cache[path] = cfg
		}
	}
	return
}

func GetConf(sec, key string) string {
	InitConfig()
	if g_cfg == nil {
		logger.W("Conf", "NOT_FOUND[sec:%s,key:%s] AT:%v", sec, key, flagutil.GetConfig())
		return ""
	}
	return g_cfg.MustValue(sec, key, "")
}

func GetConfDefault(sec, key, def string) string {
	InitConfig()
	if g_cfg == nil {
		logger.W("Conf", "NOT_FOUND[sec:%s,key:%s] AT:%v", sec, key, flagutil.GetConfig())
		return ""
	}
	return g_cfg.MustValue(sec, key, def)
}

func GetConfs(sec, key string) []string {
	InitConfig()
	if g_cfg == nil {
		logger.W("Conf", "NOT_FOUND[sec:%s,key:%s]", sec, key)
		return []string{}
	}
	return g_cfg.MustValueArray(sec, key, " ")
}

func GetConfStringMap(sec string) (ret map[string]string) {
	InitConfig()
	if g_cfg == nil {
		logger.W("Conf", "NOT_FOUND[sec:%s]", sec)
		return nil
	}
	var err error
	if ret, err = g_cfg.GetSection(sec); err != nil {
		logger.W("Conf", err)
		ret = make(map[string]string, 0)
	}
	return
}

func GetConfArrayMap(sec string) (ret map[string][]string) {
	InitConfig()
	if g_cfg == nil {
		logger.W("Conf", "NOT_FOUND[sec:%s]", sec)
		return nil
	}
	ret = make(map[string][]string, 0)
	confs := g_cfg.GetKeyList(sec)
	for _, k := range confs {
		ret[k] = g_cfg.MustValueArray(sec, k, " ")
	}
	return
}
