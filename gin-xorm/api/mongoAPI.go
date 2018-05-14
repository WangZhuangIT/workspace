package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Persons struct {
	Name  string
	Phone string
}

func MongoAPI(context *gin.Context) {

	fmt.Println(123)

	//可本地可远程，不指定协议时默认为http协议访问，此时需要设置 mongodb 的nohttpinterface=false来打开httpinterface。
	//也可以指定mongodb协议，如 "mongodb://127.0.0.1:27017"
	var MOGODB_URI = "mongodb://10.99.1.151:27017"
	//连接
	session, err := mgo.Dial(MOGODB_URI)
	//连接失败时终止
	if err != nil {
		panic(err)
	}
	//延迟关闭，释放资源
	defer session.Close()
	//设置模式
	session.SetMode(mgo.Monotonic, true)
	//选择数据库与集合
	c := session.DB("adatabase").C("acollection")
	//插入文档
	err = c.Insert(&Persons{Name: "Ale", Phone: "+55 53 8116 9639"},
		&Persons{Name: "Cla", Phone: "+55 53 8402 8510"})
	//出错判断
	if err != nil {
		log.Fatal(err)
	}
	//查询文档
	result := Persons{}
	//注意mongodb存储后的字段大小写问题
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	//出错判断
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone:", result.Phone)
}
