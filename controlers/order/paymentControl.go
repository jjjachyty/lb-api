package order

import (
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

type PaymentControl struct {
}

func (PaymentControl) CheckPay(c *gin.Context) {
	var err error
	var payment = new(order.Payment)
	var order = new(order.Order)
	var orderid = c.Param("id")
	if "" != orderid {
		payment.One(bson.M{"order": orderid})
		if payment.State == "0" {
			err = &util.GError{Code: -1, Err: "该单还未支付,请即时支付"}
		} else {
			err = models.One("order", bson.ObjectIdHex(payment.Order), order)
			if nil == err {
				if "0" == order.State {
					//err = &util.GError{Code: -1, Err: "该单已支付，系统订单未更新,请稍后再试一试"}
					err = models.Update("order", bson.M{"_id": bson.ObjectIdHex(payment.Order)}, bson.M{"$set": bson.M{"state": "1"}})
					if nil != err {
						err = &util.GError{Code: -1, Err: "更新订单失败,请稍后再试,或联系工作人员处理"}
					}
				}
			}
		}
	} else {
		err = &util.GError{Code: -1, Err: "订单号不能为空"}
	}

	util.JSON(c, util.ResponseMesage{Message: "检查支付状态", Data: payment, Error: err})

}
