package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"synct/output"
	"time"
)

type ServicePv struct {
	Time     int64         `json:"time"`
	Logid    string        `json:"logid"`
	Hostname string        `json:"hostname"`
	Prelogid string        `json:"prelogid"`
	Pageid   string        `json:"pageid"`
	Sessid   string        `json:"sessid"`
	Appid    interface{}   `json:"appid"`
	Token    string        `json:"token"`
	Ver      string        `json:"ver"`
	Os       string        `json:"os"`
	Ua       string        `json:"ua"`
	Ip       string        `json:"ip"`
	Xesid    string        `json:"xesid"`
	Userid   string        `json:"userid"`
	Devid    string        `json:"devid"`
	Data     []interface{} `json:"data"`
	// Data []struct {
	// 	Req  map[string]interface{} `json:"req"`
	// 	Resp interface{}            `json:"resp"`
	// }
}

var wg sync.WaitGroup

func main() {
	// logFile, err := os.OpenFile("/tmp/err_line.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// defer logFile.Close()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.CommitInterval = 5
	config.Group.Topics.Whitelist = regexp.MustCompile(`flume-service-pv-.*`)
	// config.Group.Heartbeat.Interval = 1

	brokers := []string{"10.10.9.107:9092", "10.10.9.108:9092", "10.10.9.109:9092", "10.10.9.181:9092", "10.10.9.182:9092"}
	// brokers := []string{"10.97.14.111:9092", "10.97.14.112:9092", "10.97.14.112:9092"}
	// topics := []string{"flume-test-.*"}

	consumer, err := cluster.NewConsumer(brokers, "consumer-go", nil, config)
	if err != nil {
		panic(err)
	}

	dbConn := output.NewDBConn("10.10.142.142", "9000", "beta")
	// dbConn := output.NewDBConn("10.97.14.111", "9000", "beta")

	defer consumer.Close()

	signals := make(chan os.Signal, 10)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for err := range consumer.Errors() {
			fmt.Println("consumer error ", err)
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			fmt.Println("consumer not ", ntf)
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go SaveMsg(consumer, signals, dbConn)

	}
	wg.Wait()
}

func SaveMsg(consumer *cluster.Consumer, signals chan os.Signal, dbConn *output.DB) {
	var param []map[string]interface{}
	defer wg.Done()
	for {
		select {
		case <-signals:
			fmt.Println(1)
			return
		case line, ok := <-consumer.Messages():
			if ok {
				msg, err := GetMsg(line.Value)
				if err != nil {
					fmt.Println(err, string(line.Value))
				} else {
					data, err := DealMsg(msg)
					if err != nil {
						fmt.Println(err, string(line.Value))
					} else {
						// TODO data map => slice
						if len(param) < 100 {
							param = append(param, data)
						} else {
							if err = dbConn.ServicePv(param); err != nil {
								fmt.Println(err)
							}
							param = param[:0]
						}
					}
				}
				consumer.MarkOffset(line, "")
			}
		}
	}
}

func GetMsg(line []byte) ([]byte, error) {
	var msg []byte
	index := strings.Index(string(line), "{\"time")
	if index == -1 {
		return nil, errors.New("line error")
	}

	msg = line[index:]

	return msg, nil
}

func DealMsg(log []byte) (map[string]interface{}, error) {
	content := &ServicePv{}
	var param map[string]interface{}
	// var ok bool
	err := json.Unmarshal(log, &content)
	if err != nil {
		return nil, err
	} else {
		param = make(map[string]interface{})

		switch appid := content.Appid.(type) {
		case float64:
			param["appid"] = uint64(appid)
		case string:
			if param["appid"], err = strconv.ParseUint(appid, 10, 64); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknow type appid")
		}

		param["url"] = ""
		param["data"] = ""

		for _, val := range content.Data {
			if data, ok := val.(map[string]interface{}); ok {
				if _, ok = data["req"]; ok {
					if req, ok := data["req"].(map[string]interface{}); ok {
						if _, ok = req["url"]; ok {
							if url, ok := req["url"].(string); ok {
								param["url"] = url
								break
							}
						}
					}
				}
			}
		}

		if data, err := json.Marshal(content.Data); err != nil {
			return nil, err
		} else {
			param["data"] = string(data)
		}

		param["date"] = time.Now().Format("2006-01-02")
		param["time"] = content.Time
		param["logid"] = content.Logid
		param["hostname"] = content.Hostname
		param["prelogid"] = content.Prelogid
		param["pageid"] = content.Pageid
		param["sessid"] = content.Sessid
		param["token"] = content.Token
		param["ver"] = content.Ver
		param["os"] = content.Os
		param["ua"] = content.Ua
		param["xesid"] = content.Xesid
		param["userid"] = content.Userid
		param["devid"] = content.Devid

		return param, nil
	}
}
