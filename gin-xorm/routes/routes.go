package routes

import (
	"fmt"
	. "gin-xorm/api"

	. "gin-xorm/config"

	configManager "git.xesv5.com/senior/lib/go/configManager"
	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {

		ReadBasicConfig(configManager.GetConfigMap()["basic.yml"])
		basicConfig := GetBasicConfig()
		fmt.Println(basicConfig.Appid)
		fmt.Println(basicConfig.GinMode)
		fmt.Println(basicConfig.HttpPort)
		fmt.Println(basicConfig.LogLevel)
	})
	//POST 		增加
	//GET 		获取
	//PUT 		更新
	//DELETE 	删除
	e.POST("/student/", AddStuAPI)

	e.GET("/student/", GetStuAPI)
	e.GET("/student/:id", AddStuAPI)

	e.PUT("/student/:id", AddStuAPI)

	e.DELETE("/student/", AddStuAPI)
	e.DELETE("/student/:id", AddStuAPI)
}
