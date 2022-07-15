package model

type ImageListReq struct {
	Dir        string  `form:"dir" binding:"required,min=1"`
	HWOperator string  `form:"hwoperator" binding:"oneof=gte lte"`
	HWRatio    float64 `form:"hwratio" binding:"gte=0,lte=1"`
	PageSize   int     `form:"pagesize" binding:"gte=0"`
	PageNum    int     `form:"pagenum" binding:"gte=0"`
}

type DownloadImageReq struct {
	Path string `form:"path" binding:"required,min=1"`
}
