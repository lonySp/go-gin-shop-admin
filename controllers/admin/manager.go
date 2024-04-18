package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"
	"strings"
)

type ManagerController struct {
	BaseController
}

// Index
// @Description 查询管理员列表
// @Author xYuan 2024-04-16 20:26:23
// @Param c
func (con ManagerController) Index(c *gin.Context) {
	var managerList []models.Manager
	models.DB.Preload("Role").Find(&managerList)
	c.HTML(http.StatusOK, "admin/manager/index.html", gin.H{"managerList": managerList})
}

// Add
// @Description 跳转管理员增加页面
// @Author xYuan 2024-04-16 20:26:40
// @Param c
func (con ManagerController) Add(c *gin.Context) {
	var roleList []models.Role
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/manager/add.html", gin.H{"roleList": roleList})
}

// DoAdd
// @Description 新增管理员
// @Author xYuan 2024-04-16 20:28:43
// @Param c
func (con ManagerController) DoAdd(c *gin.Context) {
	roleId, err1 := models.Int(c.PostForm("role_id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/manager/add")
	}
	username := strings.Trim(c.PostForm("username"), "")
	password := strings.Trim(c.PostForm("password"), "")
	email := strings.Trim(c.PostForm("email"), "")
	mobile := strings.Trim(c.PostForm("mobile"), "")
	// 用户名和密码长度是否合法
	if len(username) < 2 || len(password) < 6 {
		con.Error(c, "用户名或者密码长度不合法", "/admin/manager/add")
		return
	}
	// 判断管理是否存在
	var managerList []models.Manager
	models.DB.Where("username = ?", username).Find(&managerList)
	if len(managerList) > 0 {
		con.Error(c, "管理员已经存在", "/admin/manager/add")
		return
	}
	if err := models.DB.Create(&models.Manager{
		Username: username,
		Password: models.Md5(password),
		Email:    email,
		Mobile:   mobile,
		RoleId:   roleId,
		Status:   1,
		AddTime:  int(models.GetUnix()),
	}).Error; err != nil {
		con.Error(c, "新增管理员失败", "/admin/manager/add")
		return
	}
	con.Success(c, "新增管理员成功", "/admin/manager")
}

// Edit
// @Description 跳转管理员编辑页面
// @Author xYuan 2024-04-16 20:26:40
// @Param c
func (con ManagerController) Edit(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/manager")
	}
	var manager models.Manager
	if err := models.DB.Where(&models.Manager{Id: id}).First(&manager).Error; err != nil {
		con.Error(c, "管理员不存在", "/admin/manager")
	}
	var roleList []models.Role
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/manager/edit.html", gin.H{"manager": manager, "roleList": roleList})
}

// DoEdit
// @Description 执行修改管理员
// @Author xYuan 2024-04-16 21:07:09
// @Param c
func (con ManagerController) DoEdit(c *gin.Context) {
	managerId, err := models.Int(c.PostForm("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/manager/add")
		return
	}
	roleId, err := models.Int(c.PostForm("role_id"))
	if err != nil {
		con.Error(c, "角色ID错误", "/admin/manager/add")
		return
	}

	username := strings.TrimSpace(c.PostForm("username"))
	password := strings.TrimSpace(c.PostForm("password"))
	email := strings.TrimSpace(c.PostForm("email"))
	mobile := strings.TrimSpace(c.PostForm("mobile"))

	managerData := models.Manager{
		Username: username,
		Email:    email,
		Mobile:   mobile,
		RoleId:   roleId,
		AddTime:  int(models.GetUnix()),
	}

	if password != "" {
		if len(password) < 6 {
			con.Error(c, "密码长度不合法", "/admin/manager/edit?id="+models.String(managerId))
			return
		}
		managerData.Password = models.Md5(password)
	}

	if err := models.DB.Model(&models.Manager{Id: managerId}).Updates(managerData).Error; err != nil {
		con.Error(c, "修改管理员失败", "/admin/manager/edit.html")
		return
	}
	con.Success(c, "修改管理员成功", "/admin/manager")
}

// Delete
// @Description 删除管理员数据
// @Author xYuan 2024-04-16 21:26:48
// @Param c
func (con ManagerController) Delete(c *gin.Context) {
	managerId, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/manager/add")
		return
	}
	if err := models.DB.Delete(&models.Manager{Id: managerId}).Error; err != nil {
		con.Error(c, "删除管理员失败", "/admin/manager")
		return
	}
	con.Success(c, "删除管理员成功", "/admin/manager")
}
