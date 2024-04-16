package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/controllers/admin"
	"github.com/lonySp/go-gin-shop-admin/middlewares"
)

func AdminRoutersInit(r *gin.Engine) {
	//middlewares.InitMiddleware中间件
	adminRouters := r.Group("/admin", middlewares.InitAdminAuthMiddleware)
	{
		adminRouters.GET("/", admin.MainController{}.Index)
		adminRouters.GET("/welcome", admin.MainController{}.Welcome)

		adminRouters.GET("/login", admin.LoginController{}.Index)
		adminRouters.POST("/doLogin", admin.LoginController{}.DoLogin)
		adminRouters.GET("/loginOut", admin.LoginController{}.LoginOut)
		adminRouters.GET("/captcha", admin.LoginController{}.Captcha)

		adminRouters.GET("/focus", admin.FocusController{}.Index)
		adminRouters.GET("/focus/add", admin.FocusController{}.Add)
		adminRouters.GET("/focus/edit", admin.FocusController{}.Edit)
		adminRouters.GET("/focus/delete", admin.FocusController{}.Delete)

		adminRouters.GET("/role", admin.RoleController{}.Index)
		adminRouters.GET("/role/add", admin.RoleController{}.Add)
		adminRouters.POST("/role/doAdd", admin.RoleController{}.DoAdd)
		adminRouters.POST("/role/doEdit", admin.RoleController{}.DoEdit)
		adminRouters.GET("/role/edit", admin.RoleController{}.Edit)
		adminRouters.GET("/role/delete", admin.RoleController{}.Delete)

		adminRouters.GET("/manager", admin.ManagerController{}.Index)
		adminRouters.GET("/manager/add", admin.ManagerController{}.Add)
		adminRouters.POST("/manager/doAdd", admin.ManagerController{}.DoAdd)
		adminRouters.POST("/manager/doEdit", admin.ManagerController{}.DoEdit)
		adminRouters.GET("/manager/edit", admin.ManagerController{}.Edit)
		adminRouters.GET("/manager/delete", admin.ManagerController{}.Delete)

	}
}
