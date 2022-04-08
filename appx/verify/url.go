package verify

import (
	"errors"
	"fmt"

	"github.com/fushuilu/golibrary"
)

// 全局域名
var Domains = make([]string, 0)

func URLsInvalid(url ...string) error {
	for _, u := range url {
		if err := URLInvalid(u); err != nil {
			return err
		}
	}
	return nil
}

//
//  URLInvalid
//  @Description: 是否不被允许
//  @param url
//  @return error
//
func URLInvalid(url string) error {
	if url == "" {
		return nil
	}
	host, err := golibrary.HttpURLHost(url)
	if err != nil {
		return err
	}
	for _, v := range Domains {
		if v == host {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("不被允许的地址:%s", url))
}

func IsURLsDisable(url ...string) bool {
	for _, u := range url {
		if u == "" {
			continue
		}
		host, _ := golibrary.HttpURLHost(u)
		if !golibrary.IsInString(Domains, host) {
			return false
		}
	}
	return true
}
