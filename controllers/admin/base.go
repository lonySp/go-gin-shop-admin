package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

// Success
// @Description 重定向成功页面
// @Author xYuan 2024-04-18 11:13:35
// @Param c
// @Param message
// @Param redirectUrl
func (con BaseController) Success(c *gin.Context, message string, redirectUrl string) {
	c.HTML(http.StatusOK, "admin/public/success.html", gin.H{
		"message":     message,
		"redirectUrl": redirectUrl,
	})
}

// Error
// @Description 重定向失败页面
// @Author xYuan 2024-04-18 11:14:00
// @Param c
// @Param message
// @Param redirectUrl
func (con BaseController) Error(c *gin.Context, message string, redirectUrl string) {
	c.HTML(http.StatusOK, "admin/public/error.html", gin.H{
		"message":     message,
		"redirectUrl": redirectUrl,
	})
}
