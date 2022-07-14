package util

import (
	"github.com/chai2010/webp"
	"io/ioutil"
	"log"
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
	width, height, _, err := webp.GetInfo(data)
	if err != nil {
		log.Println(err)
		return ImageSize{}, err
	}
	return ImageSize{int32(height), int32(width)}, err
}
