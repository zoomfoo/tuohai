package util

import (
	"bytes"
	"github.com/nfnt/resize"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func Compress(buf bytes.Buffer, file_out string, width uint, quality int, image_type string) {
	switch image_type {
	case ".png":
		pngCompress(buf, file_out, width, quality)
		break
	case ".jpg":
		jpegCompress(buf, file_out, width, quality)
		break
	case ".gif":
		gifCompress(buf, file_out, width, quality)
		break
	}
}

func pngCompress(buf bytes.Buffer, file_out string, width uint, quality int) {
	origin, err := png.Decode(&buf)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
	canvas := resize.Thumbnail(width, width, origin, resize.Lanczos3)
	if file, err := os.Create(file_out); err == nil {
		png.Encode(file, canvas)
	} else {
		log.Println("ERROR: ", err)
		return
	}
}

func jpegCompress(buf bytes.Buffer, file_out string, width uint, quality int) {
	origin, err := jpeg.Decode(&buf)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
	canvas := resize.Thumbnail(width, width, origin, resize.Lanczos3)
	if file, err := os.Create(file_out); err == nil {
		jpeg.Encode(file, canvas, &jpeg.Options{quality})
	} else {
		log.Println("ERROR: ", err)
		return
	}
}

func gifCompress(buf bytes.Buffer, file_out string, width uint, quality int) {
	origin, err := gif.Decode(&buf)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
	canvas := resize.Thumbnail(width, width, origin, resize.Lanczos3)
	if file, err := os.Create(file_out); err == nil {
		gif.Encode(file, canvas, &gif.Options{})
	} else {
		log.Println("ERROR: ", err)
		return
	}
}
