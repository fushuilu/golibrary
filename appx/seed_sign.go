package appx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx/errorx"
)

type SeedSign interface {
	Create(seed string) string
	MustNotExpired(sign string, activeSeconds int64) (err error)
	Invalid(seed, sign string) error
	GetCreatedAt(sign string) (int64, error)
}

func SeedSignInstance() SeedSign {
	return &seedSignInstance
}

var seedSignInstance = seedSign{}

type seedSign struct {
}

/// 为随机种子生成一个签名： 签名-时间戳
/// days 天数有效期
func (ss *seedSign) Create(seed string) string {
	at := golibrary.AnyToString(time.Now().Unix())
	s, _ := golibrary.MD5(seed + at)
	return strings.Join([]string{s[1:8], at}, "-")
}

// 检查签名是否过期
func (ss *seedSign) MustNotExpired(sign string, activeSeconds int64) (err error) {
	cs := strings.Split(sign, "-")
	if len(cs) != 2 {
		return errors.New("签名格式错误")
	}
	if time.Now().Unix()-activeSeconds > golibrary.AnyToInt64(cs[1]) {
		return errors.New("签名已过期")
	}
	return
}

/// 签名验证
func (ss *seedSign) Invalid(seed, sign string) error {
	cs := strings.Split(sign, "-")
	if len(cs) != 2 {
		return errors.New("签名格式错误")
	}
	signCmp, _ := golibrary.MD5(seed + cs[1])
	if signCmp[1:8] != cs[0] {
		return errorx.New("签名错误", fmt.Sprintf("签名验证：s1: %s; s2:%s", signCmp, cs[0]))

	}
	return nil
}

// 提取签名生成的时间
func (ss *seedSign) GetCreatedAt(sign string) (int64, error) {
	cs := strings.Split(sign, "-")
	if len(cs) != 2 {
		return 0, errors.New("签名格式错误")
	}
	return strconv.ParseInt(cs[1], 10, 0)
}
