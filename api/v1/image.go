package v1

import (
	"log"
	"path"
	"path/filepath"
	"strconv"
	"wallpaper-filter/util"

	"github.com/gin-gonic/gin"
)

// 用map模拟cache，避免重复扫文件；无过期时间
var imagePathCache = make(map[string][]string)

// HwOperatorType 宽高比较运算符
type HwOperatorType int32

const (
	GreaterOrEqualThan HwOperatorType = 0 // 大于等于
	LessOrEqualThan    HwOperatorType = 1 // 小于等于
)

func GetImageList(c *gin.Context) {
	var imageList []string
	var err error
	relPath := c.Query("dir")
	// load
	if len(relPath) > 0 {
		imageList = loadImagePaths(relPath)
	}
	// filter
	var hwratio float64
	hwratio, err = strconv.ParseFloat(c.Query("hwratio"), 64)
	if err != nil {
		log.Fatal(err)
	}
	switch c.DefaultQuery("hwoperator", "") {
	case "ge":
		imageList, err = filterImagePathsByHwRatio(imageList, GreaterOrEqualThan, hwratio)
		break
	case "le":
		imageList, err = filterImagePathsByHwRatio(imageList, GreaterOrEqualThan, hwratio)
		break
	default:
		break
	}
	if err != nil {
		log.Fatal(err)
	}
	// pager
	pageNum := c.DefaultQuery("pagenum", "")
	pageSize := c.DefaultQuery("pagesize", "")
	var pageNumInt int
	var pageSizeInt int
	if len(pageNum) > 0 && len(pageSize) > 0 {
		pageNumInt, err = strconv.Atoi(pageNum)
		if err != nil {
			log.Fatal(err)
		}
		pageSizeInt, err = strconv.Atoi(pageSize)
		if err != nil {
			log.Fatal(err)
		}
		var idxStart int
		if pageNumInt*pageSizeInt > len(imageList)-1 {
			idxStart = len(imageList) - 1
		} else {
			idxStart = pageNumInt * pageSizeInt
		}
		var idxEnd int
		if pageNumInt*(pageSizeInt+1) > len(imageList) {
			idxEnd = len(imageList)
		} else {
			idxEnd = pageNumInt * (pageSizeInt + 1)
		}
		imageList = imageList[idxStart:idxEnd]
	}
	c.JSON(200, gin.H{
		"paths": imageList,
	})
}

// 优先从cache中取图片地址集合，如不存在则遍历
func loadImagePaths(relPath string) []string {
	if _, ok := imagePathCache[relPath]; !ok {
		imagePathCache[relPath] = fetchAllImagePaths(relPath)
	}
	return imagePathCache[relPath]
}

func fetchAllImagePaths(relDir string) []string {
	baseDir := util.Config().BaseDir
	filePaths := util.ListDirRecur(path.Join(baseDir, relDir))
	var imagePaths []string
	for _, p := range filePaths {
		switch filepath.Ext(p) {
		case ".webp":
		case ".jpg":
		case ".jpeg":
		case ".png":
			imagePaths = append(imagePaths, p)
			break
		default:
			break
		}
	}
	return imagePaths
}

func filterImagePathsByHwRatio(paths []string, hwOperator HwOperatorType, hwRatio float64) ([]string, error) {
	var pathFiltered []string
	for _, p := range paths {
		size, err := util.GetImageSize(p)
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}
		if size.Width == 0 {
			log.Println("[ERROR] zero division: size.width")
			continue
		}
		ratio := float64(size.Height) / float64(size.Width)
		switch hwOperator {
		case GreaterOrEqualThan:
			if ratio >= hwRatio {
				pathFiltered = append(pathFiltered, p)
			}
			break
		case LessOrEqualThan:
			if ratio <= hwRatio {
				pathFiltered = append(pathFiltered, p)
			}
			break
		}
	}
	return pathFiltered, nil
}
