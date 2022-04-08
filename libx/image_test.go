package cmn

import (
	"github.com/fushuilu/golibrary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageDraw(t *testing.T) {

	// 当前 libx 目录
	path := golibrary.PathGoLibrary() + "libx/"

	// 底图图片
	dst := path + "qrcode_test.png"
	exist, err := golibrary.FileExist(dst)
	assert.Nil(t, err)
	assert.True(t, exist)

	dstImage, _, err := ImageDecode(dst)
	assert.Nil(t, err)

	// 生成的二维码
	qrcode, err := NewQrcode("https://www.baidu.com")
	assert.Nil(t, err)
	qrImage := qrcode.Image(84)

	//  合成
	rgba := ImagesMerge(dstImage, qrImage, 140, 140)
	err = ImagesMergeSave(rgba, path+"qrcode_testrst.jpg")
	assert.Nil(t, err)
}
