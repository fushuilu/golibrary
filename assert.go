package golibrary

import (
	"reflect"
	"regexp"
	"strings"
)


func IsBoolean(text string) bool {
	switch strings.ToLower(text) {
	case "true", "yes", "t", "ok", "success":
		return true
	}
	return false
}
func IsStruct(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Struct
}
func IsPointer(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Ptr
}
func IsMap(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Map
}

func IsArray(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Array
}

func IsSlice(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Slice
}

// domain 需要协议，http 或 https，最好以 '/' 结尾，防止冒充
func InDomain(domains []string, src string) bool {
	for i := range domains {
		index := strings.Index(src, domains[i])
		switch index {
		case 4, 5:
			return true
		}
	}
	return false
}

func IsInInt64(slice []int64, find int64) bool {
	return IndexOfInt64(slice, find) != -1
}
func IsInInt(slice []int, find int) bool {
	return IndexOfInt(slice, find) != -1
}

// 字符串查询
func IsInString(texts []string, find string) (isIn bool) {
	for _, v := range texts {
		if v == find {
			return true
		}
	}
	return false
}

var int8DateReg = regexp.MustCompile(`^\d{8}$`)

func IsInt8Date(text string) bool {
	return int8DateReg.MatchString(text)
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
