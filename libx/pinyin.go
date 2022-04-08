package cmn

import (
	"github.com/mozillazg/go-pinyin"
)

// 汉字首字母，中国人 => z
func FirstPinYin(name string) string {
	if name == "" {
		return ""
	}

	var a = pinyin.NewArgs()
	pys := pinyin.Pinyin(name, a)
	if len(pys) > 0 && len(pys[0]) > 0 {
		first := pys[0][0]
		return string(first[0])
	}
	return ""
}

// 中国人 => zgr
func PinYin(word string) string {
	pp := pinyin.LazyConvert(word, nil)
	cs := make([]uint8, len(pp))
	for i := range pp {
		cs[i] = pp[i][0]
	}
	return string(cs)
}
