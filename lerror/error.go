package lerror

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

/*
gerror.Wrap 会记录太多的无效信息，特别是放在 BindHandler 之中
示例
0). /xxx/mainDemo/api/index.go:31
        [1642043802] datacenter/mainDemo/api.(*Index).Error
        error in api # 即 error.message 记录错误信息，也用于响应
        -->[] # 即 error.detail 通常记录参数
1). /xxx/mainDemo/service/log.go:27
        [1642043802] datacenter/mainDemo/service.(*Log).Error
        error in service
        -->[debug 15]
2). /xxx/mainDemo/action/log.go:23
        [1642043802] datacenter/mainDemo/action.(*Log).Error
        error in action
        -->[]
*/

// 记录调用者
type Error struct {
	preErr  error       // 上一个错误
	time    int64       // 时间
	file    string      // 文件名称
	name    string      // 方法名称
	line    int         // 位置
	detail  interface{} // 更多的细节信息，例如参数
	message string      // 错误信息
	code    int         // 业务错误码
	debug   bool        // 是否为 debug 错误，此参数会传递给下一个 error；使用场景，action/service 中的参数错误之类不需要记录到日志中
}

func (err *Error) Error() string {
	return err.message
}

// 返回业务错误码
func (err *Error) Code() int {
	return err.code
}

func (err *Error) IsDebug() bool {
	return err.debug
}

func New(message string, detail ...interface{}) error {
	return Skip(2, http.StatusInternalServerError, message, false, detail)
}

func NewCode(code int, message string) error {
	return Skip(2, code, message, true, nil)
}

func NewIf(err error, message string, detail ...interface{}) error {
	if err == nil {
		return nil
	}
	return skipWrap(2, err, false, message, detail)
}

func DebugIf(err error, message string, detail ...interface{}) error {
	if err == nil {
		return nil
	}
	return skipWrap(2, err, true, message, detail)
}

func Skip(skip int, code int, message string, debug bool, detail interface{}) error {
	pc, file, line, _ := runtime.Caller(skip)
	return &Error{
		time:    time.Now().Unix(),
		file:    file,
		name:    runtime.FuncForPC(pc).Name(),
		line:    line,
		detail:  detail,
		code:    code,
		message: message,
		debug:   debug,
	}
}

func Wrap(err error, message string, detail ...interface{}) error {
	return skipWrap(2, err, false, message, detail)
}

func skipWrap(skip int, err error, debug bool, message string, detail interface{}) error {
	if err == nil {
		err = errors.New(message)
	} else {
		if e, ok := err.(*Error); ok {
			debug = e.debug
		}
	}
	pc, file, line, _ := runtime.Caller(skip)

	return &Error{
		name:    runtime.FuncForPC(pc).Name(),
		time:    time.Now().Unix(),
		preErr:  err, // 可能是普通错误
		file:    file,
		line:    line,
		message: message,
		detail:  detail,
		debug:   debug,
	}
}

func (err *Error) Stack() string {

	if err == nil {
		return ""
	}

	buffer := bytes.NewBuffer(nil)

	cur := err
	index := 0
	code := 0

	for cur != nil {
		buffer.WriteString(fmt.Sprintf(
			"\n%d). %s:%d\n   \t[%d] %s\n   \t%s\n   \t-->%+v",
			index, cur.file, cur.line,
			cur.time, cur.name,
			cur.message,
			cur.detail,
		))
		if cur.code > 0 { // 错误码
			code = cur.code
		}

		if cur.preErr != nil {
			if e, ok := cur.preErr.(*Error); ok {
				cur = e
			} else { // 其它的错误
				buffer.WriteString(fmt.Sprintf("\n|<- reason:%+v", cur.preErr))
				break
			}
		} else {
			break
		}

		index += 1
	}
	err.code = code
	return buffer.String()
}

func PanicIfError(err error) {
	if err != nil {
		switch err.(type) {
		case *Error:
			fmt.Println("error stack")
			fmt.Println(err.(*Error).Stack())
		}
		panic(err)
	}
}
