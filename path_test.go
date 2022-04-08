package golibrary

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestFileExist(t *testing.T) {
	path := PathGoLibrary()
	fmt.Println("file exists:", path)

	exist, err := FileExist(path)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = FileExist(path + "abc.txt")
	assert.Nil(t, err)
	assert.False(t, exist)
}

func TestPathAppendName(t *testing.T) {
	data := []struct {
		Dir  string
		Name string
		Rst  string
	}{
		{Dir: "a/", Name: "/b", Rst: "a/b"},
		{Dir: "c/", Name: "b", Rst: "c/b"},
		{Dir: "d", Name: "/b", Rst: "d/b"},
		{Dir: "e", Name: "b", Rst: "e/b"},
		{Dir: "f/b", Name: "/", Rst: "f/b/"},
	}

	for _, v := range data {
		rst := FilePath(v.Dir, v.Name)
		assert.Equal(t, v.Rst, rst)
	}
}

func TestFile(t *testing.T) {
	testPath := "/a/b/c/d.png"
	dir := filepath.Dir(testPath)
	assert.Equal(t, "/a/b/c", dir)
}
