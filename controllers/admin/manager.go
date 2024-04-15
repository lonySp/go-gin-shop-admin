package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	BaseController
}

func (con ManagerController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/manager/index.html", gin.H{})

}
func (con ManagerController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/manager/add.html", gin.H{})
}

func (con ManagerController) Edit(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/manager/edit.html", gin.H{})
}
func (con ManagerController) Delete(c *gin.Context) {
	c.String(http.StatusOK, "-add--文章-")
}
