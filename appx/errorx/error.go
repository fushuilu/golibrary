package errorx

import (
	"github.com/fushuilu/golibrary/lerror"
)

// 虽然这个 error 不重要，但我也想知道它出现的错误
// 当你想直接返回 errors.New 时，你可能需要这个方法
func New(msg string, v ...interface{}) error {
	return lerror.Skip(2, 0, msg, true, v)
}
