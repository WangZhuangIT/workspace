package controller

import (
	"go_demo/model"
	"go_libs/dao/redisdao"
	"go_libs/service/logger"
	"go_libs/utils/confutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Entry struct {
}

func init() {
	this := new(Entry)
	entry := Gin.Group("/entry")
	entry.GET("/test", this.entryTest)
	logger.I("DEBUG", confutil.GetConf("DEFAULT", "word"))
}
func (this *Entry) entryTest(ctx *gin.Context) {
	logger.I("DEBUG", "entry test")
	student := model.Student{}
	students, err := student.StuInfoList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Error(http.StatusInternalServerError, err))
	} else {
		//ctx.JSON(http.StatusOK, Success(students))
		key := "demo-student"
		logger.I("DEBUG", "Student:%v", students)
		cache := redisdao.NewRedisCache("demo/cache")
		str, _ := cache.Id(key).Get().String()
		ctx.JSON(http.StatusOK, Success(str))
	}
	return
}
