package api

import (
	"fmt"
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
	studentModel := models.Student{Name: name, Sex: sex, Age: age, Address: addr}
	err := studentModel.AddStu()
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
	var id int
	var err error
	if id, err = strconv.Atoi(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  "id must be int",
		})
	}

	studentModel := models.Student{Id: id}
	err = studentModel.GetStu()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 1,
			"msg":  studentModel,
		})
	}
}

func GetStuSliceAPI(ctx *gin.Context) {
	studentModel := models.Student{}
	stuSlice, err := studentModel.GetStuSlice(1, 10)
	fmt.Println(stuSlice)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 1,
			"msg":  stuSlice,
		})
	}
}

func UpStuAPI(ctx *gin.Context) {
	studentModel := models.Student{Age: 24}
	err := studentModel.UpStu()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 1,
			"msg":  "更新成功",
		})
	}
}

func DelStuAPI(ctx *gin.Context) {
	studentModel := models.Student{Age: 44}
	err := studentModel.DelStu()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 0,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"sign": 1,
			"msg":  "删除成功",
		})
	}
}
