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

// Auth
// @Description 重定向授权界面
// @Author xYuan 2024-04-17 11:09:08
// @Param c
func (con RoleController) Auth(c *gin.Context) {
	// 1、获取角色ID
	roleId, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/role")
		return
	}

	// 2、获取所有的权限并预加载关联的AccessItem
	var accessList []models.Access
	if err := models.DB.Where("module_id = ?", 0).Preload("AccessItem").Find(&accessList).Error; err != nil {
		con.Error(c, "无法加载权限数据", "/admin/role")
		return
	}

	// 3、获取当前角色拥有的权限，并把权限id放在一个map对象里面
	var roleAccess []models.RoleAccess
	if err := models.DB.Where("role_id = ?", roleId).Find(&roleAccess).Error; err != nil {
		con.Error(c, "无法加载角色权限", "/admin/role")
		return
	}
	roleAccessMap := make(map[int]bool)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = true
	}

	// 4、遍历所有的权限数据，判断当前权限的id是否在角色权限的Map对象中，如果是的话给当前数据加入checked属性
	for i, access := range accessList {
		if _, ok := roleAccessMap[access.Id]; ok {
			accessList[i].Checked = true
		}
		for j, item := range accessList[i].AccessItem {
			if _, ok := roleAccessMap[item.Id]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}
	c.HTML(http.StatusOK, "admin/role/auth.html", gin.H{
		"roleId":     roleId,
		"accessList": accessList,
	})
}

// DoAuth
// @Description 执行角色权限选择
// @Author xYuan 2024-04-17 11:08:25
// @Param c
func (con RoleController) DoAuth(c *gin.Context) {
	// 获取角色ID
	roleId, err := models.Int(c.PostForm("role_id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/role")
		return
	}

	// 获取权限id切片
	accessIds := c.PostFormArray("access_node[]")

	// 删除当前角色对应的所有权限
	if err := models.DB.Where("role_id = ?", roleId).Delete(&models.RoleAccess{}).Error; err != nil {
		con.Error(c, "删除旧权限失败: "+err.Error(), "/admin/role")
		return
	}

	// 准备批量插入新的权限
	roleAccesses := make([]models.RoleAccess, len(accessIds))
	for i, v := range accessIds {
		accessId, err := models.Int(v)
		if err != nil {
			con.Error(c, "无效的权限ID: "+err.Error(), "/admin/role")
			return
		}
		roleAccesses[i] = models.RoleAccess{RoleId: roleId, AccessId: accessId}
	}

	// 批量创建新权限
	if err := models.DB.Create(&roleAccesses).Error; err != nil {
		con.Error(c, "添加新权限失败: "+err.Error(), "/admin/role")
		return
	}

	// 重定向到授权页面，显示授权成功
	con.Success(c, "授权成功", "/admin/role/auth?id="+models.String(roleId))
}
