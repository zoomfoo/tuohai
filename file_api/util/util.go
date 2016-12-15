package util

import (
	"fmt"
	"image"
	"io"
	"strings"

	"tuohai/internal/console"
)

func ImgDimension(r io.Reader) (width, height int) {
	Config, _, err := image.DecodeConfig(r)
	if err != nil {
		console.StdLog.Error(err)
		return
	}
	width, height = Config.Width, Config.Height
	return
}

func IsImg(filename string) bool {
	names := strings.Split(filename, "/")
	if len(names) == 0 {
		return false
	}
	fmt.Println(names[0] == "image")
	return names[0] == "image"
}
