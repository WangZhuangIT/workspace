package output

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/kshvakov/clickhouse"
)

type DB struct {
	dbConn *sql.DB
}

func NewDBConn(host string, port string, database string) *DB {
	var err error
	db := new(DB)
	dsn := fmt.Sprintf("tcp://%s:%s?database=%s", host, port, database)
	db.dbConn, err = sql.Open("clickhouse", dsn)
	if err != nil {
		panic(err)
	}
	err = db.dbConn.Ping()
	if err != nil {
		panic(err)
	}
	db.dbConn.SetMaxOpenConns(30)

	return db
}

func (cDb *DB) ServicePv(param []map[string]interface{}) error {
	sqlStr := `INSERT INTO xes_service_pv (date, time, logid, hostname, prelogid, pageid, sessid, appid, token, ver, os, ua, xesid, userid, devid, url, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if l := len(param); l <= 0 {
		return errors.New("empty data")
	}

	tx, err := cDb.dbConn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}

	for _, row := range param {
		_, err = stmt.Exec(row["date"], row["time"], row["logid"], row["hostname"], row["prelogid"], row["pageid"], row["sessid"], row["appid"], row["token"], row["ver"], row["os"], row["ua"], row["xesid"], row["userid"], row["devid"], row["url"], row["data"])
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (cDb *DB) Insert(table string, param map[string]interface{}) error {
	l := len(param)
	if l <= 0 {
		return errors.New("empty data")
	}

	var columns, holder string
	var values []interface{}
	i := 1
	sql := "INSERT INTO " + table
	for col, val := range param {
		columns += col
		holder += "?"
		values = append(values, val)
		if i < l {
			columns += ", "
			holder += ", "
		}
		i++
	}

	sql += " (" + columns + ") VALUES (" + holder + ")"
	// fmt.Println(sql, values)
	tx, err := cDb.dbConn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
