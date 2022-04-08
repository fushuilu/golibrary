package datax

import (
	"errors"
	"github.com/fushuilu/golibrary"
)

// 通常用于修改属性
type RecordAttr struct {
	Name  string      `json:"name"` // 修改的键
	Value interface{} `json:"value"`
}

func (p *RecordAttr) StringValue() string {
	return golibrary.AnyToString(p.Value)
}

func (p *RecordAttr) BoolValue() bool {
	return golibrary.AnyToBool(p.Value)
}

func (p *RecordAttr) StringsValue() []string {
	return golibrary.AnyToStrings(p.Value)
}

func (p *RecordAttr) IntValue() int {
	return golibrary.AnyToInt(p.Value)
}

type RecordAttrs struct {
	Id    int64        `json:"id"`
	Attrs []RecordAttr `json:"attrs"`
}

func (pd *RecordAttrs) Invalid() error {
	if pd.Id < 1 {
		return errors.New("必须指定 ID")
	}
	if len(pd.Attrs) == 0 {
		return errors.New("待修改属性不能为空")
	}
	return nil
}
