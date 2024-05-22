package util

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// ListDirRecur 递归列出文件夹中所有文件
func ListDirRecur(dirPath string) []string {
	// fmt.Println("[DEBUG] ListDirRecur path=", dirPath)
	files, err := os.ReadDir(dirPath)
	var filePaths []string
	if err != nil {
		log.Printf("[ERROR] failed to iterate %s.", dirPath)
	}
	// 遍历输入目录
	for _, f := range files {
		fullPath := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			innerPaths := ListDirRecur(fullPath)
			filePaths = append(filePaths, innerPaths...)
		} else {
			filePaths = append(filePaths, fullPath)
		}
	}
	return filePaths
}

func CopyFile(srcPath string, dstDir string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Printf("[ERROR] failed to open file %s.", srcPath)
		return err
	}
	defer src.Close()

	dstPath := filepath.Join(dstDir, filepath.Base(srcPath))
	// 如果文件已存在则直接返回
	if _, err := os.Stat(dstPath); err == nil {
		return nil
	}
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("[ERROR] failed to create file %s.", dstDir)
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Printf("[ERROR] failed to copy file %s.", srcPath)
		return err
	}
	return nil
}
