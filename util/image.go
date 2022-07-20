package util

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
)

type ImageSize struct {
	Height int32
	Width  int32
}

func GetImageSize(path string) (ImageSize, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return ImageSize{}, err
	}

	var width, height int

	switch filepath.Ext(path) {
	case ".webp":
		width, height, _, err = webp.GetInfo(data)
	default:
		reader, _ := os.Open(path)
		defer reader.Close()
		var im image.Config
		im, _, err = image.DecodeConfig(reader)
		width, height = im.Width, im.Height
	}

	if err != nil {
		log.Println(err)
		return ImageSize{}, err
	}
	return ImageSize{int32(height), int32(width)}, err
}
