package datax

import (
	"errors"
	"fmt"

	"github.com/fushuilu/golibrary"
)

// 通常用于模型中常量的转换
type Dict struct {
	Name     string
	DefText  string
	DefValue interface{} // int 必须是能够转为 int
	Data     map[string]interface{}
}

// 将前转提交的字符串，转为指定类型
func (d *Dict) MustGetValueIfNotEmpty(text string, typeTo func(int)) error {
	if text == "" {
		return nil
	}
	return d.MustGetValue(text, typeTo)
}

func (d *Dict) MustGetValue(text string, typeTo func(int)) error {
	if text == "" {
		return errors.New(fmt.Sprintf("必须设置%s", d.Name))
	}
	if e, ok := d.Data[text]; ok {
		typeTo(golibrary.AnyToInt(e))
		return nil
	} else if text != "" {
		return errors.New("匹配错误:" + text)
	} else if d.DefValue == nil {
		return errors.New("参数错误:" + d.Name)
	} else {
		typeTo(golibrary.AnyToInt(d.DefValue))
		return nil
	}
}

// 通常用在校验
func (d *Dict) MustInclude(text string) error {
	if text == "" {
		return errors.New(fmt.Sprintf("必须设置%s", d.Name))
	}
	if _, ok := d.Data[text]; ok {
		return nil
	}
	return errors.New("参数错误:" + d.Name)
}

func (d *Dict) MustIncludeIfNotEmpty(text string) error {
	if text == "" {
		return nil
	}
	return d.MustInclude(text)
}

func (d *Dict) MustMapValueIfNotEmpty(text string, typeTo func(int)) error {
	if text != "" {
		return d.MustMapValue(text, typeTo)
	}
	return nil
}

/**
 * @Description: 将 text 转换为值(如果没有对应值，则报错，注意不使用默认值)
 */
func (d *Dict) MustMapValue(text string, typeTo func(int)) error {
	if text == "" {
		return errors.New(fmt.Sprintf("必须设置 %s", d.Name))
	}
	if e, ok := d.Data[text]; ok {
		typeTo(golibrary.AnyToInt(e))
	} else {
		return errors.New(fmt.Sprintf("参数 %s 值不被允许 %s", d.Name, text))
	}
	return nil
}

func (d *Dict) GetValueIfNotEmpty(text string, typeTo func(i int)) {
	if text != "" {
		if e, ok := d.Data[text]; ok {
			typeTo(golibrary.AnyToInt(e))
		} else if d.DefValue != nil {
			typeTo(golibrary.AnyToInt(d.DefValue))
		}
	}
}

// 直接取值，如果不存在，则使用默认值
func (d *Dict) GetIntValue(text string) int {
	if e, ok := d.Data[text]; ok {
		return golibrary.AnyToInt(e)
	} else if d.DefValue == nil {
		return 0
	}
	return golibrary.AnyToInt(d.DefValue)
}

func (d *Dict) GetValue(text string, typeTo func(int)) {
	typeTo(d.GetIntValue(text))
}

// 将指定常量输出为字符串
func (d *Dict) ToText(t interface{}) string {
	if t == 0 {
		return ""
	}
	i := golibrary.AnyToInt(t)
	if i == 0 {
		return ""
	}
	for k, v := range d.Data {
		if i == golibrary.AnyToInt(v) {
			return k
		}
		if v == t || v == i {
			return k
		}
	}
	return d.DefText
}

func (d *Dict) AppendData(text string, value interface{}, strictValue bool) error {
	if _, ok := d.Data[text]; ok {
		return errors.New(fmt.Sprintf("%s 已经存在", text))
	} else if strictValue {
		if find := d.ToText(value); find != "" {
			return errors.New(fmt.Sprintf("%s: %+v 值已经存在", text, value))
		}
	}
	d.Data[text] = value
	return nil
}

var MapStatusText = Dict{
	Name:     "状态",
	DefText:  StatusActive,
	DefValue: IndexStatusActive,
	Data: map[string]interface{}{
		StatusActive:   IndexStatusActive,
		StatusDisabled: IndexStatusDisabled,
		StatusLock:     IndexStatusLock,
		StatusWarning:  IndexStatusWaring,
		StatusDelete:   IndexStatusDelete,
		StatusCheck:    IndexStatusCheck,
	},
}

// 资源域名类型
type MediaKind int

const (
	_              MediaKind = iota
	MediaKindText            = 1
	MediaKindImage           = 2
	MediaKindVideo           = 3
	MediaKindAudio           = 4
	MediaKindUrl             = 5
	MediaKindMap             = 6
	MediaKindFile            = 7
)

var MapMediaKind = Dict{
	Name:     "媒体类型",
	DefText:  "image",
	DefValue: MediaKindImage,
	Data: map[string]interface{}{
		"text":  MediaKindText,
		"image": MediaKindImage,
		"img":   MediaKindImage,
		"video": MediaKindVideo,
		"audio": MediaKindAudio,
		"url":   MediaKindUrl,
		"map":   MediaKindMap,
		"file":  MediaKindFile,
	},
}
