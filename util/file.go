package util

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

// 递归列出文件夹中所有文件
func ListDirRecur(dirPath string) []string {
	files, err := ioutil.ReadDir(dirPath)
	var filePaths []string
	if err != nil {
		log.Fatal(err)
	}
	// 遍历输入目录
	for _, f := range files {
		fullPath := filepath.Join(dirPath, f.Name())
		filePaths = append(filePaths, fullPath)
	}
	return filePaths
}
