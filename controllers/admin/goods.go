package admin

import (
	"fmt"
	"github.com/lonySp/go-gin-shop-admin/models"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

type GoodsController struct {
	BaseController
}

func (con GoodsController) Index(c *gin.Context) {

	goodsList := []models.Goods{}
	models.DB.Find(&goodsList)

	c.HTML(http.StatusOK, "admin/goods/index.html", gin.H{
		"goodsList": goodsList,
	})
}
func (con GoodsController) Add(c *gin.Context) {
	//获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	//获取所有颜色信息
	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)

	//获取商品规格包装
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "admin/goods/add.html", gin.H{
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
	})
}

func (con GoodsController) GoodsTypeAttribute(c *gin.Context) {
	cateId, err1 := models.Int(c.Query("cateId"))
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	err2 := models.DB.Where("cate_id = ?", cateId).Find(&goodsTypeAttributeList).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"result":  "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  goodsTypeAttributeList,
		})
	}
}

func (con GoodsController) DoAdd(c *gin.Context) {

	//1、获取表单提交过来的数据 进行判断

	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	goodsSn := c.PostForm("goods_sn")
	cateId, _ := models.Int(c.PostForm("cate_id"))
	goodsNumber, _ := models.Int(c.PostForm("goods_number"))
	//注意小数点
	marketPrice, _ := models.Float(c.PostForm("market_price"))
	price, _ := models.Float(c.PostForm("price"))

	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr")
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color")

	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	isDelete, _ := models.Int(c.PostForm("is_delete"))
	isHot, _ := models.Int(c.PostForm("is_hot"))
	isBest, _ := models.Int(c.PostForm("is_best"))
	isNew, _ := models.Int(c.PostForm("is_new"))
	goodsTypeId, _ := models.Int(c.PostForm("goods_type_id"))
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))
	addTime := int(models.GetUnix())

	//2、获取颜色信息 把颜色转化成字符串
	goodsColorStr := strings.Join(goodsColorArr, ",")

	//3、上传图片   生成缩略图
	goodsImg, _ := models.UploadImg(c, "goods_img")

	//4、增加商品数据

	goods := models.Goods{
		Title:         title,
		SubTitle:      subTitle,
		GoodsSn:       goodsSn,
		CateId:        cateId,
		ClickCount:    100,
		GoodsNumber:   goodsNumber,
		MarketPrice:   marketPrice,
		Price:         price,
		RelationGoods: relationGoods,
		GoodsAttr:     goodsAttr,
		GoodsVersion:  goodsVersion,
		GoodsGift:     goodsGift,
		GoodsFitting:  goodsFitting,
		GoodsKeywords: goodsKeywords,
		GoodsDesc:     goodsDesc,
		GoodsContent:  goodsContent,
		IsDelete:      isDelete,
		IsHot:         isHot,
		IsBest:        isBest,
		IsNew:         isNew,
		GoodsTypeId:   goodsTypeId,
		Sort:          sort,
		Status:        status,
		AddTime:       addTime,
		GoodsColor:    goodsColorStr,
		GoodsImg:      goodsImg,
	}
	err := models.DB.Create(&goods).Error
	if err != nil {
		con.Error(c, "增加失败", "/admin/goods/add")
	}
	//5、增加图库 信息
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()
	//6、增加规格包装
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, attributeIdErr := models.Int(attrIdList[i])
			if attributeIdErr == nil {
				//获取商品类型属性的数据
				goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
				models.DB.Find(&goodsTypeAttributeObj)
				//给商品属性里面增加数据  规格包装
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
				goodsAttrObj.AttributeValue = attrValueList[i]
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}

		}
		wg.Done()
	}()
	wg.Wait()
	con.Success(c, "增加数据成功", "/admin/goods")
}

// 修改
func (con GoodsController) Edit(c *gin.Context) {

	// 1、获取要修改的商品数据
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数错误", "/admin/goods")
	}
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)

	// 2、获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	// 3、获取所有颜色 以及选中的颜色
	goodsColorSlice := strings.Split(goods.GoodsColor, ",")
	goodsColorMap := make(map[string]string)
	for _, v := range goodsColorSlice {
		goodsColorMap[v] = v
	}

	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)
	for i := 0; i < len(goodsColorList); i++ {
		if _, ok := goodsColorMap[models.String(goodsColorList[i].Id)]; ok {
			goodsColorList[i].Checked = true
		}
	}

	// 4、商品的图库信息
	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsImageList)

	// 5、获取商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	// 6、获取规格信息
	goodsAttr := []models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsAttr)
	goodsAttrStr := ""

	for _, v := range goodsAttr {
		if v.AttributeType == 1 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: </span> <input type="hidden" name="attr_id_list" value="%v" />   <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else if v.AttributeType == 2 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: 　</span><input type="hidden" name="attr_id_list" value="%v" />  <textarea cols="50" rows="3" name="attr_value_list">%v</textarea></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else {
			//获取当前类型对应的值
			goodsTypeArttribute := models.GoodsTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&goodsTypeArttribute)
			attrValueSlice := strings.Split(goodsTypeArttribute.AttrValue, "\n")

			goodsAttrStr += fmt.Sprintf(`<li><span>%v: 　</span>  <input type="hidden" name="attr_id_list" value="%v" /> `, v.AttributeTitle, v.AttributeId)
			goodsAttrStr += fmt.Sprintf(`<select name="attr_value_list">`)
			for i := 0; i < len(attrValueSlice); i++ {
				if attrValueSlice[i] == v.AttributeValue {
					goodsAttrStr += fmt.Sprintf(`<option value="%v" selected >%v</option>`, attrValueSlice[i], attrValueSlice[i])
				} else {
					goodsAttrStr += fmt.Sprintf(`<option value="%v">%v</option>`, attrValueSlice[i], attrValueSlice[i])
				}
			}
			goodsAttrStr += fmt.Sprintf(`</select>`)
			goodsAttrStr += fmt.Sprintf(`</li>`)

		}
	}

	c.HTML(http.StatusOK, "admin/goods/edit.html", gin.H{
		"goods":          goods,
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
		"goodsAttrStr":   goodsAttrStr,
		"goodsImageList": goodsImageList,
	})
}
func (con GoodsController) ImageUpload(c *gin.Context) {
	//上传图片
	imgDir, err := models.UploadImg(c, "file") //注意：可以在网络里面看到传递的参数
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"link": "/" + imgDir,
		})
	}
}
