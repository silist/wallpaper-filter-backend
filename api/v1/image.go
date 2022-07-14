package v1

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetImageList(c *gin.Context) {
	filePaths, err := filepath.Glob("**/*.[jpg|png|jpeg|webp]")
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
