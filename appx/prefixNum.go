package appx

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/fushuilu/golibrary"
)

var prefixReg = regexp.MustCompile("^\\d+\\w+$")
var preNumberReg = regexp.MustCompile("^(\\d+)")

// totalLen 总长度
func CreatePrefixNum(num int64, totalLen int) string {
	ids := strconv.FormatInt(num, 10)
	return ids + strings.ToLower(golibrary.RandLetters(totalLen-len(ids))) // 全部为小写（因为小程序 screen 参数只能使用小写)
}

// 返回由 CreatePrefixNum 生成的字符串的前辍数字
func ExplodePrefixNum(numStr string) int64 {
	if numStr == "" {
		return 0
	}
	if prefixReg.MatchString(numStr) {
		id, _ := strconv.ParseInt(preNumberReg.FindString(numStr), 10, 64)
		return id
	} else {
		return 0
	}
}

func MustExplodePrefixNum(numStr string) (int64, error) {
	id := ExplodePrefixNum(numStr)
	if id < 1 {
		return 0, errors.New("参数解析错误")
	}
	return id, nil
}
