package golibrary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockInt int

func TestAnyToInt64(t *testing.T) {
	rst := int64(15)
	datas := []interface{}{
		int8(15), int16(15), int32(15), int64(15),
		uint(15), uint8(15), uint16(15), uint32(15), uint32(15),
		float32(15), float64(15),
	}
	for _, v := range datas {
		assert.Equal(t, rst, AnyToInt64(v))
	}

	assert.Equal(t, int64(1), AnyToInt64(true))
	assert.Equal(t, rst, AnyToInt64("15"))
	assert.Equal(t, 15, AnyToInt(MockInt(15)))
}
