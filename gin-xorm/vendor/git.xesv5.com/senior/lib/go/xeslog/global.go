package xeslog

import (
	"fmt"
	"os"
)

var (
	CallDepth      int = 4
	logLevelToName map[int]string
	logger         LoggerIF
	level          int
)

type LoggerIF interface {
	write(int, string)
	close()
}

const (
	LOGLEVEL_SILENT = -1
	LOGLEVEL_FATAL  = iota
	LOGLEVEL_ERROR
	LOGLEVEL_ALERT
	LOGLEVEL_WARN
	LOGLEVEL_INFO
	LOGLEVEL_DEBUG
	NR_LOGLEVELS
)

func init() {
	logLevelToName = make(map[int]string, NR_LOGLEVELS)
	logLevelToName[LOGLEVEL_DEBUG] = "DEBUG"
	logLevelToName[LOGLEVEL_INFO] = "INFO"
	logLevelToName[LOGLEVEL_WARN] = "WARNING"
	logLevelToName[LOGLEVEL_ALERT] = "ALERT"
	logLevelToName[LOGLEVEL_ERROR] = "ERROR"
	logLevelToName[LOGLEVEL_FATAL] = "FATAL"
}

func Set(logLevel int, handler LoggerIF) {
	if handler == nil {
		handler = NewDefaultHandler()
	}
	level = logLevel
	logger = handler
}

func Getlevel() int {
	return level
}

func Debug(v ...interface{}) {
	if level < LOGLEVEL_DEBUG {
		return
	}
	logger.write(LOGLEVEL_DEBUG, fmt.Sprintln(v...))
}

func Debugf(format string, v ...interface{}) {
	if level < LOGLEVEL_DEBUG {
		return
	}
	logger.write(LOGLEVEL_DEBUG, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	if level < LOGLEVEL_INFO {
		return
	}
	logger.write(LOGLEVEL_INFO, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	if level < LOGLEVEL_INFO {
		return
	}
	logger.write(LOGLEVEL_INFO, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	if level < LOGLEVEL_WARN {
		return
	}
	logger.write(LOGLEVEL_WARN, fmt.Sprintln(v...))

}

func Warnf(format string, v ...interface{}) {
	if level < LOGLEVEL_WARN {
		return
	}
	logger.write(LOGLEVEL_WARN, fmt.Sprintf(format, v...))
}
func Alert(v ...interface{}) {
	if level < LOGLEVEL_ALERT {
		return
	}
	logger.write(LOGLEVEL_ALERT, fmt.Sprintln(v...))
}

func Alertf(format string, v ...interface{}) {
	if level < LOGLEVEL_ALERT {
		return
	}
	logger.write(LOGLEVEL_ALERT, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	if level < LOGLEVEL_ERROR {
		return
	}
	logger.write(LOGLEVEL_ERROR, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	if level < LOGLEVEL_ERROR {
		return
	}
	logger.write(LOGLEVEL_ERROR, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	if level < LOGLEVEL_FATAL {
		return
	}
	logger.write(LOGLEVEL_FATAL, fmt.Sprintln(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if level < LOGLEVEL_FATAL {
		return
	}
	logger.write(LOGLEVEL_FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}
