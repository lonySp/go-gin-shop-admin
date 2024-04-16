package admin

import (
	"fmt"
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login/login.html", gin.H{})

}
func (con LoginController) DoLogin(c *gin.Context) {
	captchaId := c.PostForm("captchaId")
	verifyValue := c.PostForm("verifyValue")

	if models.VerifyCaptcha(captchaId, verifyValue) {
		// c.String(http.StatusOK, "验证码验证成功")
		con.success(c, "验证成功", "/admin")
	} else {
		// c.String(http.StatusOK, "验证码验证失败")
		con.error(c, "验证失败", "/admin/login")
	}
}
func (con LoginController) Captcha(c *gin.Context) {
	captcha, captchaImg, _, err := models.MackCaptcha()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"captchaId": captcha, "captchaImg": captchaImg, "err": err})
}
