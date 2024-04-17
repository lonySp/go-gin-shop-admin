package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"
	"os"
)

type FocusController struct {
	BaseController
}

// Index
// @Description 查询轮播图数据
// @Author xYuan 2024-04-17 15:28:15
// @Param c
func (con FocusController) Index(c *gin.Context) {
	var focus []models.Focus
	models.DB.Find(&focus)
	c.HTML(http.StatusOK, "admin/focus/index.html", gin.H{
		"focusList": focus,
	})
}
func (con FocusController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/focus/add.html", gin.H{})
}

// Edit
// @Description 重定向到修改页面
// @Author xYuan 2024-04-17 16:03:23
// @Param c
func (con FocusController) Edit(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/focus")
		return
	}
	fmt.Println("id:", id)
	focus := models.Focus{Id: id}
	if err := models.DB.Find(&focus).Error; err != nil {
		con.Error(c, "数据不存在", "/admin/focus")
		return
	}
	c.HTML(http.StatusOK, "admin/focus/edit.html", gin.H{
		"focus": focus,
	})
}

// DoAdd
// @Description 执行添加轮播图信息
// @Author xYuan 2024-04-17 15:58:43
// @Param c
func (con FocusController) DoAdd(c *gin.Context) {
	title := c.PostForm("title")
	focusType, err1 := models.Int(c.PostForm("focus_type"))
	focusImgSrc, err4 := models.UploadImg(c, "focus_img")
	link := c.PostForm("link")
	sort, err2 := models.Int(c.PostForm("sort"))
	status, err3 := models.Int(c.PostForm("status"))
	if err1 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/focus/add")
		return
	}
	if err2 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/add")
		return
	}
	if err4 != nil {
		fmt.Println(err4)
		return
	}
	if err := models.DB.Create(&models.Focus{
		Title:     title,
		FocusType: focusType,
		FocusImg:  focusImgSrc,
		Link:      link,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}).Error; err != nil {
		con.Error(c, "增加轮播图失败", "/admin/focus/add")
	} else {
		con.Success(c, "增加轮播图成功", "/admin/focus")
	}
}

// DoEdit
// @Description 执行修改
// @Author xYuan 2024-04-17 15:59:04
// @Param c
func (con FocusController) DoEdit(c *gin.Context) {
	id, err := models.Int(c.PostForm("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/focus")
		return
	}
	title := c.PostForm("title")
	focusType, err1 := models.Int(c.PostForm("focus_type"))
	focusImgSrc, err4 := models.UploadImg(c, "focus_img")
	link := c.PostForm("link")
	sort, err2 := models.Int(c.PostForm("sort"))
	status, err3 := models.Int(c.PostForm("status"))
	if err1 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/focus/add")
		return
	}
	if err2 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/add")
		return
	}
	if err4 != nil {
		fmt.Println(err4)
		return
	}
	focus := models.Focus{
		Title:     title,
		FocusType: focusType,
		Link:      link,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}
	if focusImgSrc != "" {
		focus.FocusImg = focusImgSrc
	}
	if err := models.DB.Where("id = ?", id).Updates(&focus).Error; err != nil {
		con.Error(c, "修改轮播图失败，请重新尝试", "/admin/focus/edit?id="+models.String(id))
	} else {
		con.Success(c, "修改轮播图成功", "/admin/focus")
	}
}

// Delete
// @Description 删除轮播图
// @Author xYuan 2024-04-17 16:03:08
// @Param c
func (con FocusController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "参数错误", "/admin/focus")
		return
	}
	// 先查询出要删除的记录
	var focus models.Focus
	if err := models.DB.Where("id = ?", id).First(&focus).Error; err != nil {
		con.Error(c, "找不到指定记录", "/admin/focus")
		return
	}
	// 删除图片文件，如果有的话
	if focus.FocusImg != "" {
		if err := os.Remove(focus.FocusImg); err != nil {
			con.Error(c, "删除图片失败: "+err.Error(), "/admin/focus")
			return
		}
	}
	// 删除数据库记录
	if err := models.DB.Delete(&focus).Error; err != nil {
		con.Error(c, "删除记录失败: "+err.Error(), "/admin/focus")
		return
	}
	con.Success(c, "删除成功", "/admin/focus")
}
