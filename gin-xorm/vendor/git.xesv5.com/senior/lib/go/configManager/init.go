package config

import (
	"sync"
)

var (
	lock    sync.RWMutex
	dataMap map[string][]byte
)

type readChan struct {
}

type LoaderIF interface {
	Read() (map[string][]byte, error)
	Watch(onChange func(data map[string][]byte))
}

func Init(dir string, load LoaderIF) {
	if load == nil {
		load = NewFileLoader(dir)
	}
	dataMap, err := load.Read()
	if err != nil {
		panic(err)
	}
	loadConfig(dataMap)
	go load.Watch(loadConfig)
}

func loadConfig(data map[string][]byte) {
	lock.Lock()
	defer lock.Unlock()
	dataMap = data
}

func GetConfigMap() map[string][]byte {
	lock.RLock()
	defer lock.RUnlock()
	return dataMap
}
