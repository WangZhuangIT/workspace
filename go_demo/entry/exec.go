package main

import (
	"runtime"

	"go_libs/service/logger"
)

type Configure struct {
	Http string
	FDS  struct {
		AppKey    string
		AppSecret string
	}
}

func exec(quit chan struct{}) {
	defer recoverPanic()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	ctrl := NewServer()
	for {
		select {
		case <-quit:
			if ctrl != nil {
				ctrl.Close()
				logger.I("SERVER_FINISHED", "shutdown")
			}
			return
		}
	}
}
func recoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(error)
		logger.E("PANIC_RECOVERED", "Unhandled error: %v\n", err.Error())
	}
}
