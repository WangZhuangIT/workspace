package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name    string `json:"Name"`
	Sex     string `json:"Sex"`
	Age     int    `json:"Age"`
	Address string `json:"Address"`
}

func main() {
	stu := new(Student)
	stu.Name = "wangdazhuang"
	stu.Sex = "man"
	stu.Age = 16
	stu.Address = "河北省保定市"

	data, _ := json.Marshal(stu)
	jdata := string(data)
	fmt.Println(jdata)

	std := new(Student)
	json.Unmarshal([]byte(jdata), &std)
	fmt.Println(std.Address)
}
