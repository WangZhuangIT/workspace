package database

import (
	"database/sql"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

func init() {
	var err error
	SqlDB, err = sql.Open("mysql", "root:123456@tcp(10.99.2.56:3306)/test?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	SqlDB.SetMaxIdleConns(20) //数据库最大空闲连接
	SqlDB.SetMaxOpenConns(20) //数据库最大连接

	if err := SqlDB.Ping(); err != nil {
		log.Fatalln(err)
	}
}
