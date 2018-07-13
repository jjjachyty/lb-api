package order

import (
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type OrderControl struct{}

const orderCN = "order"

func (OrderControl) List(c *gin.Context) {
	var err error
	var cond bson.M
	orderType := c.Query("type")
	state := c.Query("state")
	identity := c.Query("identity")
	result := make([]order.Order, 0)
	if "" != orderType && "" != state {
		userid := middlewares.GetUserIDFromToken(c)
		cond = bson.M{"type": orderType, "state": state}
		if identity == "0" { // 我买的
			cond["buyBy"] = userid
		} else { //我卖的
			cond["sellBy"] = userid
		}
		err = models.Find(orderCN, &result, "-createAt", 10, bson.M{}, cond)
	} else {
		err = &util.GError{Code: 0, Err: "数据完整性错误"}
	}

	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}
