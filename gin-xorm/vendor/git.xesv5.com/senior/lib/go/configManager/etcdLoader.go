package config

import (
	"context"
	"strings"
	"time"

	"git.xesv5.com/senior/lib/go/xeslog"
	"github.com/coreos/etcd/clientv3"
)

type etcdLoader struct {
	dir string
	cli *clientv3.Client
}

func NewEtcdLoader(endPoints []string, dir string) *etcdLoader {
	e := new(etcdLoader)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		xeslog.Fatal(err)
	}
	e.cli = cli
	e.dir = strings.Trim(dir, "/")
	return e
}

func (e *etcdLoader) Read() (map[string][]byte, error) {
	//获取所有etcd的keys，设置value
	data := make(map[string][]byte)
	var err error
	res, err := e.cli.Get(context.Background(), e.dir+"/", clientv3.WithPrefix())
	if err != nil {
		return data, err
	}
	for _, kv := range res.Kvs {
		data[strings.Replace(string(kv.Key), e.dir+"/", "", -1)] = kv.Value
	}
	return data, err
}
func (e *etcdLoader) Watch(onChange func(data map[string][]byte)) {
	rch := e.cli.Watch(context.Background(), e.dir+"/", clientv3.WithPrefix())
	for ev := range rch {
		if len(ev.Events) > 0 {
			xeslog.Debug(ev.Events)
			data, err := e.Read()
			xeslog.Debug(data)
			if err != nil {
				xeslog.Error(err)
			} else {
				onChange(data)
			}
		}
	}
}
