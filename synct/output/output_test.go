package output

import (
    // "github.com/stretchr/testify/assert"
    "fmt"
    "sync"
    "testing"
)

var dbConn *DB

func TestMain(m *testing.M) {
    dbConn = NewDBConn("10.97.14.111", "9000", "test")
    m.Run()
}

func Test_Insert(t *testing.T) {
    param := make(map[string]interface{})
    var wg sync.WaitGroup
    param["date"] = "2018-03-08"
    param["time"] = 1520478473091
    param["logid"] = "00016a5259837e26f8024ac6dfa6ecc9"
    param["hostname"] = "www-laoshi-8-214"
    param["prelogid"] = ""
    param["pageid"] = ""
    param["sessid"] = ""
    param["appid"] = 1000006
    param["token"] = ""
    param["ver"] = "1"
    param["os"] = ""
    param["ua"] = "CommonFramework 1.0 rv:1 (iPhone; iOS 11.2.1; en_HK)9233"
    param["xesid"] = ""
    param["userid"] = "6279599"
    param["devid"] = ""
    param["url"] = "laoshi.xueersi.com/StudyCenter/liveOutlinesListNew"
    param["req"] = `{"url": "laoshi.xueersi.com/StudyCenter/liveOutlinesListNew", "params": {"url": "StudyCenter/liveOutlinesListNew", "stuCourseId": "30314130"}}`
    param["resp"] = `{"result": {"status": 0, "rows": 1, "data": "\u8bfe\u7a0bID\u4e0d\u5408\u6cd5"}}`
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(param map[string]interface{}) {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                if err := dbConn.Insert("xes_service_pv", param); err != nil {
                    fmt.Println(err)
                    t.Error(err)
                }
            }
        }(param)
    }
    wg.Wait()
}
