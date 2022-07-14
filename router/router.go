package router

import (
	"github.com/gin-gonic/gin"
	"log"
	v1 "wallpaper-filter/api/v1"
	"wallpaper-filter/util"
)

func InitRouter() {
	r := gin.Default()

	rv1 := r.Group("api/v1")
	{
		rv1.GET("image_list", v1.GetImageList)
	}

	err := r.Run(util.Config().Addr)
	if err != nil {
		log.Fatal(err)
	}
}
