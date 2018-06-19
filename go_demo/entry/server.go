package main

import (
	"go_demo/controller"
	"go_libs/utils/confutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/alecthomas/log4go"
	"github.com/tylerb/graceful"
)

type Server struct {
	srv *graceful.Server
}

func NewServer() *Server {
	this := new(Server)
	//mux := http.NewServeMux()
	/*for index, handle := range controller.ActionsMap {
		mux.HandleFunc(index, handle.Exec)
	}*/
	// mux.Handle("/api/", http.FileServer(http.Dir("swagger")))
	go func() {
		http.Handle("/memprofile", new(MemProf))
		http.Handle("/cpuprofile", new(CpuProf))
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	port := confutil.GetConf("HTTP", "port")
	if len(port) == 0 {
		log4go.Crash("INVALID_PORT")
	}
	this.srv = &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    ":" + port,
			Handler: controller.Gin,
		},
	}
	go this.srv.ListenAndServe()
	return this
}

func (this *Server) Close() {
	go this.srv.Stop(2 * time.Second)
	<-this.srv.StopChan()
}

type MemProf struct {
}

type CpuProf struct {
}

func (this *MemProf) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := "mem-" + strconv.Itoa(os.Getpid()) + ".memprof"
	f, _ := os.Create(filename)
	pprof.WriteHeapProfile(f)
	f.Close()
}
func (this *CpuProf) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	seconds, _ := strconv.Atoi(r.FormValue("seconds"))
	filename := "cpu-" + strconv.Itoa(os.Getpid()) + ".cupprof"
	f, _ := os.Create(filename)
	err := pprof.StartCPUProfile(f)
	if err != nil {
		log4go.Crash("CpuProf:", err)
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	pprof.StopCPUProfile()
	f.Close()
}
