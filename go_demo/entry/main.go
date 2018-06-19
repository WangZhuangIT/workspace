package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"go_libs/service/logger"
	"go_libs/utils/confutil"
	"go_libs/utils/flagutil"

	//"git.n.xiaomi.com/miot_shop/go_libs/util/swaggergen"
	"github.com/alecthomas/log4go"
)

var (
	exec_do    = make(chan struct{}, 0)
	done       = make(chan struct{}, 0)
	sigstop    = flagutil.GetSignal()
	foreground = flagutil.GetForeground()
	//gendoc     = flagutil.GetGendoc()
	g_cfg   *Configure
	g_cntxt *daemon.Context
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		version()
		return
	}
	log4go.LoadConfiguration("conf/log.xml")
	confutil.InitConfig()

	if *foreground {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			logger.I("QUIT", "terminating...")
			close(exec_do)
		}()
		worker()
		logger.I("QUIT", "terminated")
		fmt.Println("QUIT", "terminated")
		return
	}

	/*if gendoc {
		curpath, _ := os.Getwd()
		swaggergen.GenerateDocs(curpath)
	}*/
	daemon.AddCommand(daemon.StringFlag(sigstop, "stop"), syscall.SIGQUIT, termHandler)
	binhome := filepath.Dir(filepath.Dir(gRealPath(os.Args[0])))
	os.Mkdir(binhome+"/var", 0775)
	os.Mkdir(binhome+"/data", 0775)

	g_cntxt = &daemon.Context{
		PidFileName: binhome + "/var/pid",
		PidFilePerm: 0644,
		LogFileName: binhome + "/var/stdout." + strconv.Itoa(os.Getpid()),
		LogFilePerm: 0644,
		WorkDir:     binhome,
		Umask:       027,
		Args:        os.Args,
	}

	if len(daemon.ActiveFlags()) > 0 {
		if d, err := g_cntxt.Search(); err != nil {
			fmt.Println("Unable send signal to the daemon:", err)
		} else {
			daemon.SendCommands(d)
		}
		return
	}

	if d, err := g_cntxt.Reborn(); err != nil {
		logger.E("Reborn:", err)
	} else if d != nil {
		return
	}
	defer g_cntxt.Release()

	logger.I("MAIN", "daemon started")

	go worker()

	if err := daemon.ServeSignals(); err != nil {
		logger.E("MAIN_WAIT", err)
	}
	logger.I("MAIN", "daemon terminated")
}

func worker() {
	exec(exec_do)
	close(done)
}

func termHandler(sig os.Signal) error {
	logger.I("MAIN", "terminating...")
	close(exec_do)
	if sig == syscall.SIGQUIT {
		//TODO release
		g_cntxt.Release()
		<-done
	}
	return daemon.ErrStop
}

func gRealPath(path string) string {
	if path[0] == '~' {
		home := os.Getenv("HOME")
		path = home + path[1:]
	}
	rpath, err := filepath.Abs(path)
	if err == nil {
		path = rpath
	}
	return strings.TrimSpace(path)
}
