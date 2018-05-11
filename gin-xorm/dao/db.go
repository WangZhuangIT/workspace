package dao

import (
	"fmt"
	. "gin-xorm/config"
	"log"

	configManager "git.xesv5.com/senior/lib/go/configManager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var Db *xorm.Engine

func InitDb() {
	configMap := configManager.GetConfigMap()
	ReadDbConfig(configMap["db.yml"])
	master := GetDbConfig().Master
	var err error
	Db, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v", master.User, master.Pwd, master.Host, master.Port, master.Db, master.Charset))
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//如果需要设置连接池的空闲数大小
	Db.SetMaxIdleConns(10)
	//如果需要设置最大打开连接数
	Db.SetMaxOpenConns(10)
	Db.ShowSQL(true)
}
