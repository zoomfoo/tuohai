package util

import (
	"bytes"
	"fmt"
	"image"
	"path/filepath"
	"strings"

	"tuohai/internal/console"
)

func ImgDimension(b []byte) (width, height int) {
	Config, _, err := image.DecodeConfig(bytes.NewReader(b))
	if err != nil {
		console.StdLog.Error(err)
		return
	}
	fmt.Println("--------------------------")
	fmt.Println("width:", Config.Width, "height:", Config.Height)
	fmt.Println("--------------------------")
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

func GetExt(name string) string {
	return filepath.Ext(name)[1:]
}
