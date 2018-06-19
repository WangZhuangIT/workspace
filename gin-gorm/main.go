package main

import (
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "this is a msg",
		})
	})
	go func() {
		for {
			laddr := net.TCPAddr{net.ParseIP("127.0.0.1"), 8080, ""}
			l, err := net.ListenTCP("tcp", &laddr)
			if err != nil {
				fmt.Println(err)
			}
			fd, err := l.File()
			conn, _ := l.Accept()
			fmt.Println(conn)
			fmt.Println("*************************")
			fmt.Println(fd)
			fmt.Println(err)
			fmt.Println("*************************")

			// key := os.Getenv("LISTEN_FDS")
			// fmt.Println(key)
			time.Sleep(time.Second * 3)
		}
	}()
	route.Run()
}
