package golibrary

import "runtime"

// skip = 1 XXX.Caller
// skip = 2 XXX.XXX.Caller
func Caller(skip int) (file string, line int, funcName string) {
	pc, file, line, _ := runtime.Caller(skip)
	return file, line, runtime.FuncForPC(pc).Name()
}
