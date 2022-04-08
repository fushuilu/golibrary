package golibrary

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func HasString(data string) (res string, has bool) {
	if data == "" || data == " " {
		return "", false
	}
	res = strings.TrimSpace(data)
	switch res {
	case "null", "undefined", "", "0":
		return "", false
	}
	return res, true
}

func FilterString(data string) string {
	res := strings.TrimSpace(data)
	switch strings.ToLower(res) {
	case "null", "undefined", "", "0":
		return ""
	}
	return res
}

var tagRule = regexp.MustCompile("[、，]")

// 字符串分割，分割符为中文的、， 英文的 ,
func Split(tag string) []string {
	if tag == "" {
		return []string{}
	}
	tag = tagRule.ReplaceAllString(tag, ",")
	return strings.Split(tag, ",")
}

func PadStartZero(data int, len int) string {
	return fmt.Sprintf("%0*d", len, data)
}

func SplitToInt64s(tag string) ([]int64, error) {
	ids := make([]int64, 0)
	for _, v := range Split(tag) {
		if i, err := strconv.ParseInt(v, 10, 64); err != nil {
			return ids, errors.New(fmt.Sprintf("%s 无法转为整型", v))
		} else {
			ids = append(ids, i)
		}
	}
	return ids, nil
}
func SplitToInts(tag string) ([]int, error) {
	ids := make([]int, 0)
	for _, v := range Split(tag) {
		if i, err := strconv.Atoi(v); err != nil {
			return ids, errors.New(fmt.Sprintf("%s 无法转为整型", v))
		} else {
			ids = append(ids, i)
		}
	}
	return ids, nil
}

// 使用星号代替
func ReplaceWithStar(account string) string {
	if IsPhone(account) {
		return string(account[0:6]) + "****"
	} else if IsEmail(account) {
		ss := strings.Split(account, "@")
		ss0Len := len(ss[0])

		if ss0Len > 6 {
			return strings.Join([]string{
				ss[0][0:4],
				ss[1],
			}, "****@")
		} else if ss0Len > 2 {
			return strings.Join([]string{
				ss[0][0:2],
				ss[1],
			}, "****@")
		} else {
			return "****@" + ss[1]
		}
	}

	return account
}

func SubStringEllipsis(text string, maxLen int) string {
	textRune := []rune(text)
	l := len(textRune)
	if l <= maxLen {
		return text
	}
	return string(textRune[0:maxLen]) + "..."
}

func StringLen(text string) int {
	textRune := []rune(text)
	return len(textRune)
}
