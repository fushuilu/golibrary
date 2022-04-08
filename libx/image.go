package cmn

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/gabriel-vasile/mimetype"
)

type ImageType int

const (
	_ ImageType = iota
	ImageTypeJpg
	ImageTypePng
	ImageTypeGif
)

// https://yourbasic.org/golang/create-image/
// 创建一张空白的图片
func NewImage(width, height int) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	return image.NewRGBA(image.Rectangle{upLeft, lowRight})
}

// 获取图片的类型
func GetImageType(path string) (it ImageType, err error) {
	dstImage, e := os.Open(path)
	if e != nil {
		err = errors.New(fmt.Sprintf("failed to open image file: %s", e))
		return
	}
	defer dstImage.Close()

	mime, e := Mimetype(path)
	if e != nil {
		err = errors.New(fmt.Sprintf("failed to check image file mimetype: %s", e))
		return
	}

	switch mime.String() {
	case "image/jpeg":
		return ImageTypeJpg, nil
	case "image/png":
		return ImageTypePng, nil
	case "image/gif":
		return ImageTypeGif, nil
	default:
		err = errors.New(fmt.Sprintf("unsupport image mime type: %s", mime.String()))
		return
	}
}

// https://github.com/gabriel-vasile/mimetype/blob/master/mimetype_test.go
func Mimetype(path string) (*mimetype.MIME, error) {
	return mimetype.DetectFile(path)
}

// 将指定路径文件转为图片
func ImageDecode(path string) (img image.Image, it ImageType, err error) {
	dstImage, e := os.Open(path)
	if e != nil {
		err = errors.New(fmt.Sprintf("failed to open image file: %s", e))
		return
	}
	defer dstImage.Close()

	mime, e := Mimetype(path)
	if e != nil {
		err = errors.New(fmt.Sprintf("failed to check image file mimetype: %s", e))
		return
	}

	switch mime.String() {
	case "image/jpeg":
		img, e = jpeg.Decode(dstImage)
		it = ImageTypeJpg
	case "image/png":
		img, e = png.Decode(dstImage)
		it = ImageTypePng
	case "image/gif":
		img, e = gif.Decode(dstImage)
		it = ImageTypeGif
	default:
		err = errors.New(fmt.Sprintf("unsupport image mime type: %s", mime.String()))
		return
	}
	if e != nil {
		err = errors.New(fmt.Sprintf("failed to decode image: %s", e))
		return
	}
	return
}

/*
dstPath 底图图片地址
srcPath 水印图片地址
top, left 绘制起点
savePath 合成图片保存地址
*/
func ImageDraw(dstPath, srcPath string, top, left int, savePath string) (err error) {
	first, _, e := ImageDecode(dstPath)
	if e != nil {
		return e
	}

	second, _, e := ImageDecode(srcPath)
	if e != nil {
		return e
	}

	rgba := ImagesMerge(first, second, top, left)
	return ImagesMergeSave(rgba, savePath)
}

// left = x, top = y
func ImagesMerge(dstImage, srcImage image.Image, left, top int) *image.RGBA {
	offset := image.Pt(left, top)

	b := dstImage.Bounds()

	rst := image.NewRGBA(b)
	draw.Draw(rst, b, dstImage, image.ZP, draw.Src)
	draw.Draw(rst, srcImage.Bounds().Add(offset), srcImage, image.ZP, draw.Over)

	return rst
}

func ImagesMergeSave(rgba *image.RGBA, savePath string) error {
	file, e := os.Create(savePath)
	if e != nil {
		return errors.New(fmt.Sprintf("failed to create save image: %s", e))
	}

	if e = jpeg.Encode(file, rgba, &jpeg.Options{jpeg.DefaultQuality}); e != nil {
		return errors.New(fmt.Sprintf("failed to decode save image: %s", e))
	}

	defer file.Close()
	return nil
}

func ImageOutput(w io.Writer, rgba image.Image, it ImageType) error {
	if it == ImageTypeJpg {
		return jpeg.Encode(w, rgba, &jpeg.Options{jpeg.DefaultQuality})
	} else if it == ImageTypePng {
		return png.Encode(w, rgba)
	}
	return errors.New("只支持输入 jpg/png 图片")
}

func ImageScale(img image.Image, width int) (image.Image, error) {
	bound := img.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	// 缩略图大小
	dst := image.NewRGBA(image.Rect(0, 0, width, width*dy/dx))
	if err := graphics.Scale(dst, img); err != nil {
		return nil, err
	}
	return dst, nil
}
