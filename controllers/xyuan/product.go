package xyuan

import (
	"github.com/lonySp/go-gin-shop-admin/models"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	BaseController
}

func (con ProductController) Category(c *gin.Context) {

	//分类id
	cateId, _ := models.Int(c.Param("id"))
	//当前页
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	//每一页显示的数量
	pageSize := 5
	//获取当前分类
	currentCate := models.GoodsCate{}
	models.DB.Where("id=?", cateId).Find(&currentCate)
	subCate := []models.GoodsCate{}
	var tempSlice []int
	if currentCate.Pid == 0 {
		//获取二级分类
		models.DB.Where("pid=?", currentCate.Id).Find(&subCate)
		for i := 0; i < len(subCate); i++ {
			tempSlice = append(tempSlice, subCate[i].Id)
		}
	} else {
		//兄弟分类
		models.DB.Where("pid=?", currentCate.Pid).Find(&subCate)
	}
	tempSlice = append(tempSlice, cateId)
	where := "cate_id in ?"
	goodsList := []models.Goods{}
	models.DB.Where(where, tempSlice).Offset((page - 1) * pageSize).Limit(pageSize).Find(&goodsList)

	//获取总数量
	var count int64
	models.DB.Where(where, tempSlice).Table("goods").Count(&count)

	//自定义模板
	//https://www.mi.com/p/3469.html
	tpl := "xyuan/product/list.html"
	if currentCate.Template != "" {
		tpl = currentCate.Template
	}

	con.Render(c, tpl, gin.H{
		"page":        page,
		"goodsList":   goodsList,
		"subCate":     subCate,
		"currentCate": currentCate,
		"totalPages":  math.Ceil(float64(count) / float64(pageSize)),
	})

}

func (con ProductController) Detail(c *gin.Context) {

	id, err := models.Int(c.Query("id"))

	if err != nil {
		c.Redirect(302, "/")
		c.Abort()
	}

	//1、获取商品信息
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)

	//2、获取关联商品  RelationGoods
	relationGoods := []models.Goods{}
	goods.RelationGoods = strings.ReplaceAll(goods.RelationGoods, "，", ",")
	relationIds := strings.Split(goods.RelationGoods, ",")

	models.DB.Where("id in ?", relationIds).Select("id,title,price,goods_version").Find(&relationGoods)

	//3、获取关联赠品 GoodsGift

	goodsGift := []models.Goods{}
	goods.GoodsGift = strings.ReplaceAll(goods.GoodsGift, "，", ",")
	giftIds := strings.Split(goods.GoodsGift, ",")
	models.DB.Where("id in ?", giftIds).Select("id,title,price,goods_version").Find(&goodsGift)

	//4、获取关联颜色 GoodsColor
	goodsColor := []models.GoodsColor{}
	goods.GoodsColor = strings.ReplaceAll(goods.GoodsColor, "，", ",")
	colorIds := strings.Split(goods.GoodsColor, ",")
	models.DB.Where("id in ?", colorIds).Find(&goodsColor)

	//5、获取关联配件 GoodsFitting
	goodsFitting := []models.Goods{}
	goods.GoodsFitting = strings.ReplaceAll(goods.GoodsFitting, "，", ",")
	fittingIds := strings.Split(goods.GoodsFitting, ",")
	models.DB.Where("id in ?", fittingIds).Select("id,title,price,goods_version").Find(&goodsFitting)

	//6、获取商品关联的图片 GoodsImage
	goodsImage := []models.GoodsImage{}
	models.DB.Where("goods_id=?", goods.Id).Limit(6).Find(&goodsImage)

	//7、获取规格参数信息 GoodsAttr
	goodsAttr := []models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsAttr)

	//8、获取更多属性

	/*
			颜色:红色,白色,黄色 | 尺寸:41,42,43

			切片

			[
				{
					Cate:"颜色",
					List:[红色,白色,黄色]
				},
				{
					Cate:"尺寸",
					List:[41,42,43]
				}
			]


		goodsAttrStrSlice[0]	尺寸:41,42,43

				tempSlice[0]    尺寸

				tempSlice[1]	41,42,43

		goodsAttrStrSlice[1]	套餐:套餐1,套餐2

	*/

	// goodsAttrStr := "尺寸:41,42,43|套餐:套餐1,套餐2"
	goodsAttrStr := goods.GoodsAttr
	goodsAttrStr = strings.ReplaceAll(goodsAttrStr, "，", ",")
	goodsAttrStr = strings.ReplaceAll(goodsAttrStr, "：", ":")

	var goodsItemAttrList []models.GoodsItemAttr
	if strings.Contains(goodsAttrStr, ":") {
		goodsAttrStrSlice := strings.Split(goodsAttrStr, "|")
		//创建切片的存储空间
		goodsItemAttrList = make([]models.GoodsItemAttr, len(goodsAttrStrSlice))
		for i := 0; i < len(goodsAttrStrSlice); i++ {
			tempSlice := strings.Split(goodsAttrStrSlice[i], ":")
			goodsItemAttrList[i].Cate = tempSlice[0]
			listSlice := strings.Split(tempSlice[1], ",")
			goodsItemAttrList[i].List = listSlice
		}
	}

	// c.JSON(200, gin.H{
	// 	"goodsItemAttrList": goodsItemAttrList,
	// })

	// c.String(200, "Detail")
	tpl := "xyuan/product/detail.html"

	con.Render(c, tpl, gin.H{
		"goods":             goods,
		"relationGoods":     relationGoods,
		"goodsGift":         goodsGift,
		"goodsColor":        goodsColor,
		"goodsFitting":      goodsFitting,
		"goodsImage":        goodsImage,
		"goodsAttr":         goodsAttr,
		"goodsItemAttrList": goodsItemAttrList,
	})
}

func (con ProductController) GetImgList(c *gin.Context) {

	goodsId, err1 := models.Int(c.Query("goods_id"))
	colorId, err2 := models.Int(c.Query("color_id"))

	//查询商品图库信息

	goodsImageList := []models.GoodsImage{}
	err3 := models.DB.Where("goods_id=? AND color_id=?", goodsId, colorId).Find(&goodsImageList).Error

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(200, gin.H{
			"success": false,
			"result":  "",
			"message": "参数错误",
		})
	} else {

		//判断 goodsImageList的长度 如果goodsImageList没有数据，那么我们需要返回当前商品所有的图库信息
		if len(goodsImageList) == 0 {
			models.DB.Where("goods_id=?", goodsId).Find(&goodsImageList)
		}
		c.JSON(200, gin.H{
			"success": true,
			"result":  goodsImageList,
			"message": "获取数据成功",
		})
	}
}
