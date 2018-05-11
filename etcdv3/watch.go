package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.99.1.151:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("连接失败", err)
		return
	}

	fmt.Println("连接成功")
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cli.Put(ctx, "nihao", "hahhaa")
	cancel()

	for {
		rch := cli.Watch(context.Background(), "nihao")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
