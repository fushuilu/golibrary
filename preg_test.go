package golibrary

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func TestIsUrl(t *testing.T) {
	data := []struct {
		url string
		rst bool
	}{
		{url: "http://www.golangcode.com", rst: true},
		{url: "http://golangcode.com", rst: true},
		{url: "golangcode.com", rst: false},
		{url: "", rst: false},
	}

	for _, v := range data {
		assert.Equal(t, v.rst, IsUrl(v.url))
	}
}

func TestIsEmail(t *testing.T) {
	data := []struct {
		email string
		rst   bool
	}{
		{email: "123@qq.com", rst: true},
		{email: "12345678901@qq.com", rst: true},
		{email: "abc@163.com.cn", rst: true},
		{email: "abc.com.cn", rst: false},
	}

	for _, v := range data {
		assert.Equal(t, v.rst, IsEmail(v.email), fmt.Sprintf("%s check failed\n", v.email))
	}
}

func TestGetPrefixNumber(t *testing.T) {
	data := []struct {
		Text string
		Num  int
	}{
		{Text: "3月", Num: 3},
		{Text: "1年", Num: 1},
		{Text: "活期", Num: 0},
	}
	for _, v := range data {
		rst := GetPrefixNumber(v.Text)
		assert.Equal(t, v.Num, rst)
	}
}

func TestPhone(t *testing.T) {
	data := []struct {
		Phone string
		Rst   bool
	}{
		{Phone: "1340000", Rst: false},
		{Phone: "13400000000", Rst: true},
		{Phone: "86-13400000000", Rst: true},
		{Phone: "00852-11112222", Rst: false},
	}
	for _, v := range data {
		assert.Equal(t, IsPhone(v.Phone), v.Rst)
	}
}
