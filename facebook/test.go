package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("123")
	wd, _ := os.Getwd()
	fmt.Println(wd)
	// envCountKey := "LISTEN_FDS"
	fmt.Printf("%s lives in %s.\n", os.Getenv("USER"), os.Getenv("HOME"))
}
