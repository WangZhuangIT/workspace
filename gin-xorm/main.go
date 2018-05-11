package main

import (
	. "gin-xorm/config"
	. "gin-xorm/dao"
	. "gin-xorm/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	InitConfigManager("config/conf")
	InitDb()
	gin.SetMode("debug")
	e := gin.New()
	InitRoutes(e)
	e.Run(":8180")
}
