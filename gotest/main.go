package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	. "git.xesv5.com/senior/lib/go/xesSdk"
	"git.xesv5.com/senior/oneToOneBackend/config"
)

func main() {

	test()
	// testPost()
	// testReq()
	test()
	fmt.Println(999)

}

func testConfig() {
	dir := "./dev"
	config.Init(dir, nil)
	fmt.Println(config.GetBasicConfig().Appid)
	fmt.Println(123)
}

func testSdk() {
	username := "xuelin"
	pwd := "Q7@n~2sZ"
	appid := 44
	var rep PwdLoginResp
	var req PwdLoginReq = PwdLoginReq{Username: username, Pwd: pwd, Appid: appid}
	var sdk SdkMethods
	sdk = new(DefaultHandler)
	err := sdk.AdminPermission("pwdLogin", req, &rep)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println(rep.TeacherId)
}

func test() {
	files, err := ioutil.ReadDir("../synct")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func testReq() {
	username := "xuelin"
	pwd := "Q7@n~2sZ"
	appid := "44"
	dataurl := "http://adminapi.xesv5.com/Admin/Login/checkPwd"
	//表单的数据
	v := url.Values{}
	v.Set("username", username)
	v.Set("pwd", pwd)
	v.Set("appid", appid)
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))

	client := &http.Client{}

	//利用指定的method,url以及可选的body返回一个新的请求.如果body参数实现了io.Closer接口，Request返回值的Body 字段会被设置为body，并会被Client类型的Do、Post和PostFOrm方法以及Transport.RoundTrip方法关闭。
	req, err := http.NewRequest("POST", dataurl, body)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	//必须设定该参数，POST请求才能正常提交
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))

	pc, _, _, _ := runtime.Caller(1)
	funFullname := runtime.FuncForPC(pc).Name()
	fmt.Println(funFullname)
}

func testPost() {
	username := "xuelin"
	pwd := "Q7@n~2sZ"
	appid := "44"
	dataurl := "http://adminapi.xesv5.com/Admin/Login/checkPwd"
	resp, err := http.PostForm(dataurl, url.Values{"appid": {appid}, "pwd": {pwd}, "username": {username}})

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer resp.Body.Close()

	buf := make([]byte, 4*1024)
	var result string
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}

	fmt.Println(result)
}
