package xyuan

import (
	"fmt"
	"github.com/lonySp/go-gin-shop-admin/models"
	"time"

	"github.com/gin-gonic/gin"
)

type DefaultController struct {
	BaseController
}

func (con DefaultController) Index(c *gin.Context) {

	//演示：设置cookie 获取cookie
	// models.Cookie.Set(c, "username", "李四")
	// var username string
	// models.Cookie.Get(c, "username", &username)
	// fmt.Println(username)

	timeStart := time.Now().UnixNano()
	//1、获取顶部导航 挪到了base.go里面

	//2、获取轮播图数据
	focusList := []models.Focus{}
	if hasFocusList := models.CacheDb.Get("focusList", &focusList); !hasFocusList {
		models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
		models.CacheDb.Set("focusList", focusList, 60*60)
	}

	//3、获取分类的数据 挪到了base.go里面

	//4、获取中间导航 挪到了base.go里面

	//手机
	phoneList := []models.Goods{}
	if hasPhoneList := models.CacheDb.Get("phoneList", &phoneList); !hasPhoneList {
		phoneList = models.GetGoodsByCategory(1, "best", 8)
		models.CacheDb.Set("phoneList", phoneList, 60*60)
	}

	//配件
	otherList := []models.Goods{}
	if hasOtherList := models.CacheDb.Get("otherList", &otherList); !hasOtherList {
		otherList = models.GetGoodsByCategory(9, "all", 1)
		models.CacheDb.Set("otherList", otherList, 60*60)
	}

	timeEnd := time.Now().UnixNano()

	fmt.Printf("执行时间：%v 毫秒", (timeEnd-timeStart)/1000000)

	con.Render(c, "xyuan/index/index.html", gin.H{
		"focusList": focusList,
		"phoneList": phoneList,
		"otherList": otherList,
	})
}
