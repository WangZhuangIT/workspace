package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "github.com/coreos/etcd/clientv3"
)

func main111() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.99.1.151:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("cuole ")
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp_put, err := cli.Put(ctx, "sample_key", "sample_value")
	cancel()
	if err != nil {
		fmt.Println("a ...interface{}")
		log.Fatalln(err)
	}

	fmt.Println(resp_put)

	//设置1秒超时，访问etcd有超时控制
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	resp_get, err := cli.Get(ctx, "sample_key")
	//操作完毕，取消etcd
	cancel()

	if err != nil {
		fmt.Println("get error")
		log.Fatalln(err)
	}

	for _, ev := range resp_get.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
