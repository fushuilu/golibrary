package cmn

import (
	"errors"
	"fmt"
	"image/png"
	"io"

	"github.com/skip2/go-qrcode"
)
// 根据 url 生成二维码，并渲染到 io 中
func QrcodeFromURL(url string, size int, w io.Writer) error {
	if qrCode, err := NewQrcode(url); err != nil {
		return err
	} else {
		return png.Encode(w, qrCode.Image(size))
	}
}
// 使用示例，查看 image_test.go
func NewQrcode(content string) (*qrcode.QRCode, error) {
	if q, err := qrcode.New(content, qrcode.Medium); err != nil {
		return q, errors.New(fmt.Sprintf("create qrcode failed: %s", err))
	} else {
		return q, nil
	}
}
