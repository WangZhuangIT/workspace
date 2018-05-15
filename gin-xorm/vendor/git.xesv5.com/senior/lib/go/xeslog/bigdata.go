package xeslog

type BigDataShowLog struct {
	Timestamp int64       `json:"timestamp"`
	Logid     string      `json:"logid"`
	Prelogid  string      `json:"prelogid"`
	Pageid    string      `json:"pageid"`
	Sessid    string      `json:"sessid"`
	Appid     string      `json:"appid"`
	Os        string      `json:"os"`
	Ua        string      `json:"ua"`
	Xesid     string      `json:"xesid"`
	Userid    int64       `json:"userid"`
	Data      interface{} `json:"data"`
}
