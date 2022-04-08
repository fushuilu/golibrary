package dbx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type article struct {
	Id    int64
	Title string
}

func (a article) GetId() int64 {
	return a.Id
}

func (a article) GetTitle() string {
	return a.Title
}

func TestNewMapInt64String(t *testing.T) {

	rows := []article{
		{Id: 1, Title: "name1"},
		{Id: 2, Title: "name2"},
		{Id: 3, Title: "name3"},
	}


	dd, ok := NewMapInt64String(rows)
	assert.True(t, ok)
	assert.Equal(t, "name1", dd.Get(1))
	assert.Equal(t, "name3", dd.Get(3))
}
