package routes

import (
	"fmt"
	. "gin-xorm/api"
	"gin-xorm/middleware"
	"net/http"

	configManager "git.xesv5.com/senior/lib/go/configManager"

	. "gin-xorm/config"

	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.Engine) {
	e.LoadHTMLGlob("templates/*")

	e.GET("/", func(c *gin.Context) {

		ReadBasicConfig(configManager.GetConfigMap()["basic.yml"])
		basicConfig := GetBasicConfig()
		fmt.Println(basicConfig.Appid)
		fmt.Println(basicConfig.GinMode)
		fmt.Println(basicConfig.HttpPort)
		fmt.Println(basicConfig.LogLevel)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "this is a msg",
		})
	})
	//POST 		增加
	//GET 		获取
	//PUT 		更新
	//DELETE 	删除
	e.POST("/student/", AddStuAPI)

	e.GET("/student/", GetStuSliceAPI)
	e.GET("/student/:id", GetStuAPI)

	e.PUT("/student/", UpStuAPI)

	e.DELETE("/student/", middleware.Auth(), DelStuAPI)

	e.GET("/mongo/", MongoAPI)
}
