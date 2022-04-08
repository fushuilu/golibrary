package golibrary

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type ClientHttp interface {
	Post(url string, mapData interface{}, resp interface{}) error
	PostByte(url string, data []byte, resp interface{}) error
	Get(url string, mapData interface{}, resp interface{}) error
}

// 提示 ~~~~~~~~~~~~~~~~
// 查看 http_test.go 中 对 net/url 包的使用
/*
q, err := url.Parse(easy.Url)
host, _, err := net.SplitHostPort(q.Host)
*/

//
//  HttpURLHost
//  @Description: 获取 URL 地址中的 host，不包含端口
//  @param href 地址
//  @return string
//
func HttpURLHost(href string) (string, error) {
	if p, err := url.Parse(href); err != nil {
		return "", err
	} else {
		return p.Hostname(), nil
	}
}

//
//  HttpURLHostPort
//  @Description: 获取 URL 地址中的 host:port (包含端口)
//  @param src
//  @return string host or host:port
//  @return error
//
func HttpURLHostPort(src string) (string, error) {
	uri, err := url.Parse(src)
	if err != nil {
		return "", errors.New(fmt.Sprintf("not allow domain invalid:%s", src))
	}
	//host, port, err := net.SplitHostPort(uri.Host)
	return uri.Host, nil
}

//
//  HttpBuildQuery
//  @Description: 合成查询参数，并把键按从小到大排序
//  @param v
//  @return string
//
func HttpBuildQuery(v url.Values) string {
	keys := make([]string, 0)
	for i := range v {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	bb := bytes.Buffer{}
	for _, k := range keys {
		bb.WriteString(k)
		bb.WriteString("=")
		bb.WriteString(v.Get(k))
		bb.WriteString("&")
	}
	return strings.TrimRight(bb.String(), "&")
}

//
//  HttpURLSecure
//  @Description: 将 http://url 转为 https://url，如果含有 localhost，则不转化
//  @param domain
//  @return string
//
func HttpURLSecure(domain string) string {
	if len(domain) > 5 && strings.Index(domain, "localhost") == -1 {
		return strings.Replace(domain, "http://", "https://", 1)
	}
	return domain
}

//
//  HttpURLContact
//  @Description: 将 domain 和 uri 进行连接
//  @param domain
//  @param uri
//  @return string
//
func HttpURLContact(domain, uri string) string {
	if uri == "" {
		return domain
	}
	domainEnd := strings.HasSuffix(domain, "/")
	uriStart := strings.HasPrefix(uri, "/")

	if domainEnd {
		if uriStart {
			return domain + uri[1:]
		} else {
			return domain + uri
		}
	} else {
		if uriStart {
			return domain + uri
		} else {
			return domain + "/" + uri
		}
	}
}

func HttpQueryContact(uri string, query string) string {
	if strings.Contains(uri, "?") {
		return uri + "&" + query
	} else {
		return uri + "?" + query
	}
}

//
//  HttpURLQueryGetWith
//  @Description: 从 url 读取查询参数的值，只想获取单独某个值时可用
//  @param uri string url 地址
//  @param key string 想要查询的 query 的名称
//  @return string 值
//  @return error
//
func HttpURLQueryGetWith(uri string, key string) (string, error) {
	if uri == "" {
		return "", errors.New("getQuery:url 地址为空")
	}
	if key == "" {
		return "", errors.New("getQuery:key 为空")
	}
	if vs, err := url.Parse(uri); err != nil {
		return "", err
	} else {
		return vs.Query().Get(key), nil
	}
}

//
//  HttpFileName
//  @Description: 从 http 下载地址中提取出文件名
//  @param httpUrl string 下载地址
//  @return string 文件名
//  @return error
//
func HttpFileName(httpUrl string) (string, error) {
	// build filename form url
	fileURL, e := url.Parse(httpUrl)
	if e != nil {
		return "", errors.New(fmt.Sprintf("parse http url failed: %s", e))
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	return segments[len(segments)-1], nil
}

// https://golangdocs.com/golang-download-files
//  HttpDownload
//  @Description:  文件下载
//  @param httpUrl 文件访问路径
//  @param filePath 保存地址，需要带上文件名
//  @return written
//  @return err
//
func HttpDownload(httpUrl string, filePath string) (written int64, err error) {

	// create a blank file，你可能需要提前创建目录
	file, e := os.Create(filePath)
	if e != nil {
		err = errors.New(fmt.Sprintf("http download: create file failed: %s", e))
		return
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// put content on file
	resp, e := client.Get(httpUrl)
	if e != nil {
		err = errors.New(fmt.Sprintf("http download: get from url failed: %s", e))
		return
	}

	defer resp.Body.Close()

	size, e := io.Copy(file, resp.Body)
	if e != nil {
		err = errors.New(fmt.Sprintf("http download: copy file content failed: %s", e))
		return
	}
	return size, nil
}
