package dbdao

import (
	"log"
	"strings"
	"sync/atomic"

	"go_libs/utils/confutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBDao struct {
	Engine *xorm.Engine
	quiter chan struct{}
}

var (
	db_instance map[string][]*DBDao
	curDbPoses  map[string]*uint64 // 当前选择的数据库
)

func newDBDaoWithParams(host string, driver string, max, maxIdle int) (Db *DBDao) {
	Db = new(DBDao)
	engine, err := xorm.NewEngine(driver, host)
	Db.Engine = engine
	//TODO: 增加存活检查
	if err != nil {
		log.Fatal(err)
	}
	/*Db.Engine.Logger.SetLevel(core.LOG_DEBUG)
	  Db.Engine.ShowSQL = true
	  Db.Engine.ShowInfo = true
	  Db.Engine.ShowDebug = true
	  Db.Engine.ShowErr = true
	  Db.Engine.ShowWarn = true*/
	Db.Engine.SetMaxOpenConns(max)
	Db.Engine.SetMaxIdleConns(maxIdle)
	return
}

func GetDefault(cluster string) *DBDao {
	return GetDbInstance("default", cluster)
}

func init() {
	db_instance = make(map[string][]*DBDao, 0)
	curDbPoses = make(map[string]*uint64)
	//idc := confdao.GetIDC()
	idc := ""
	for cluster, hosts := range confutil.GetConfArrayMap("MysqlCluster") {
		items := strings.Split(cluster, ".")
		//必须包含 writer 和 reader
		if len(items) < 2 {
			continue
		}
		//过滤IDC
		if len(items) > 2 && items[2] != idc {
			continue
		}
		instance := items[0] + "." + items[1]
		dbs := make([]*DBDao, 0)
		for _, host := range hosts {
			dbs = append(dbs, newDBDaoWithParams(host, "mysql", 100, 30))
		}
		db_instance[instance] = dbs
		curDbPoses[instance] = new(uint64)
	}
}

func GetDbInstance(db, cluster string) *DBDao {
	key := db + "." + cluster
	if instances, ok := db_instance[key]; ok {
		// round-robin选择数据库
		cur := atomic.AddUint64(curDbPoses[key], 1) % uint64(len(instances))
		return instances[cur]
	} else {
		return nil
	}
}

func (this *DBDao) Close() {
	this.Engine.Close()
}
