package router

import (
	"fmt"
	"lb-api/controlers"
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
	api.POST("/login", middlewares.JWT().LoginHandler)
	api.POST("/refresh_token", middlewares.JWT().RefreshHandler)
	// api.POST("/emailvalid", controlers.ValidEmail)

	user := api.Group("/user")

	user.Use(middlewares.JWT().MiddlewareFunc())
	//用户实名认证
	user.POST("/scanidcard", controlers.UserControl{}.IDCardOCR)
	user.POST("/valididcard", controlers.UserControl{}.ValidIDCard)
	//用户设置
	user.GET("/info", controlers.UserControl{}.GetUserInfo)
	user.PUT("/info", controlers.UserControl{}.Update)
	//用户充值
	user.POST("/recharge", controlers.UserControl{}.NewRecharge) //生成充值订单
	user.GET("/recharges", controlers.UserControl{}.AllRecharge) //生成充值订单
	//用户取现
	user.POST("/applycash", controlers.UserControl{}.NewApplyCash) //生成取现订单
	user.GET("/applycash", controlers.UserControl{}.AllApplyCash)  //生成取现订单

	user.POST("/modifypasswd", controlers.UserControl{}.ModifyPasswd)

	user.POST("/exparticle", controlers.UserControl{}.NewExposureArticle)
	user.PUT("/exparticle", controlers.UserControl{}.UpdateExposureArticles)
	user.GET("/exparticle", controlers.UserControl{}.GetUserExpArtById)

	user.GET("/myexparticle", controlers.UserControl{}.GetUserExposureArticles)
	user.GET("/article/:id", controlers.UserControl{}.GetUserExposureArticle)

	user.DELETE("/myexparticle", controlers.UserControl{}.DeleteExposureArticles)
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
	user.GET("/thumbups", controlers.ThumbsUpControl{}.List)
	user.GET("/uptoken", middlewares.GetQnToken)
	user.DELETE("/images", middlewares.DeleteFile)

	user.POST("/tipoffs", controlers.UserControl{}.NewTipOffs)
	user.GET("/msg", controlers.MessageControl{}.GetUserMessage)
	user.DELETE("/msg", controlers.MessageControl{}.Remove)
	user.PUT("/msg", controlers.MessageControl{}.Update)
	user.GET("/msg/news", controlers.MessageControl{}.GetNewMessageCount)

	//物流代购

	purch := api.Group("/purch")

	purch.GET("/list", purchase.PurchaseControl{}.List)
	purch.GET("/get", purchase.PurchaseControl{}.Get)

	purch.Use(middlewares.JWT().MiddlewareFunc())
	purch.POST("/quotation", purchase.QuotationOrderControl{}.NewQuotationOrder)
	// user.Any("/text", func(c *gin.Context) error {
	// 	user := c.Get("user").(*jwt.Token)
	// 	claims := user.Claims.(*util.JwtCustomClaims)
	// 	return c.String(http.StatusOK, "Welcome "+claims.Id+"!")
	// })

}
