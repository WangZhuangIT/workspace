package logger

import (
	"context"
	"fmt"
	"os"
	//"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	//"go_libs/dao/confdao"
	//"go_libs/utils/kafkautil"

	"github.com/alecthomas/log4go"
	//"github.com/kardianos/osext"
)

var hostname string

//var env string = confdao.GetEnv()
//var Log_local bool = confdao.LogLocal()
//var level confdao.Level = confdao.GetLogLevel()
var STDOUT bool

func init() {
}

func I(tag string, arg0 interface{}, v ...interface{}) {
	loggerX(nil, 'I', tag, arg0, v...)
}
func Ix(ctx context.Context, tag string, arg0 interface{}, v ...interface{}) {
	loggerX(ctx, 'I', tag, arg0, v...)
}

func D(tag string, arg0 interface{}, v ...interface{}) {
	loggerX(nil, 'D', tag, arg0, v...)
}

func Dx(ctx context.Context, tag string, arg0 interface{}, v ...interface{}) {
	loggerX(ctx, 'D', tag, arg0, v...)
}

func W(tag string, arg0 interface{}, v ...interface{}) {
	loggerX(nil, 'W', tag, arg0, v...)
}

func Wx(ctx context.Context, tag string, arg0 interface{}, v ...interface{}) {
	loggerX(ctx, 'W', tag, arg0, v...)
}

func E(tag string, arg0 interface{}, v ...interface{}) {
	loggerX(nil, 'E', tag, arg0, v...)
}

func Ex(ctx context.Context, tag string, arg0 interface{}, v ...interface{}) {
	loggerX(ctx, 'E', tag, arg0, v...)
}

func C(tag string, arg0 interface{}, v ...interface{}) {
	loggerX(nil, 'C', tag, arg0, v...)
}

func M(tag string, arg0 interface{}, v ...interface{}) {
	loggerX(nil, 'M', tag, arg0, v...)
}

func Close() {
	log4go.Close()
}

func loggerX(ctx context.Context, lvl byte, tag string, arg0 interface{}, v ...interface{}) {
	if ctx == nil {
		id := strconv.Itoa(Id())
		ctx = context.WithValue(context.Background(), "logid", id)
	}
	logid, ok := ctx.Value("logid").(string)
	if !ok {
		logid = strconv.Itoa(Id())
	}
	if tag == "" {
		tag = "NoTagError"
	}
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	var errMessage string
	var position string
	switch t := arg0.(type) {
	case *StackErr:
		errMessage = t.ErrorMessage
		position = t.Position
	case error:
		errMessage = t.Error()
	case string:
		if len(v) > 0 {
			errMessage = fmt.Sprintf(t, v...)
		} else {
			errMessage = t
		}
	default:
		errMessage = fmt.Sprint(t)
	}

	if position == "" {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			position = filepath.Base(file) + ":" + strconv.Itoa(line) + ":" + logid
		} else {
			position = "EMPTY"
		}
	}
	tag = Filter(tag)
	errMessage = Filter(errMessage)
	//n := time.Now().Format("2006/01/02 15:04:05")

	//	if env == "dev" || Log_local {
	switch lvl {
	case 'D':
		log4go.Log(log4go.DEBUG, position, tag+"\t"+errMessage)
	case 'I':
		if start, ok := ctx.Value("start").(time.Time); ok {
			cost := time.Now().Sub(start)
			errMessage = errMessage + " COST:" + cost.String()
		}
		log4go.Log(log4go.INFO, position, tag+"\t"+errMessage)
	case 'W':
		log4go.Log(log4go.WARNING, position, tag+"\t"+errMessage)
	case 'E':
		log4go.Log(log4go.ERROR, position, tag+"\t"+errMessage)
	case 'C':
		log4go.Log(log4go.CRITICAL, position, tag+"\t"+errMessage)

	}
	/*	} else {
		program, _ := osext.Executable()
		msg := "alog" + "\t" + path.Base(program) + "\t" + n + "\t" + string(lvl) + "\t" + hostname + "\t" + position + "\t" + tag + "\t" + errMessage
		switch lvl {
		case 'D':
			sendLogByKafka(confdao.DEBUG, logid, msg)
		case 'I':
			perf.IncI(tag)
			sendLogByKafka(confdao.INFO, logid, msg)
		case 'W':
			perf.IncW(tag)
			sendLogByKafka(confdao.WARNING, logid, msg)
		case 'E':
			perf.IncE(tag)
			sendLogByKafka(confdao.ERROR, logid, msg)
		case 'C':
			sendLogByKafka(confdao.CRITICAL, logid, msg)
		}
	}*/
}

/*
func sendLogByKafka(lvl confdao.Level, logid, msg string) {
	if lvl >= level {
		if STDOUT {
			fmt.Println(msg, logid)
		} else {
			if err := kafkautil.Send2Proxy("miotstore_alog", msg, logid); err != nil {
				mlog.E("KafkaProxyError", err)
			}
		}
	}
}
*/
var defaultReplacer *strings.Replacer

func init() {
	defaultReplacer = strings.NewReplacer("\t", "", "\r", "", "\n", "")
}

func Filter(msg string, r ...string) string {
	replacer := defaultReplacer
	if len(r) > 0 {
		replacer = strings.NewReplacer("\t", r[0], "\r", r[0], "\n", r[0])
	}
	return replacer.Replace(msg)
}
