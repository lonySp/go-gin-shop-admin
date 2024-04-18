package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lonySp/go-gin-shop-admin/controllers/xyuan"
	"github.com/lonySp/go-gin-shop-admin/middlewares"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", xyuan.DefaultController{}.Index)
		defaultRouters.GET("/category:id", xyuan.ProductController{}.Category)
		defaultRouters.GET("/detail", xyuan.ProductController{}.Detail)
		defaultRouters.GET("/product/getImgList", xyuan.ProductController{}.GetImgList)

		defaultRouters.GET("/cart", xyuan.CartController{}.Get)
		defaultRouters.GET("/cart/addCart", xyuan.CartController{}.AddCart)

		defaultRouters.GET("/cart/successTip", xyuan.CartController{}.AddCartSuccess)
		defaultRouters.GET("/cart/decCart", xyuan.CartController{}.DecCart)
		defaultRouters.GET("/cart/incCart", xyuan.CartController{}.IncCart)

		defaultRouters.GET("/cart/changeOneCart", xyuan.CartController{}.ChangeOneCart)
		defaultRouters.GET("/cart/changeAllCart", xyuan.CartController{}.ChangeAllCart)
		defaultRouters.GET("/cart/delCart", xyuan.CartController{}.DelCart)

		defaultRouters.GET("/pass/login", xyuan.PassController{}.Login)
		defaultRouters.GET("/pass/captcha", xyuan.PassController{}.Captcha)

		defaultRouters.GET("/pass/registerStep1", xyuan.PassController{}.RegisterStep1)
		defaultRouters.GET("/pass/registerStep2", xyuan.PassController{}.RegisterStep2)
		defaultRouters.GET("/pass/registerStep3", xyuan.PassController{}.RegisterStep3)
		defaultRouters.GET("/pass/sendCode", xyuan.PassController{}.SendCode)
		defaultRouters.GET("/pass/validateSmsCode", xyuan.PassController{}.ValidateSmsCode)
		defaultRouters.POST("/pass/doRegister", xyuan.PassController{}.DoRegister)
		defaultRouters.POST("/pass/doLogin", xyuan.PassController{}.DoLogin)
		defaultRouters.GET("/pass/loginOut", xyuan.PassController{}.LoginOut)
		//判断用户权限
		defaultRouters.GET("/buy/checkout", middlewares.InitUserAuthMiddleware, xyuan.BuyController{}.Checkout) //判断用户权限
		defaultRouters.POST("/buy/doCheckout", middlewares.InitUserAuthMiddleware, xyuan.BuyController{}.DoCheckout)
		defaultRouters.GET("/buy/pay", middlewares.InitUserAuthMiddleware, xyuan.BuyController{}.Pay)
		defaultRouters.GET("/buy/orderPayStatus", middlewares.InitUserAuthMiddleware, xyuan.BuyController{}.OrderPayStatus)

		defaultRouters.POST("/address/addAddress", middlewares.InitUserAuthMiddleware, xyuan.AddressController{}.AddAddress)
		defaultRouters.POST("/address/editAddress", middlewares.InitUserAuthMiddleware, xyuan.AddressController{}.EditAddress)
		defaultRouters.GET("/address/changeDefaultAddress", middlewares.InitUserAuthMiddleware, xyuan.AddressController{}.ChangeDefaultAddress)
		defaultRouters.GET("/address/getOneAddressList", middlewares.InitUserAuthMiddleware, xyuan.AddressController{}.GetOneAddressList)

		defaultRouters.GET("/alipay", middlewares.InitUserAuthMiddleware, xyuan.AlipayController{}.Alipay)
		defaultRouters.POST("/alipayNotify", xyuan.AlipayController{}.AlipayNotify)
		defaultRouters.GET("/alipayReturn", middlewares.InitUserAuthMiddleware, xyuan.AlipayController{}.AlipayReturn)

		defaultRouters.GET("/wxpay", middlewares.InitUserAuthMiddleware, xyuan.WxpayController{}.Wxpay)
		defaultRouters.POST("/wxpay/notify", xyuan.WxpayController{}.WxpayNotify)

		defaultRouters.GET("/user", middlewares.InitUserAuthMiddleware, xyuan.UserController{}.Index)
		defaultRouters.GET("/user/order", middlewares.InitUserAuthMiddleware, xyuan.UserController{}.OrderList)
		defaultRouters.GET("/user/orderinfo", middlewares.InitUserAuthMiddleware, xyuan.UserController{}.OrderInfo)

		defaultRouters.GET("/search", xyuan.SearchController{}.Index)
		defaultRouters.GET("/search/getOne", xyuan.SearchController{}.GetOne)
		defaultRouters.GET("/search/addGoods", xyuan.SearchController{}.AddGoods)
		defaultRouters.GET("/search/updateGoods", xyuan.SearchController{}.UpdateGoods)
		defaultRouters.GET("/search/deleteGoods", xyuan.SearchController{}.DeleteGoods)
		defaultRouters.GET("/search/query", xyuan.SearchController{}.Query)
		defaultRouters.GET("/search/filterQuery", xyuan.SearchController{}.FilterQuery)
		defaultRouters.GET("/search/pagingQuery", xyuan.SearchController{}.PagingQuery)
	}

}
