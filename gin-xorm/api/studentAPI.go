package api

import (
	"gin-xorm/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddStuAPI(ctx *gin.Context) {
	name := ctx.Request.FormValue("name")
	sex := ctx.Request.FormValue("sex")
	age, _ := strconv.Atoi(ctx.Request.FormValue("age"))
	addr := ctx.Request.FormValue("addr")
	stu := models.Student{Name: name, Sex: sex, Age: age, Address: addr}
	err := stu.AddStu()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 1,
			"msg":  "插入成功",
		})
	}
}

func GetStuAPI(ctx *gin.Context) {
	studentModel := models.Student{}
	err := studentModel.GetStu()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 1,
			"msg":  "插入成功",
		})
	}
}
