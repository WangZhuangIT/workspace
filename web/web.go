package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	time.Sleep(5 * time.Second)
	fmt.Println("ni hao a xiong di")
	http.HandleFunc("/", son)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatalln(err)
	}

}

func son(w http.ResponseWriter, r *http.Request) {

	fmt.Println("this is a msg")
}
