package api

import (
	. "gin-xorm/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddStuAPI(ctx *gin.Context) {
	stu := Student{}
	stu.Name = "wangdazhuang"
	stu.Sex = "man"
	stu.Age = 16
	stu.Address = "河北省保定市"
	err := stu.AddStuModel()
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
