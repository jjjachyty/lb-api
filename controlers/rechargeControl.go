package controlers

import (
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/util"
	"time"

	"github.com/gin-gonic/gin"

	"labix.org/v2/mgo/bson"
)

func (UserControl) NewRecharge(c *gin.Context) {
	var rech = new(models.Recharge)
	var err error
	if err = c.Bind(rech); nil == err {
		if 0 < rech.Amount && 0 == int(rech.Amount)%100 {
			rech.TradeNumber = string(util.RandNumber(4))
			rech.State = "0"
			rech.CreateAt = time.Now()
			rech.UserID = middlewares.GetUserIDFromToken(c)
			err = rech.Insert()
		} else {
			err = &util.GError{Code: 4002, Err: "保证金金额必须大于0且为100的整数倍"}
		}
	}
	util.JSON(c, util.ResponseMesage{Message: "生成充值订单", Data: rech, Error: err})
}
func (UserControl) AllRecharge(c *gin.Context) {
	var rech = &models.Recharge{UserID: middlewares.GetUserIDFromToken(c)}
	recharges, err := rech.FindAllByCondition(bson.M{"userID": rech.UserID})

	util.JSON(c, util.ResponseMesage{Message: "获取我的充值订单", Data: recharges, Error: err})
}
