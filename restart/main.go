package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sign := <-signs
		fmt.Print("接收到了信号：")
		fmt.Print(sign)
		startBg()
		done <- true
	}()

	fmt.Println("waitting for you , sign")
	pid := os.Getpid()
	fmt.Printf("The process id is %v", pid)
	<-done
	fmt.Println("exiting")
}

func startBg() {
	env := os.Environ()
	procAttr := &os.ProcAttr{
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	pid, err := os.StartProcess("/usr/local/go/bin/go", []string{"go", "run", "main.go"}, procAttr)
	if err != nil {
		fmt.Printf("Error %v starting process!", err) //
		os.Exit(1)
	}
	fmt.Printf("The process id is %v", pid)
}
