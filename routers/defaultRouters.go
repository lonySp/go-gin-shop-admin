package routers

import (
	"github.com/gin-gonic/gin"
	xyuan "github.com/lonySp/go-gin-shop-admin/controllers/xyuan"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", xyuan.DefaultController{}.Index)
		defaultRouters.GET("/thumbnail1", xyuan.DefaultController{}.Thumbnail1)
		defaultRouters.GET("/thumbnail2", xyuan.DefaultController{}.Thumbnail2)
		defaultRouters.GET("/qrcode1", xyuan.DefaultController{}.Qrcode1)
		defaultRouters.GET("/qrcode2", xyuan.DefaultController{}.Qrcode2)

	}
}
