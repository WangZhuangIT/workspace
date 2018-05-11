package main

import "github.com/gin-gonic/gin"

func main() {
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "this is a msg",
		})
	})
	route.Run()

}
