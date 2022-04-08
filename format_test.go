package golibrary

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FormatFloat(t *testing.T) {

	rst := FormatFloat(1/3.0, 6)
	fmt.Println(rst)
	assert.Equal(t, 0.333333, rst)
}

func Test_FormatPrice(t *testing.T) {

	data := []struct {
		Data int
		Rst  string
	}{
		{Data: 100, Rst: "1.00"},
		{Data: 101, Rst: "1.01"},
	}

	for _, v := range data {
		assert.Equal(t, v.Rst, FormatPrice(v.Data))
	}
}
