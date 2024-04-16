package admin

import (
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	BaseController
}

// Index
// @Description 查询所有角色
// @Author xYuan 2024-04-16 17:26:26
// @Param c
func (con RoleController) Index(c *gin.Context) {
	var roleList []models.Role
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/role/index.html", gin.H{"roleList": roleList})

}

// Add
// @Description 重定向到添加页面
// @Author xYuan 2024-04-16 17:26:33
// @Param c
func (con RoleController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/role/add.html", gin.H{})
}

// DoAdd
// @Description 执行添加角色
// @Author xYuan 2024-04-16 17:26:43
// @Param c
func (con RoleController) DoAdd(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), "")
	description := strings.Trim(c.PostForm("description"), "")
	if title == "" {
		con.Error(c, "角色名称不能为空", "/admin/role/add")
		return
	}
	err := models.DB.Create(&models.Role{
		Title:       title,
		Description: description,
		Status:      1,
		AddTime:     int(models.GetUnix()),
	}).Error
	if err != nil {
		con.Error(c, "添加角色失败", "/admin/role/add")
	} else {
		con.Success(c, "添加角色成功", "/admin/role")
	}
}

// Edit
// @Description 重定向到修改角色页面
// @Author xYuan 2024-04-16 17:39:46
// @Param c
func (con RoleController) Edit(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/role")
		return
	} else {
		role := models.Role{Id: id}
		models.DB.Find(&role)
		c.HTML(http.StatusOK, "admin/role/edit.html", gin.H{"role": role})
	}
}

// DoEdit
// @Description 执行修改角色
// @Author xYuan 2024-04-16 17:40:06
// @Param c
func (con RoleController) DoEdit(c *gin.Context) {
	id, err := models.Int(c.PostForm("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/role")
		return
	}
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	if title == "" {
		con.Error(c, "角色名称不能为空", "/admin/role/add")
		return
	}
	updateErr := models.DB.Model(&models.Role{}).Where("id = ?", id).Updates(&models.Role{
		Title:       title,
		Description: description,
	}).Error
	if updateErr != nil {
		con.Error(c, "修改角色失败", "/admin/role")
		return
	}
	con.Success(c, "修改角色成功", "/admin/role")
}

// Delete
// @Description 删除角色
// @Author xYuan 2024-04-16 18:02:39
// @Param c
func (con RoleController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/role")
		return
	}
	if err := models.DB.Delete(&models.Role{}, id).Error; err != nil {
		con.Error(c, "删除角色失败", "/admin/role")
		return
	}
	con.Success(c, "删除角色成功", "/admin/role")
}
