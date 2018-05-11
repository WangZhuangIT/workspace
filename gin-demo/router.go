package main

import (
	. "gin-demo/apis"
	. "gin-demo/middleware"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func initRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/index", IndexApi)
	//增加
	router.POST("/person", AddPersonApi)

	//查询
	router.GET("/person", QueryPersonApi)
	router.GET("/person/:id", QueryPersonApiById)

	//更新
	router.PUT("/person/:id", UpdatePersonApi)

	//删除
	router.DELETE("/person/:id", DelPersonApi)

	router.POST("/form", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
	})

	//路由组
	someGroup := router.Group("/someGroup")
	{
		someGroup.GET("/someGet", func(c *gin.Context) {
			var store = sessions.NewCookieStore([]byte(""))
			session, _ := store.Get(c.Request, "wz")
			fmt.Println(session)
			fmt.Println("get group")
		})
		someGroup.POST("/somePost", func(c *gin.Context) {
			fmt.Println("post group")
		})
	}

	// 绑定JSON的例子 ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login

		if c.BindJSON(&json) == nil {
			if json.User == "123" && json.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	// 绑定普通表单的例子 (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// 根据请求头中 content-type 自动推断.
		if c.Bind(&form) == nil {
			if form.User == "123" && form.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	router.POST("/upload", func(c *gin.Context) {

		file, header, err := c.Request.FormFile("upload")
		filename := header.Filename
		fmt.Println(header.Filename)
		out, err := os.Create("." + filename + ".png")
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}

	})

	//加载模板
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	//定义路由
	router.GET("/indexs", func(c *gin.Context) {
		//根据完整文件名渲染模板，并传递参数
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/redirect", Logger(), func(c *gin.Context) {
		//支持内部和外部的重定向
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
		fmt.Println(123)
	})

	router.GET("/mongo", MongoAPI)

	router.POST("/mid", MiddleWare(), func(c *gin.Context) {
		fmt.Println("mid aaa data")
	})

	router.GET("/version/:id", MD5API)

	return router
}
