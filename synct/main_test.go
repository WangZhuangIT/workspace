package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"regexp"
	"synct/output"
	"testing"
)

func Test_DealMsg(t *testing.T) {
	line := `{"time":1519980185217,"logid":"b7d48713d7cf0565bffe268e52359d6f","hostname":"www-laoshi-9-91","prelogid":"","pageid":"","sessid":"7adc44caedb6e94c07cb7bcad04aae7e9f94f4x49b","appid":"1000031","token":"","ver":"1","os":"","ua":"CommonFramework 1.0 rv:1 (iPad; iOS 10.3.3; zh_CN)","ip":"36.106.19.54","xesid":"ab79cc49c5d9dee013d568ac400f5307","userid":"4307179","devid":"","data":[{"req":{"url":"laoshi.xueersi.com\/Task\/getUserTaskInfo","params":{"url":"Task\/getUserTaskInfo","requesttime":"1519980185139","appVersionNumber":"61032","datalogId":"4672482760463439780","logId":"1136880558568734575","systemName":"iOS","appVersion":"6.3.02","systemVersion":"10.3.3","sessid":"7adc44caedb6e94c07cb7bcad04aae7e9f94f4x49b","identifierForClient":"450AAFF7-7D2C-4B2A-BF50-44E31AFB8280","deviceModel":"iPad"}}},{"resp":{"result":{"status":1,"rows":1,"data":{"pendingNum":0,"status":0}}}}]}`
	log := []byte(line)
	// _, err := DealMsg(log)
	// if err != nil {
	// 	t.Error(err)
	// }
	DealMsg(log)
}

func Test_SaveMsg(t *testing.T) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Group.Topics.Whitelist = regexp.MustCompile(`flume-service-pv-.*`)
	// config.Group.Heartbeat.Interval = 1

	brokers := []string{"10.97.14.111:9092", "10.97.14.112:9092", "10.97.14.112:9092"}
	// topics := []string{"flume-test-.*"}

	consumer, err := cluster.NewConsumer(brokers, "test", nil, config)
	if err != nil {
		panic(err)
	}
	dbConn := output.NewDBConn("10.97.14.111", "9000", "test")

	defer consumer.Close()

	signals := make(chan os.Signal, 1)
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

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go SaveMsg(consumer.Messages(), signals, dbConn)

	}
	wg.Wait()
}
