package xeslog

import (
	"fmt"
	"log"
	"os"
)

const calldepth = 3

type DefaultHandler struct {
}

func NewDefaultHandler() *DefaultHandler {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	return new(DefaultHandler)
}
func (p *DefaultHandler) close() {}

func (p *DefaultHandler) write(level int, data string) {
	log.SetPrefix(fmt.Sprintf("[%s] ", logLevelToName[level]))
	log.Output(calldepth, data)
}
