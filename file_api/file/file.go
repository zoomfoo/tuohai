package file

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func ChangeImgeSize(r io.Reader, x, y, w, h int) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	img, format, err := image.Decode(r)
	fmt.Println("1111:", format)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var subImg image.Image
	switch img.(type) {
	case *image.YCbCr:
		subImg = img.(*image.YCbCr).SubImage(image.Rect(x, y, w, h))
	case *image.CMYK:
		subImg = img.(*image.CMYK).SubImage(image.Rect(x, y, w, h))
	case *image.RGBA:
		subImg = img.(*image.RGBA).SubImage(image.Rect(x, y, w, h))
	case *image.NRGBA:
		subImg = img.(*image.NRGBA).SubImage(image.Rect(x, y, w, h))
	default:
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		return nil, fmt.Errorf("不正确的色域")
	}

	if err := encode(buf, subImg, format); err != nil {
		return nil, err
	}
	return buf, nil
}

func encode(w io.Writer, img image.Image, format string) error {
	switch format {
	case "jpeg":
		return jpeg.Encode(w, img, nil)
	case "png":
		return png.Encode(w, img)
	default:
		return fmt.Errorf("未知图片")
	}
}
