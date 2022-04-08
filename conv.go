package golibrary

import (
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
	"strings"
)

func Int64SToString(ids []int64) string {
	s := make([]string, len(ids))
	for i := range ids {
		s[i] = strconv.FormatInt(ids[i], 10)
	}
	return strings.Join(s, ",")
}

func AnyToFloat64(any interface{}) float64 {
	return gconv.Float64(any)
}

func AnyToInt64(any interface{}) int64 {
	return gconv.Int64(any)
}

func AnyToInt(any interface{}) int {
	return gconv.Int(any)
}

func AnyToBool(any interface{}) bool {
	return gconv.Bool(any)
}

func AnyToString(any interface{}) string {
	return gconv.String(any)
}

func AnyToStrings(any interface{}) []string {
	return gconv.Strings(any)
}

func AnyToBytes(any interface{}) []byte {
	return gconv.Bytes(any)
}
