package golibrary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Int64Unique(t *testing.T) {

	data := []int64{1, 3, 5, 7, 3, 5}

	rst := Int64Unique(data)
	assert.Equal(t, t, len(rst) == 4)
}

func Test_IsPointer(t *testing.T) {
	pd := struct {
		Name string
	}{}

	assert.False(t, IsPointer(pd))
	assert.True(t, IsPointer(&pd))
}
