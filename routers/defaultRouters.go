package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/controllers/xyuan"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", xyuan.DefaultController{}.Index)
		defaultRouters.GET("/news", xyuan.DefaultController{}.News)

	}
}
