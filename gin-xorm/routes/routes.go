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
	e.PUT("/student/", AddStuAPI)
}
