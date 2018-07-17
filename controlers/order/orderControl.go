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

	util.JSON(c, util.ResponseMesage{Message: "获取订单列表", Data: result, Error: err})

}

func (OrderControl) Update(c *gin.Context) {
	var err error
	var order = new(order.Order)
	var orderID = c.Param("id")
	var update bson.M
	var currentUser = middlewares.GetUserIDFromToken(c)
	if err = c.ShouldBindJSON(order); nil == err && "" != orderID {
		switch order.Type {
		case "1": //代购订单
			switch order.State {
			case "1": //更新购买
				if "" != order.Ticket {
					if currentUser == order.SellBy {
						update = bson.M{"$set": bson.M{"ticket": orderID, "state": "1"}}
					} else {
						util.Glog.Warnf("更新订单-操作人%s-非本人操作", currentUser)
						err = &util.GError{Code: 0, Err: "非法操作已被系统记录"}
					}
				} else {
					err = &util.GError{Code: 0, Err: "上传凭证不能为空"}
				}
			}
		}
		err = models.Update(orderCN, bson.M{"_id": bson.ObjectIdHex(orderID)}, update)
	}
	util.JSON(c, util.ResponseMesage{Message: "更新订单", Data: order, Error: err})

}
