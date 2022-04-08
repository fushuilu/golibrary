package golibrary

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strings"
	"testing"
)

// 内置的标准包
func Test_NetURL(t *testing.T) {
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]

	href := "http://test.com:5432/bower_components/console.html?name=test#frame"

	// url.Parse() 将 绝对地址或相对对址 的字段串解析为 URL 结构体
	u, err := url.Parse(href)
	assert.Nil(t, err)
	assert.Equal(t, "http", u.Scheme)
	assert.Equal(t, "", u.Opaque)
	assert.Equal(t, "", u.User.String())     // 认证信息
	assert.Equal(t, "test.com:5432", u.Host) // 主机名+端口
	assert.Equal(t, "/bower_components/console.html", u.Path)
	assert.Equal(t, "", u.RawPath)
	assert.Equal(t, false, u.ForceQuery)
	assert.Equal(t, "name=test", u.RawQuery)
	assert.Equal(t, "frame", u.Fragment) // 分段字段
	assert.Equal(t, "", u.RawFragment)   // see EscapedFragment method

	// opaque 透明类型 URL， 即 scheme 后没有 // 的 URL
	rawurl := "mailto:abc@qq.com"
	u, err = url.Parse(rawurl)
	assert.Nil(t, err)
	assert.Equal(t, "mailto", u.Scheme)
	assert.Equal(t, "abc@qq.com", u.Opaque)

	// url.ParseRequestURI() 忽略 fragment

	u, err = url.ParseRequestURI(href)
	assert.Nil(t, err)
	assert.Equal(t, "http", u.Scheme)
	assert.Equal(t, "", u.Opaque)
	assert.Equal(t, "", u.User.String())
	assert.Equal(t, "test.com:5432", u.Host)
	assert.Equal(t, "/bower_components/console.html", u.Path)
	assert.Equal(t, "", u.RawPath)
	assert.Equal(t, false, u.ForceQuery)
	assert.Equal(t, "name=test#frame", u.RawQuery) // 1. 不解析 frame
	assert.Equal(t, "", u.Fragment)                // 2. 空
	assert.Equal(t, "", u.RawFragment)

	// url.PathEscape 转义，以便安全地放置到 URL 路径中
	// / => %2F, ? => %3F, # => %23, 空格 => %20
	// : 和 & 不转义编码
	rawurl = "http://www.baidu.com/search?name=管理 员#header" // 含特殊符号
	escape := url.PathEscape(rawurl)
	// http:%2F%2Fwww.baidu.com%2Fsearch%3Fname=%E7%AE%A1%E7%90%86%20%E5%91%98%23header
	assert.False(t, strings.Contains(escape, "管"))
	assert.False(t, strings.Contains(escape, "#"))

	// : 和 & 都参与转义编码
	// 空格 => +, : => %3A, & => %26
	queryEscape := url.QueryEscape(rawurl)
	// http%3A%2F%2Fwww.baidu.com%2Fsearch%3Fname%3D%E7%AE%A1%E7%90%86+%E5%91%98%23header
	fmt.Println(queryEscape)

	// url.ParseQuery 将查询字符串解析成  url.Values 字典实例
	u, err = url.Parse(href)
	assert.Nil(t, err)
	v, err := url.ParseQuery(u.RawQuery) // 相当于 u.Query()
	assert.Nil(t, err)
	assert.Equal(t, "test", v.Get("name"))

	v.Add("age", "5")
	assert.Equal(t, "age=5&name=test", v.Encode())
}

func Test_HttpURLContact(t *testing.T) {
	data := []struct {
		Domain string
		Uri    string
		Rst    string
	}{
		{Domain: "http://d.com", Uri: "a/b", Rst: "https://d.com/a/b"},
		{Domain: "https://d.com", Uri: "a/b", Rst: "https://d.com/a/b"},
		{Domain: "http://d.com/", Uri: "a/b", Rst: "https://d.com/a/b"},
		{Domain: "https://d.com/", Uri: "a/b", Rst: "https://d.com/a/b"},
		{Domain: "https://d.com/", Uri: "/a/b", Rst: "https://d.com/a/b"},
	}

	for _, v := range data {
		rst := HttpURLContact(v.Domain, v.Uri)
		assert.Equal(t, v.Rst, rst)
	}
}

func Test_HttpDownload(t *testing.T) {
	path := PathGoLibrary()
	fmt.Println("save file to:", path)
	written, err := HttpDownload("https://www.baidu.com/img/flexible/logo/pc/result@2.png", path+"result@2.png")
	assert.Nil(t, err)
	assert.True(t, written > 1)

	exist, err := FileExist(path + "result@2.png")
	assert.Nil(t, err)
	assert.True(t, exist)
}

func Test_HttpURLRemoveScheme(t *testing.T) {

	data := []struct {
		url string
		rst string
	}{
		{url: "https://demo.com/", rst: "demo.com"},
		{url: "http://demo.com", rst: "demo.com"},
		{url: "http://demo.com/a/", rst: "demo.com"},
	}

	for _, v := range data {
		host, err := HttpURLHost(v.url)
		assert.Nil(t, err)
		assert.Equal(t, v.rst, host)
	}
}
