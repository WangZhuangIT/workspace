package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-xorm/xorm"
)

type Logdata struct {
	Id                   int32     `xorm:"pk autoincr"`
	Timestamp            time.Time `xorm:"created"`
	Hostname             string    `json:"hostname" xorm:"varchar(25) notnull 'hostname'"`
	ServerName           string    `json:"server_name"`
	HttpXForwardedFor    string    `json:"http_x_forwarded_for"`
	XesApp               string    `json:"xes-app"`
	RemoteAddr           string    `json:"remote_addr"`
	RemoteUser           string    `json:"remote_user"`
	BodyBytesSent        int       `json:"body_bytes_sent"`
	RequestTime          float64   `json:"request_time"`
	UpstreamResponseTime string    `json:"upstream_response_time"`
	Status               int       `json:"status"`
	UpstreamStatus       string    `json:"upstream_status"`
	ConnectionRequests   int       `json:"connection_requests"`
	Request              string    `json:"request"`
	RequestMethod        string    `json:"request_method"`
	RequestBody          string    `json:"request_body"`
	HttpReferrer         string    `json:"http_referrer"`
	HttpCookie           string    `json:"http_cookie"`
	HttpXRequestId       string    `json:"http_x_request_id"`
	HttpXAppId           string    `json:"http_x_app_id"`
	HttpUserAgent        string    `json:"http_user_agent"`
}

func main() {

	data := `{"hostname": "gw-88-51", "server_name": "hwapi.xesv5.com", "http_x_forwarded_for": "-", "xes-app": "xes-app/hwapi-136-38", "remote_addr": "10.10.133.91", "remote_user": "-", "body_bytes_sent": 700, "request_time": 0.004, "upstream_response_time": "0.004", "status": 200, "upstream_status": "200", "connection_requests": 1, "request": "POST /StuHomework/getHomeworkInfo HTTP/1.1", "request_method": "POST", "request_body": "stuCouPlan%5B0%5D%5BstuCouId%5D=7213262&stuCouPlan%5B0%5D%5BplanId%5D=43675&stuCouPlan%5B0%5D%5Bcourse_id", "http_referrer": "-", "http_cookie": "-", "http_x_request_id": "f0830e283b5e0160b775e76ef060448e_0", "http_x_app_id": "100", "http_user_agent": "-"}`
	logdata := &Logdata{}
	jerr := json.Unmarshal([]byte(data), &logdata)

	if jerr != nil {
		log.Fatalln(jerr)
		return
	}

	var engine *xorm.Engine
	var err error
	engine, err = xorm.NewEngine("mysql", "root:123456@tcp(10.99.1.151:3306)/test?charset=utf8")
	if err = engine.Sync2(new(Logdata)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
	bools, error := engine.IsTableExist("log_data")
	if error != nil {
		log.Fatalln(error)
	}

	if bools {
		log.Println("cunzai")
	} else {
		log.Println("bucunzai")
	}

	if err != nil {
		log.Fatalln(err)
	}
	err = engine.Ping()

	if err != nil {
		log.Fatalln(err)
	}

	//如果需要设置连接池的空闲数大小
	engine.SetMaxIdleConns(10)
	//如果需要设置最大打开连接数
	engine.SetMaxOpenConns(10)

	engine.ShowSQL(true)
	// engine.SetTableMapper(core.SnakeMapper{})
	// engine.SetColumnMapper(core.SameMapper{})

	datas := []Logdata{}
	fmt.Println(cap(datas))
	fmt.Println(len(datas))
	//目前需要自行分割成每150条插入一次
	datas = append(datas, *logdata)
	id, errs := engine.Insert(datas)

	fmt.Println(cap(datas))

	if errs != nil {
		log.Fatalln(err)
	}

	fmt.Println(id)

	numbers := make(map[string]int)
	numbers["one"] = 1  //赋值
	numbers["ten"] = 10 //赋值
	numbers["three"] = 3

	//查询

	var querydata Logdata

	engine.Id(1).Cols("Timestamp").Get(&querydata)

	if error != nil {
		log.Fatalln(error)
	}

	timeNow := querydata.Timestamp.Format("2006-01-02 15:04:05")

	fmt.Println(timeNow)

	fmt.Println(querydata.Timestamp)

	querydata.ServerName = "wangdazhuang787887878"

	engine.AllCols().Id(5).Update(&querydata)

	finddata := []Logdata{}
	engine.Where("id>?", 10).Desc("id").Find(&finddata)

	for index, data := range finddata {
		fmt.Printf("%v, %v \n", index, data)
	}

	fmt.Println(len(finddata))
	fmt.Println(finddata[0].Id)

	total, _ := engine.Count(new(Logdata))

	fmt.Println(total)

	sum, _ := engine.Sum(new(Logdata), "id")

	fmt.Println(sum)

	_, errss := engine.Update(&Logdata{ServerName: "xinde", HttpXForwardedFor: "123"})

	if errss != nil {
		log.Fatalln(errss)
	}

	engine.Id(4).Delete(Logdata{})

	sql := "select * from logdata"
	results, err111 := engine.Query(sql)

	if err111 != nil {
		log.Fatalln(err111)
	}
	fmt.Println(results)

	sql = "update `logdata` set server_name=? where id=?"
	res, er := engine.Exec(sql, "xiaolun", 1)
	if er != nil {
		log.Fatalln(er)
	}

	fmt.Println(res)

	doSession(engine)

}

func doSession(engine *xorm.Engine) {
	session := engine.NewSession()
	defer session.Close()
	// add Begin() before any action
	err := session.Begin()
	user1 := Logdata{ServerName: "ddddd"}
	_, err = session.Insert(&user1)
	if err != nil {
		session.Rollback()
		return
	}
	user2 := Logdata{ServerName: "yyy"}
	_, err = session.Where("id = ?", 8888).Update(&user2)
	if err != nil {
		session.Rollback()
		return
	}

	_, err = session.Exec("delete from logdata where id = ?", 8888)
	if err != nil {
		session.Rollback()
		return
	}

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		return
	}
}
