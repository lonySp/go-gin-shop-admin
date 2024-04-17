package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/goccy/go-json"
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MainController struct{}

// Index
// @Description 重定向到管理页面，并进行权限筛查
// @Author xYuan 2024-04-17 12:10:35
// @Param c
func (con MainController) Index(c *gin.Context) {
	// 获取 userinfo 对应的 session
	session := sessions.Default(c)
	userinfo, ok := session.Get("userinfo").(string)
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	// 类型断言后解析用户信息
	var managers []models.Manager
	if err := json.Unmarshal([]byte(userinfo), &managers); err != nil || len(managers) == 0 {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	manager := managers[0]

	// 获取所有顶级权限
	var accessList []models.Access
	models.DB.Where("module_id = ?", 0).Preload("AccessItem").Find(&accessList)

	// 获取当前角色的权限，并存储在 map 中
	var roleAccess []models.RoleAccess
	models.DB.Where("role_id = ?", manager.RoleId).Find(&roleAccess)
	roleAccessMap := make(map[int]bool)
	for _, access := range roleAccess {
		roleAccessMap[access.AccessId] = true
	}

	// 标记拥有的权限
	for i, access := range accessList {
		if _, ok := roleAccessMap[access.Id]; ok {
			accessList[i].Checked = true
		}
		for j, item := range access.AccessItem {
			if _, ok := roleAccessMap[item.Id]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}

	// 返回页面
	c.HTML(http.StatusOK, "admin/main/index.html", gin.H{
		"username":   manager.Username,
		"accessList": accessList,
		"isSuper":    manager.IsSuper,
	})
}

func (con MainController) Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/main/welcome.html", gin.H{})
}
