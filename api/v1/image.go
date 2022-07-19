package v1

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"wallpaper-filter/model"
	"wallpaper-filter/util"

	"github.com/gin-gonic/gin"
)

// 用map模拟cache，避免重复扫文件；无过期时间
var imagePathCache = make(map[string][]string)

func GetImageList(c *gin.Context) {
	// fmt.Print("[DEBUG] GetImageList", c.Request)
	// bind req
	var req model.ImageListReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("[ERROR] %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": util.ErrParseReq,
		})
		return
	}

	// fmt.Println("[DEBUG] req: ", req)

	// load and filter
	imageList := loadImagePaths(req.Dir)
	// fmt.Println("L34-[DEBUG] imageList: ", imageList)

	var err error
	switch req.HWOperator {
	case "gte":
		imageList, err = filterImagePathsByHwRatio(imageList, util.GreaterOrEqualThan, req.HWRatio)
	case "lte":
		imageList, err = filterImagePathsByHwRatio(imageList, util.LessOrEqualThan, req.HWRatio)
	}
	if err != nil {
		fmt.Printf("[ERROR] %v", util.ErrFilterImage)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": util.ErrFilterImage,
		})
		return
	}

	// fmt.Println("L51-[DEBUG] imageList: ", imageList)

	// pager
	imageList, err = getImageListPage(imageList, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": util.ErrParseReq,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"paths": imageList,
	})
}

// getImageListPage 分页逻辑
func getImageListPage(imageList []string, req model.ImageListReq) ([]string, error) {
	// default: return all
	if req.PageNum == 0 && req.PageSize == 0 {
		return imageList, nil
	}
	if req.PageNum == 0 || req.PageSize == 0 {
		log.Printf("[ERROR] one of page_num, page_size is emtpy.")
		return []string{}, nil
	}
	// pager
	var idxStart, idxEnd int
	if req.PageNum*req.PageSize > len(imageList)-1 {
		idxStart = len(imageList) - 1
	} else {
		idxStart = req.PageNum * req.PageSize
	}
	if (req.PageNum+1)*req.PageSize > len(imageList) {
		idxEnd = len(imageList)
	} else {
		idxEnd = req.PageNum * (req.PageSize + 1)
	}
	return imageList[idxStart:idxEnd], nil
}

// loadImagePaths 优先从cache中取图片地址集合，如不存在则遍历
func loadImagePaths(relPath string) []string {
	if _, ok := imagePathCache[relPath]; !ok {
		imagePathCache[relPath] = fetchAllImagePaths(relPath)
	}
	return imagePathCache[relPath]
}

// fetchAllImagePaths 遍历所有子目录找图片，返回相对地址
func fetchAllImagePaths(relDir string) []string {
	baseDir := util.Config().BaseDir
	filePaths := util.ListDirRecur(path.Join(baseDir, relDir))
	// fmt.Println("[DEBUG] filePaths: ", filePaths)
	var imagePaths []string
	for _, p := range filePaths {
		relPath, err := filepath.Rel(baseDir, p)
		if err != nil {
			continue
		}
		switch filepath.Ext(relPath) {
		case ".webp":
			imagePaths = append(imagePaths, relPath)
		case ".jpg":
			imagePaths = append(imagePaths, relPath)
		case ".jpeg":
			imagePaths = append(imagePaths, relPath)
		case ".png":
			imagePaths = append(imagePaths, relPath)
		}
	}
	return imagePaths
}

// filterImagePathsByHwRatio 过滤宽高比满足要求的图片地址
func filterImagePathsByHwRatio(paths []string, hwOperator util.HWOperatorType, hwRatio float64) ([]string, error) {
	var pathFiltered []string
	for _, p := range paths {
		size, err := util.GetImageSize(filepath.Join(util.Config().BaseDir, p))
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		if size.Width == 0 {
			fmt.Println("[ERROR] zero division: size.width")
			continue
		}
		ratio := float64(size.Height) / float64(size.Width)
		switch hwOperator {
		case util.GreaterOrEqualThan:
			if ratio >= hwRatio {
				pathFiltered = append(pathFiltered, p)
			}
		case util.LessOrEqualThan:
			if ratio <= hwRatio {
				pathFiltered = append(pathFiltered, p)
			}
		}
	}
	return pathFiltered, nil
}

// GetImage 响应前端请求返回图片
func GetImage(c *gin.Context) {
	relPath := c.Query("path")
	baseDir := util.Config().BaseDir
	absPath := path.Join(baseDir, relPath)
	fmt.Println("[DEBUG] absPath: ", absPath)
	c.File(absPath)
}

func DownloadImage(c *gin.Context) {
	var req model.DownloadImageReq
	err := c.ShouldBind(req)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": util.ErrParseReq,
		})
		return
	}
	srcPath := filepath.Join(util.Config().BaseDir, req.Path)
	if err = util.CopyFile(srcPath, util.Config().DownloadDir); err != nil {
		fmt.Printf("[ERROR] %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": util.ErrDownloadImage,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "okay",
	})
}
