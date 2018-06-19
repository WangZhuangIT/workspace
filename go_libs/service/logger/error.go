package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"go_libs/utils/errorutil"
	"go_libs/utils/textutil"
	"github.com/spf13/cast"
)

type StackErr struct {
	Filename      string
	CallingMethod string
	Line          int
	ErrorMessage  string
	StackTrace    string
	Code          int
	Info          string
	Tag           string
	Position      string
	Data          interface{}
}

const (
	MIOT_DB_ERROR     int = 500010001
	MIOT_JSON_ERROR   int = 500020001
	MIOT_REDIS_ERROR  int = 500030001
	MIOT_ES_ERROR     int = 500040001
	MIOT_THRIFT_ERROR int = 500050001
)

type ERROR int

const (
	ERROR_INVALID_PARAM ERROR = 401000 + iota
	ERROR_INVALID_CARD
	ERROR_INVALID_USER
	ERROR_INVALID_GOODS
	ERROR_INVALID_NOTSTART
	ERROR_NETWORK

	ERROR_DB     ERROR = 402000 + iota
	ERROR_JSON   ERROR = 403000 + iota
	ERROR_REDIS  ERROR = 404000 + iota
	ERROR_FORBID ERROR = 405000 + iota

	ERROR_VERIFY_NEEDBINDPHONE ERROR = 4000001
	ERROR_VERIFY_NEED          ERROR = 4000002
	ERROR_VERIFY_FAILED        ERROR = 4000002
)

/*
* NewError 构造错误
* err 如果err的类型是err或string,将错误信息写入ErrorMessage
* 	   如果err是StackErr,直接返回
* ext ext[0]:错误code  ext[1]:返回给调用端的错误信息
*/
func NewError(err interface{}, ext ...interface{}) *StackErr {
	return newError(err, ext...)
}
func newError(err interface{}, ext ...interface{}) *StackErr {

	var errMessage string
	switch t := err.(type) {
	case *StackErr:
		return t
	case string:
		errMessage = Filter(t)
	case error:
		errMessage = Filter(t.Error())
	default:
		errMessage = Filter(fmt.Sprintf("%v", t))
	}
	stackErr := &StackErr{}
	stackErr.Tag = "DEFAULT"
	stackErr.Code = 0

	stackErr.ErrorMessage = errMessage
	_, file, line, ok := runtime.Caller(2)
	if ok {
		stackErr.Line = line
		components := strings.Split(file, "/")
		stackErr.Filename = components[(len(components) - 1)]
		stackErr.Position = filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	const size = 1 << 12
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)
	stackErr.StackTrace = Filter(string(buf[:n]), " ")

	if len(ext) >= 1 {
		code := ext[0]
		switch c := code.(type) {
		case errorutil.ERROR:
			stackErr.Code = int(c)
		case ERROR:
			stackErr.Code = int(c)
		default:
			stackErr.Code = cast.ToInt(c)
		}
	}
	if len(ext) >= 2 {
		if msg, ok := ext[1].(string); ok {
			stackErr.Info = Filter(msg)
		}
	}
	if len(ext) >= 3 {
		stackErr.Data = ext[2]
	}

	return stackErr
}

func NewErrorWithTag(tag string, err interface{}, ext ...interface{}) *StackErr {
	stackErr := newError(err, ext...)

	isMsgChinese := textutil.IsChineseChar(stackErr.ErrorMessage)
	isInfoChinese := textutil.IsChineseChar(stackErr.Info)

	if isMsgChinese && !isInfoChinese && stackErr.Info != "" {
		stackErr.Info = stackErr.ErrorMessage
	} else if isInfoChinese && !isMsgChinese && stackErr.ErrorMessage != "" {
		stackErr.ErrorMessage = stackErr.Info
	}

	tag = Filter(tag)
	if tag != "" {
		stackErr.Tag = tag
	}
	return stackErr
}

func NewDbError(err error) (e *StackErr) {
	e = newError(err)
	e.Tag = "DB"
	e.Code = MIOT_DB_ERROR
	e.Info = "服务器繁忙"
	return e
}

func NewESError(err error) (e *StackErr) {
	e = newError(err)
	e.Tag = "ES"
	e.Code = MIOT_ES_ERROR
	e.Info = "服务器繁忙"
	return e
}
func NewJsonError(err error) (e *StackErr) {
	e = newError(err)
	e.Tag = "JSON"
	e.Code = MIOT_JSON_ERROR
	e.Info = "服务器繁忙"
	return e
}
func NewRedisError(err error) (e *StackErr) {
	e = newError(err)
	e.Tag = "REDIS"
	e.Code = MIOT_REDIS_ERROR
	e.Info = "服务器繁忙"
	return e
}

func NewThriftError(err error) (e *StackErr) {
	e = newError(err)
	e.Tag = "THRIFT"
	e.Code = MIOT_THRIFT_ERROR
	e.Info = "服务器繁忙"
	return e
}

func (this *StackErr) Msg() string {
	return this.Info
}

func (this *StackErr) Error() string {
	return this.ErrorMessage
}

func (this *StackErr) Stack() string {
	return fmt.Sprintf("(%s:%d)%s\tStack: %s", this.Filename, this.Line, this.ErrorMessage, this.StackTrace)
}

func (this *StackErr) Detail() string {
	return fmt.Sprintf("(%s:%d)%s", this.Filename, this.Line, this.ErrorMessage)
}

func (this *StackErr) Format(tag ...string) (data string) {
	var strs []string
	if len(tag) > 0 {
		this.Tag = tag[0]
	}
	strs = append(strs, this.Tag)
	strs = append(strs, cast.ToString(this.Code))
	strs = append(strs, this.Filename)
	strs = append(strs, cast.ToString(this.Line))
	strs = append(strs, this.ErrorMessage)
	data = strings.Join(strs, "\t")
	return
}

