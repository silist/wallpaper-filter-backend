package v1

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"wallpaper-filter/model"
	"wallpaper-filter/util"

	"github.com/gin-gonic/gin"
)

type ImagePath struct {
	RelPath string `json:"rel_path"`
	AbsPath string `json:"abs_path"`
}

// 用map模拟cache，避免重复扫文件；无过期时间
var ImagePathCache = make(map[string][]ImagePath)

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
	// fmt.Println("L34-[DEBUG] imageList: ", len(imageList))

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
func getImageListPage(imageList []ImagePath, req model.ImageListReq) ([]ImagePath, error) {
	// default: return all
	if req.PageNum == 0 && req.PageSize == 0 {
		return imageList, nil
	}
	// if req.PageNum == 0 || req.PageSize == 0 {
	// 	log.Printf("[ERROR] one of page_num, page_size is emtpy.")
	// 	return []ImagePath{}, nil
	// }

	// pager

	var idxStart = req.PageNum * req.PageSize
	var idxEnd = (req.PageNum + 1) * req.PageSize

	if idxStart > len(imageList) {
		idxStart = len(imageList)
	}
	if idxEnd > len(imageList) {
		idxEnd = len(imageList)
	}

	fmt.Printf("[pager] idxStart: %d idxEnd: %d\n", idxStart, idxEnd)
	return imageList[idxStart:idxEnd], nil
}

// loadImagePaths 优先从cache中取图片地址集合，如不存在则遍历
func loadImagePaths(relPath string) []ImagePath {
	if _, ok := ImagePathCache[relPath]; !ok {
		ImagePathCache[relPath] = fetchAllImagePaths(relPath)
	}
	return ImagePathCache[relPath]
}

// fetchAllImagePaths 遍历所有子目录找图片，返回相对地址
func fetchAllImagePaths(relDir string) []ImagePath {
	baseDir := util.Config().BaseDir
	filePaths := util.ListDirRecur(path.Join(baseDir, relDir))
	// fmt.Println("L109-[DEBUG] filePaths: ", len(filePaths))
	var imagePaths []ImagePath
	for _, p := range filePaths {
		relPath, err := filepath.Rel(baseDir, p)
		if err != nil {
			continue
		}
		switch filepath.Ext(relPath) {
		case ".webp":
			imagePaths = append(imagePaths, ImagePath{RelPath: relPath, AbsPath: p})
		case ".jpg":
			imagePaths = append(imagePaths, ImagePath{RelPath: relPath, AbsPath: p})
		case ".jpeg":
			imagePaths = append(imagePaths, ImagePath{RelPath: relPath, AbsPath: p})
		case ".png":
			imagePaths = append(imagePaths, ImagePath{RelPath: relPath, AbsPath: p})
		}
	}
	return imagePaths
}

// filterImagePathsByHwRatio 过滤宽高比满足要求的图片地址
func filterImagePathsByHwRatio(paths []ImagePath, hwOperator util.HWOperatorType, hwRatio float64) ([]ImagePath, error) {
	// fmt.Printf("[DEBUG] filterImagePathsByHwRatio=%v\n", paths)
	var pathFiltered []ImagePath
	for _, p := range paths {
		size, err := util.GetImageSize(filepath.Join(util.Config().BaseDir, p.RelPath))
		// fmt.Printf("[DEBUG] image_path=%v|size=%v\n", p, size)
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
	// fmt.Println("[DEBUG] absPath: ", absPath)
	c.File(absPath)
}

func DownloadImage(c *gin.Context) {
	var req model.DownloadImageReq
	err := c.ShouldBind(&req)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": util.ErrParseReq,
		})
		return
	}
	// fmt.Println("[DEBUG] req: ", req)
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
