package router

import (
	"log"
	v1 "wallpaper-filter/api/v1"
	"wallpaper-filter/util"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	rv1 := r.Group("api/v1")
	{
		rv1.GET("image_list", v1.GetImageList)
		rv1.GET("image", v1.GetImage)
	}

	err := r.Run(util.Config().Addr)
	if err != nil {
		log.Fatal(err)
	}
}
