package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"
	"strings"
)

type AccessController struct {
	BaseController
}

// Index
// @Description 查询权限列表
// @Author xYuan 2024-04-17 09:30:39
// @Param c
func (con AccessController) Index(c *gin.Context) {
	var accessList []models.Access
	models.DB.Where(&models.Access{ModuleId: 0}).Preload("AccessItem").Find(&accessList)
	c.HTML(http.StatusOK, "admin/access/index.html", gin.H{
		"accessList": accessList,
	})
}

// Add
// @Description 重定向权限添加页面
// @Author xYuan 2024-04-17 09:36:43
// @Param c
func (con AccessController) Add(c *gin.Context) {
	//获取顶级模块
	var accessList []models.Access
	models.DB.Where(&models.Access{ModuleId: 0}).Find(&accessList)
	c.HTML(http.StatusOK, "admin/access/add.html", gin.H{
		"accessList": accessList,
	})
}

// DoAdd
// @Description 执行添加权限
// @Author xYuan 2024-04-17 09:38:50
// @Param c
func (con AccessController) DoAdd(c *gin.Context) {
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := c.PostForm("action_name")
	url := c.PostForm("url")
	description := c.PostForm("description")
	accessType, err1 := models.Int(c.PostForm("type"))
	moduleId, err2 := models.Int(c.PostForm("module_id"))
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/access/add")
		return
	}
	if moduleName == "" {
		con.Error(c, "模块名称不能为空", "/admin/access/add")
		return
	}
	if err := models.DB.Create(&models.Access{
		ModuleName:  moduleName,
		Type:        accessType,
		ActionName:  actionName,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Description: description,
		Status:      status,
	}).Error; err != nil {
		con.Error(c, "增加数据失败", "/admin/access/add")
		return
	}
	con.Success(c, "增加数据成功", "/admin/access")

}

// Edit
// @Description 重定向到修改页面
// @Author xYuan 2024-04-17 09:55:13
// @Param c
func (con AccessController) Edit(c *gin.Context) {
	//获取要修改的数据
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "参数错误", "/admin/access")
	}
	access := models.Access{Id: id}
	models.DB.Find(&access)
	//获取顶级模块
	var accessList []models.Access
	models.DB.Where(&models.Access{ModuleId: 0}).Find(&accessList)
	c.HTML(http.StatusOK, "admin/access/edit.html", gin.H{
		"access":     access,
		"accessList": accessList,
	})
}

// DoEdit
// @Description 执行修改
// @Author xYuan 2024-04-17 09:56:54
// @Param c
func (con AccessController) DoEdit(c *gin.Context) {
	id, err := models.Int(c.PostForm("id"))
	if err != nil {
		con.Error(c, "传入参数错误", "/admin/access")
		return
	}

	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	if moduleName == "" {
		con.Error(c, "模块名称不能为空", fmt.Sprintf("/admin/access/edit?id=%d", id))
		return
	}

	// 收集并验证其它表单字段
	accessType, err1 := models.Int(c.PostForm("type"))
	moduleId, err2 := models.Int(c.PostForm("module_id"))
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))

	// 统一处理字段转换错误
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "表单参数类型错误", fmt.Sprintf("/admin/access/edit?id=%d", id))
		return
	}

	// 构建更新模型
	access := models.Access{
		Id:          id,
		ModuleName:  moduleName,
		Type:        accessType,
		ActionName:  c.PostForm("action_name"),
		Url:         c.PostForm("url"),
		ModuleId:    moduleId,
		Sort:        sort,
		Description: c.PostForm("description"),
		Status:      status,
	}

	// 执行更新操作
	if err := models.DB.Model(&models.Access{}).Where("id = ?", id).Updates(access).Error; err != nil {
		con.Error(c, "修改数据失败", fmt.Sprintf("/admin/access/edit?id=%d", id))
	} else {
		con.Success(c, "修改数据成功", fmt.Sprintf("/admin/access/edit?id=%d", id))
	}
}

// Delete
// @Description 删除权限操作
// @Author xYuan 2024-04-17 10:08:01
// @Param c
func (con AccessController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/access")
		return
	}

	//获取我们要删除的数据
	var access models.Access
	if err := models.DB.Where("id = ?", id).Take(&access).Error; err != nil {
		con.Error(c, "找不到指定的数据", "/admin/access")
		return
	}

	// 检查是否是顶级模块并且是否有子模块
	if access.ModuleId == 0 {
		var count int64
		models.DB.Model(&models.Access{}).Where("module_id = ?", access.Id).Count(&count)
		if count > 0 {
			con.Error(c, "当前模块下面有菜单或者操作，请删除菜单或者操作以后再来删除这个数据", "/admin/access")
			return
		}
	}

	// 删除操作或顶级模块（如果没有子模块）
	if err := models.DB.Delete(&access).Error; err != nil {
		con.Error(c, "删除数据失败", "/admin/access")
	} else {
		con.Success(c, "删除数据成功", "/admin/access")
	}
}
