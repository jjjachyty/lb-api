package purchase

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models/purchase"
	"lb-api/util"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

type QuotationOrderControl struct{}

func (QuotationOrderControl) NewQuotationOrder(c *gin.Context) {
	var err error
	var expiryTime time.Time
	qo := new(purchase.QuotationOrder)
	//处理失效时间
	// expiryTime, err = getTime(c.PostForm("expiryTime"))
	// if nil == err {

	if err = c.ShouldBindJSON(qo); nil == err {
		//处理总金额
		for _, p := range qo.Products {
			qo.Amount += p.Price
		}
		fmt.Println("middlewares.GetUserIDFromToken(c)", middlewares.GetUserIDFromToken(c))
		jwtData := middlewares.GetPalyloadFromToken(c)
		qo.BuyByID = jwtData["id"].(string)
		qo.State = "1"
		qo.ID = bson.NewObjectId()
		qo.BuyByName = jwtData["nickName"].(string)
		qo.CreateAt = time.Now().In(util.GetLocation())
		qo.Amount += qo.Charge
		fmt.Println("qo.CreateAt", qo.CreateAt)
		// qo.ExpiryTime = expiryTime
		err = qo.Insert()
	}
	// }
	fmt.Println("NewQuotationOrder", err, expiryTime)
	util.JSON(c, util.ResponseMesage{Message: "新增报价单", Data: nil, Error: err})

}
