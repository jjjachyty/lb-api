package router

import (
	"fmt"
	"lb-api/controlers"
	"lb-api/controlers/order"
	"lb-api/controlers/purchase"
	"lb-api/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	api := e.Group("/api/v1")

	api.Use(func(c *gin.Context) {
		fmt.Println("c.Request()\n\n", c.Request.Method, c.Request.URL.Path)
		if "Get" != c.Request.Method { //更改数据
			// user := controlers.GetInfoFromToken(c)
		}
	})

	api.GET("/captcha/:id", controlers.Captcha)
	api.GET("/avatars", controlers.AvatarList)
	//用户登陆
	authMiddleware := middlewares.JWT()
	api.POST("/login", authMiddleware.LoginHandler)

	// api.POST("/emailvalid", controlers.ValidEmail)

	user := api.Group("/user")
	//重新获取token
	user.Use(authMiddleware.MiddlewareFunc())
	user.GET("/refreshtoken", authMiddleware.RefreshHandler)

	//用户实名认证
	user.POST("/scanidcard", controlers.UserControl{}.IDCardOCR)
	user.POST("/valididcard", controlers.UserControl{}.ValidIDCard)
	//用户设置
	user.GET("/info", controlers.UserControl{}.GetUserInfo)
	user.PUT("/info", controlers.UserControl{}.Update)
	//用户介绍
	api.GET("/profile/:id", controlers.UserControl{}.UserProfile)

	//用户充值
	user.POST("/recharge", controlers.UserControl{}.NewRecharge) //生成充值订单
	user.GET("/recharges", controlers.UserControl{}.AllRecharge) //生成充值订单
	//用户取现
	// user.POST("/applycash", controlers.UserControl{}.NewApplyCash) //生成取现订单
	// user.GET("/applycash", controlers.UserControl{}.AllApplyCash)  //生成取现订单

	user.POST("/modifypasswd", controlers.UserControl{}.ModifyPasswd)
	/*收货地址 ----begin*/
	user.POST("/address", controlers.UserControl{}.AddAddress)
	user.PUT("/address", controlers.UserControl{}.UpdateAddress)
	user.DELETE("/address", controlers.UserControl{}.DeleteAddress)
	user.POST("/defaultaddress", controlers.UserControl{}.DefaultAddress)
	user.GET("/address", controlers.UserControl{}.GetAddress)
	/*收货地址 ----end*/

	user.POST("/exparticle", controlers.UserControl{}.NewExposureArticle)
	user.PUT("/exparticle", controlers.UserControl{}.UpdateExposureArticles)
	user.GET("/exparticle", controlers.UserControl{}.GetUserExpArtById)

	user.GET("/myexparticle", controlers.UserControl{}.GetUserExposureArticles)
	user.GET("/article/:id", controlers.UserControl{}.GetUserExposureArticle)

	user.DELETE("/myexparticle", controlers.UserControl{}.DeleteExposureArticles)
	/* 我的物流代购 begin*/
	user.GET("/purchases", purchase.PurchaseControl{}.UserList)
	user.POST("/purchase", purchase.PurchaseControl{}.Add)
	/* 我的物流代购 end*/
	//用户注册
	api.POST("/register", controlers.UserControl{}.Register)
	api.POST("/register/sms", controlers.SendSMS)
	// 找回密码
	api.POST("/phonelogin", controlers.UserControl{}.PhoneLogin)
	api.GET("/articles", controlers.UserControl{}.GetExpArt)
	api.GET("/article/:id", controlers.UserControl{}.GetExpArtById)

	api.GET("/topexparticles", controlers.GetTopExpArts)
	api.GET("/serchexparticles", controlers.SerchExpArts)

	user.POST("/comment", controlers.CommentControl{}.Add)
	api.GET("/newcomments", controlers.CommentControl{}.NewList)
	api.GET("/hotcomments", controlers.CommentControl{}.HotList)

	user.POST("/thumbup", controlers.ThumbsUpControl{}.UP)
	api.GET("/thumbups", controlers.ThumbsUpControl{}.List)
	user.GET("/uptoken", middlewares.GetQnToken)
	user.DELETE("/images", middlewares.DeleteFile)

	user.POST("/tipoffs", controlers.UserControl{}.NewTipOffs)
	user.GET("/msg", controlers.MessageControl{}.GetUserMessage)
	user.DELETE("/msg", controlers.MessageControl{}.Remove)
	user.PUT("/msg", controlers.MessageControl{}.Update)
	user.GET("/msg/news", controlers.MessageControl{}.GetNewMessageCount)

	//物流代购

	// purch := user.Group("/purch")

	api.GET("/purchases", purchase.PurchaseControl{}.List)
	api.GET("/purchase", purchase.PurchaseControl{}.Get)
	user.PUT("/purchase", purchase.PurchaseControl{}.Update)
	user.POST("/invitation", purchase.PurchaseControl{}.Invitation)
	user.DELETE("/purchase", purchase.PurchaseControl{}.Remove)

	user.POST("/quotation", purchase.QuotationOrderControl{}.NewQuotationOrder)
	user.PUT("/quotation", purchase.QuotationOrderControl{}.UpdateQuotationOrder)
	user.GET("/quotations", purchase.QuotationOrderControl{}.UserQuotation)
	user.POST("/refusequotation", purchase.QuotationOrderControl{}.RefuseQuotationOrder)
	/* 我的旅程 begin*/
	user.GET("/journeys", purchase.JourneyControl{}.UserList)
	user.POST("/journey", purchase.JourneyControl{}.Add)
	user.PUT("/journey", purchase.JourneyControl{}.Update)
	user.DELETE("/journey", purchase.JourneyControl{}.Remove)

	api.GET("/destinationuser", purchase.JourneyControl{}.DestinationList)
	api.GET("/destinationpurchase", purchase.PurchaseControl{}.DestinationList)

	/* 我的旅程 end*/
	/* 我的订单 begin*/
	user.GET("/orders", order.OrderControl{}.List)
	user.PUT("/order/:id", order.OrderControl{}.Update)
	pay := api.Group("/pay")
	wx := pay.Group("/wx")
	wx.POST("/notifiy", middlewares.NotifyCallBack)
	wx.Use(authMiddleware.MiddlewareFunc())

	user.GET("/payment/:id", order.PaymentControl{}.CheckPay)
	wx.POST("/get/:id", middlewares.WxPay)

	/* 我的订单 end*/
	// user.Any("/text", func(c *gin.Context) error {
	// 	user := c.Get("user").(*jwt.Token)
	// 	claims := user.Claims.(*util.JwtCustomClaims)
	// 	return c.String(http.StatusOK, "Welcome "+claims.Id+"!")
	// })

}
