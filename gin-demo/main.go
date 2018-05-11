package main

import (
	db "gin-demo/database"
)

func main() {
	defer db.SqlDB.Close()
	router := initRouter()
	router.Run(":8083")
}
