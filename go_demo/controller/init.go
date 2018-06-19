package controller

import (
	"fmt"
	"go_libs/service/logger"
	"go_libs/utils/errorutil"

	"github.com/gin-gonic/gin"
)

var Gin = gin.New()

func init() {
	fmt.Println("Gin init")
}

func Success(v interface{}) map[string]interface{} {
	ret := map[string]interface{}{"code": 0, "msg": "ok", "data": v}
	return ret
}

func Error(code int, msg error) map[string]interface{} {
	if code == 0 {
		er := logger.NewError(msg)
		code = er.Code
		if code == 0 {
			code = int(errorutil.ERROR_SERVER)
		}
	}
	ret := map[string]interface{}{"code": code, "msg": msg.Error()}
	return ret
}
