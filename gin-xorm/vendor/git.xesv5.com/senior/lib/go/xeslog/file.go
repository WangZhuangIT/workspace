package xeslog

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type FileHandler struct {
	path   string
	prefix string
	files  map[int]xesFile
	lock   *sync.Mutex
	app    string
}

type xesFile struct {
	fd    *os.File
	fname string
}

func NewXesFileHandler(path, prefix, app string) *FileHandler {
	this := new(FileHandler)
	this.path = path
	this.prefix = prefix
	this.lock = new(sync.Mutex)
	this.app = app
	this.files = make(map[int]xesFile)
	return this
}

func (this *FileHandler) open(level int) error {
	level_name := logLevelToName[level]
	fname := fmt.Sprintf("%s%s%s.log", this.path, this.prefix, level_name)
	need_open := false
	if v, ok := this.files[level]; !ok {
		need_open = true
	} else if _, err := os.Stat(v.fname); os.IsNotExist(err) {
		need_open = true
	}
	if need_open {
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		this.files[level] = xesFile{
			fd:    f,
			fname: fname,
		}
	}
	return nil
}

func (this *FileHandler) write(level int, data string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	err := this.open(level)
	if err != nil {
		log.Println(err)
		return
	}
	data = strings.Trim(data, "\n")
	t := time.Now().Format("2006-01-02 15:04:05")
	str := fmt.Sprintf("[%s]\x01 %s\x01 %s\x01 %s\x01 %s\x01 %d\x01 %s\x01 %s\x01\n", t, this.app, logLevelToName[level], getLocalIP(), file, line, runtime.FuncForPC(pc).Name(), data)
	f := this.files[level]
	_, err = f.fd.Write([]byte(str))
	if err != nil {
		log.Printf("write file %s error %s", f.fname, err.Error())
	}
}

func (this *FileHandler) close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, v := range this.files {
		if v.fd != nil {
			v.fd.Close()
		}
		delete(this.files, k)
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
