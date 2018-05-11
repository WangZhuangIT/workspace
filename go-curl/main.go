package main

import (
	"encoding/json"
	"fmt"

	"github.com/mikemintang/go-curl"
)

type ServerReq struct {
	Success int   `json:"success"`
	Data    Datas `json:"data"`
}

type Datas struct {
	Reqlist []Reqlist `json:"reqlist"`
}

type Reqlist struct {
	AppID      int    `json:"app_id"`
	ServerName string `json:"server_name"`
	Request    string `json:"request"`
	Tag        string `json:"tag"`
}

type Person struct {
	Age  int
	Name string
}

func main11() {

	url := "http://gwadmin.beta.xesv5.com/monitor/request/all"

	// headers := map[string]string{
	// 	"User-Agent":    "Sublime",
	// 	"Authorization": "Bearer access_token",
	// 	"Content-Type":  "application/json",
	// }

	// cookies := map[string]string{
	// 	"userId":    "12",
	// 	"loginTime": "15045682199",
	// }

	// queries := map[string]string{
	// 	"page": "2",
	// 	"act":  "update",
	// }

	// postData := map[string]interface{}{
	// 	"name":      "mike",
	// 	"age":       24,
	// 	"interests": []string{"basketball", "reading", "coding"},
	// 	"isAdmin":   true,
	// }

	// 链式操作
	req := curl.NewRequest()
	resp, err := req.
		SetUrl(url).
		// SetHeaders(headers).
		// SetCookies(cookies).
		// SetQueries(queries).
		// SetPostData(postData).
		Get()

	if err != nil {
		fmt.Println(err)
	} else {
		if resp.IsOk() {
			var r ServerReq
			err := json.Unmarshal([]byte(resp.Body), &r)
			if err != nil {
				fmt.Println("err")
				return
			}
			fmt.Println(r)

			for index, data := range r.Data.Reqlist {
				fmt.Println(index)
				fmt.Println(data.AppID)
				fmt.Println(data.Request)
				fmt.Println(data.ServerName)
				fmt.Println(data.Tag)
			}

			person := Person{}
			person.Age = 17
			person.Name = "wangdazhuang"

			result, err := json.Marshal(person)
			if err != nil {
				fmt.Println("chu cuo le")
			}
			fmt.Println("kaishile ********************************************")
			fmt.Println(string(result))

		} else {
			fmt.Println(resp.Raw)
		}
	}

}
