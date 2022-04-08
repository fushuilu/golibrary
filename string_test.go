package golibrary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillZero(t *testing.T) {
	type args struct {
		data int
		len  int
		rst  string
	}

	datas := []args{
		{data: 15, len: 3, rst: "015"},
		{data: 15, len: 1, rst: "15"},
		{data: 15, len: 2, rst: "15"},
		{data: 15, len: 4, rst: "0015"},
	}

	for _, v := range datas {
		rst := PadStartZero(v.data, v.len)
		assert.Equal(t, v.rst, rst, "want:"+v.rst)
	}

}

func TestReplaceWithStar(t *testing.T) {
	type args struct {
		data string
		rst  string
	}

	data := []args{
		{data: "12345678@qq.com", rst: "1234****@qq.com"},
		{data: "12345@qq.com", rst: "12****@qq.com"},
		{data: "12@qq.com", rst: "****@qq.com"},
		{data: "13012345678", rst: "130123****"},
	}
	for _, v := range data {
		rst := ReplaceWithStar(v.data)
		assert.Equal(t, v.rst, rst)
	}
}


func TestSubStringEllipsis(t *testing.T) {
	text := SubStringEllipsis("在golang中可以通过切片截取一个数组或字符串", 10)
	assert.Equal(t, "在golang中可以...", text)

	text = SubStringEllipsis("hello啊", 10)
	assert.Equal(t, "hello啊", text)
}
