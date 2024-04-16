package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/lonySp/go-gin-shop-admin/models"
	"strings"
)

// InitAdminAuthMiddleware
// @Description 登录权限判断
// @Author xYuan 2024-04-16 15:13:35
// @Param c
func InitAdminAuthMiddleware(c *gin.Context) {
	//进行权限判断 没有登录的用户不能进入后台管理中心
	// 1、获取Url访问地址
	// 2、获取Session里面保存的用户信息
	// 3、判断Session中的用户信息是否存在，如果不存在，则跳转到登录页面 ,如果Session中的用户信息存在，则继续执行
	// 4、如果Session不存在，判断当前访问的URl是否是login doLogin captcha，如果不是跳转到登录页面，如果是不行任何操作

	// 1、获取Url访问地址
	pathname := strings.Split(c.Request.URL.String(), "?")[0]
	// 2、获取Session里面保存的用户信息
	session := sessions.Default(c)
	userinfo := session.Get("userinfo")

	//类型断言 来判断 userinfo是不是一个string
	userinfoStr, ok := userinfo.(string)
	if ok {
		var userinfoStruct []models.Manager
		// json字符串转为结构体对象
		err := json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		// 转换成功或者没有数据并且不是登录页面 则重定向到登录页面
		if err != nil || !(len(userinfoStruct) > 0 && userinfoStruct[0].Username != "") {
			if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
				c.Redirect(302, "/admin/login")
			}
		}
	} else {
		if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
			c.Redirect(302, "/admin/login")
		}
	}
}
