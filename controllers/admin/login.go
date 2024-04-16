package admin

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/goccy/go-json"
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
	// 获取验证码隐式ID
	captchaId := c.PostForm("captchaId")
	// 获取验证码显示
	verifyValue := c.PostForm("verifyValue")

	// 获取账号和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 校验验证码是否正确
	if models.VerifyCaptcha(captchaId, verifyValue) {
		// c.String(http.StatusOK, "验证码验证成功")
		// 查询数据库，校验用户名和密码是否存在
		var userinfoList []models.Manager
		password = models.Md5(password)

		models.DB.Where("username = ? and password = ?", username, password).Find(&userinfoList)
		if len(userinfoList) > 0 {
			// 执行登录 保存用户信息 执行跳转
			session := sessions.Default(c)
			// 注意：session.Set没法直接保存结构体对应的切片 把结构体转换成json字符串
			userinfoSlice, _ := json.Marshal(userinfoList)
			session.Set("userinfo", string(userinfoSlice))
			session.Save()
			con.Success(c, "登录成功", "/admin")
		} else {
			con.Error(c, "用户名或密码错误", "/admin/login")
		}
	} else {
		// c.String(http.StatusOK, "验证码验证失败")
		con.Error(c, "验证码验证失败", "/admin/login")
	}
}

// Captcha
// @Description 校验验证码
// @Author xYuan 2024-04-16 15:29:03
// @Param c
func (con LoginController) Captcha(c *gin.Context) {
	captcha, captchaImg, answer, err := models.MackCaptcha()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"captchaId": captcha, "captchaImg": captchaImg, "answer": answer, "err": err})
}

// LoginOut 退出登录
func (con LoginController) LoginOut(c *gin.Context) {
	sessions := sessions.Default(c)
	sessions.Delete("userinfo")
	sessions.Save()
	con.Success(c, "退出成功", "/admin/login")
}
