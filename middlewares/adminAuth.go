package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/models"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

// 定义常量，用于在函数中引用，减少硬编码
const (
	AdminLoginURL   = "/admin/login"
	AdminDoLoginURL = "/admin/doLogin"
	AdminCaptchaURL = "/admin/captcha"
)

// InitAdminAuthMiddleware 初始化后台权限认证中间件
func InitAdminAuthMiddleware(c *gin.Context) {
	pathname := strings.Split(c.Request.URL.String(), "?")[0]
	session := sessions.Default(c)
	userinfo, ok := session.Get("userinfo").(string)
	if ok {
		if !validateUserSession(userinfo) {
			handleUnauthorizedAccess(pathname, c)
			return
		}
		validateUserPermissions(userinfo, pathname, c)
	} else {
		handleUnauthorizedAccess(pathname, c)
	}
}

// handleUnauthorizedAccess 处理未授权访问，如未登录用户
func handleUnauthorizedAccess(pathname string, c *gin.Context) {
	if pathname != AdminLoginURL && pathname != AdminDoLoginURL && pathname != AdminCaptchaURL {
		c.Redirect(302, AdminLoginURL)
	}
}

// validateUserSession 验证用户的会话信息是否合法
func validateUserSession(userinfo string) bool {
	var userinfoStruct []models.Manager
	if err := json.Unmarshal([]byte(userinfo), &userinfoStruct); err != nil || len(userinfoStruct) == 0 || userinfoStruct[0].Username == "" {
		return false
	}
	return true
}

// validateUserPermissions 验证用户的访问权限
func validateUserPermissions(userinfo, pathname string, c *gin.Context) {
	urlPath := strings.Replace(pathname, "/admin/", "", 1)
	if !excludeAuthPath("/" + urlPath) {
		checkPermissions(userinfo, urlPath, c)
	}
}

// checkPermissions 检查用户是否有权限访问特定URL
func checkPermissions(userinfo, urlPath string, c *gin.Context) {
	var userinfoStruct []models.Manager
	json.Unmarshal([]byte(userinfo), &userinfoStruct)
	if userinfoStruct[0].IsSuper == 1 { // 超级管理员无视权限检查
		return
	}

	var roleAccessList []models.RoleAccess
	models.DB.Where("role_id = ?", userinfoStruct[0].RoleId).Find(&roleAccessList)
	accessMap := make(map[int]struct{})
	for _, access := range roleAccessList {
		accessMap[access.AccessId] = struct{}{}
	}

	access := models.Access{}
	models.DB.Where("url = ?", urlPath).Find(&access)
	if _, exists := accessMap[access.Id]; !exists {
		c.String(403, "没有权限")
		c.Abort()
	}
}

// excludeAuthPath 判断访问路径是否在配置文件中的排除权限检查列表
func excludeAuthPath(urlPath string) bool {
	config, err := ini.Load("./conf/app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	excludePaths := config.Section("").Key("excludeAuthPath").String()
	excludePathSlice := strings.Split(excludePaths, ",")
	for _, path := range excludePathSlice {
		if path == urlPath {
			return true
		}
	}
	return false
}
