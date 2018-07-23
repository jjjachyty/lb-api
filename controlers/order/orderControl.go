package order

import (
	"fmt"
	"lb-api/controlers/pay/wx"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"
	"strconv"
	"time"

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
		cond = bson.M{"type": orderType, "state": bson.M{"$regex": state}}
		if identity == "0" { // 我买的
			cond["buyer.id"] = userid
		} else { //我卖的
			cond["seller.id"] = userid
		}
		util.Glog.Debugf("查询订单列表-查询条件%v", cond)
		err = models.Find(orderCN, &result, "-createAt", 10, bson.M{}, cond)
	} else {
		err = &util.GError{Code: 0, Err: "数据完整性错误"}
	}

	util.JSON(c, util.ResponseMesage{Message: "获取订单列表", Data: result, Error: err})

}

func (OrderControl) Update(c *gin.Context) {
	var err error
	var orderForm = new(order.Order)
	var dbOrder order.Order
	var orders []order.Order
	var orderID = c.Param("id")
	var update bson.M
	var currentUser = middlewares.GetUserIDFromToken(c)
	if err = c.ShouldBindJSON(orderForm); nil == err {
		if "" != orderID {
			//查询订单

			err = models.Find(orderCN, &orders, "-createAt", 10, bson.M{}, bson.M{"_id": bson.ObjectIdHex(orderID)})
			if len(orders) == 1 {
				dbOrder = orders[0]
				orderForm.Type = dbOrder.Type
				switch currentUser { // 当前操作用户
				case dbOrder.Buyer.ID: //买家身份
					update, err = buyerUpdate(currentUser, orderForm, dbOrder)
				case dbOrder.Seller.ID: //卖家身份
					update, err = sellerUpdate(currentUser, orderForm, dbOrder)
				default:
					util.Glog.Warnf("非法更新订单-操作人%s-,IP%s", currentUser, c.ClientIP())
					err = &util.GError{Code: -1, Err: "非法操作,不能操作他人订单,系统已记录"}
				}

			} else {
				err = &util.GError{Code: -1, Err: "订单不存在"}
			}
		} else {
			err = &util.GError{Code: -1, Err: "订单ID不能为空"}
		}
	}
	if nil == err {
		err = models.Update(orderCN, bson.M{"_id": bson.ObjectIdHex(orderID)}, update)

	}

	util.JSON(c, util.ResponseMesage{Message: "更新订单", Data: dbOrder, Error: err})

}

//buyerUpdate func买家身份更新
func buyerUpdate(currentUser string, orderForm *order.Order, dborder order.Order) (update bson.M, err error) {
	switch orderForm.Type {
	case "1": //代购订单
		switch orderForm.State {
		case "-1": //关闭订单

			if dborder.State == "0" {
				update = bson.M{"$set": bson.M{"state": "-1", "buyer.cancelReason": orderForm.Buyer.CancelReason}}
			} else if dborder.State == "1" { //已付款，取消订单
				if "" != orderForm.Buyer.CancelReason {
					update = bson.M{"$set": bson.M{"state": "-1", "buyer.cancelReason": orderForm.Buyer.CancelReason}}
				} else {
					err = &util.GError{Code: -1, Err: "取消原因不能为空"}
				}
			} else {
				err = &util.GError{Code: -1, Err: "只能关闭[待付款]和[待购买]的订单"}
			}
		case "0": //待付款
			if "" != orderForm.Buyer.CancelReason { //取消订单
				update = bson.M{"$set": bson.M{"buyer.cancelReason": orderForm.Buyer.CancelReason, "state": "-1"}}
			} else {
				err = &util.GError{Code: -1, Err: "取消订单原因不能为空"}
			}
		case "1": //更新购买
			// if "" != orderForm.BuyTicket {
			// 	if currentUser == orderForm.Seller.ID {
			// 		update = bson.M{"$set": bson.M{"ticket": orderForm.ID.Hex(), "state": "1"}}
			// 	} else {
			// 		util.Glog.Warnf("更新订单-操作人%s-非本人操作", currentUser)
			// 		err = &util.GError{Code: 0, Err: "非法操作已被系统记录"}
			// 	}
			// } else {
			// 	err = &util.GError{Code: 0, Err: "上传凭证不能为空"}
			// }
		case "2": //代发货
		case "3": //待收货
			//确认收货
			update = bson.M{"$set": bson.M{"state": "4"}}
		case "4": //已完成
			//评价卖家
			update = bson.M{"$set": bson.M{"buyer.reviews": orderForm.Buyer.Reviews}}
		case "50": //退换款
			update = bson.M{"$set": bson.M{"state": "50"}}
		case "51": //申请退货
			update = bson.M{"$set": bson.M{"state": "51"}}
		}
	}
	return update, err
}

//buyerUpdate func买家身份更新
func sellerUpdate(currentUser string, orderForm *order.Order, dborder order.Order) (update bson.M, err error) {
	switch orderForm.Type {
	case "1": //代购订单
		switch orderForm.State {
		case "-1": //关闭订单

			if dborder.State == "0" {
				update = bson.M{"$set": bson.M{"state": "-1", "seller.cancelReason": "[卖家关闭订单]"}}
			} else if dborder.State == "1" { //已付款，取消订单
				if "" != orderForm.Seller.CancelReason {
					update = bson.M{"$set": bson.M{"state": "-1", "seller.cancelReason": orderForm.Seller.CancelReason}}
					//退款给买家
					// go func() {
					//查询支付信息
					var payment = new(order.Payment)
					err := models.FindOne("payment", bson.M{"order": dborder.ID.Hex(), "payType": "pay", "state": "1"}, payment)
					fmt.Println("order", dborder.ID.Hex(), payment, err)
					if payment.ID.Valid() {
						totaFee, err := strconv.ParseInt(strconv.FormatFloat(payment.TradeAmount*100, 'f', 0, 64), 10, 64)
						returnFee, err := strconv.ParseInt(strconv.FormatFloat(payment.PayAmount*100, 'f', 0, 64), 10, 64)
						if nil == err {
							//申请微信退款
							err := wx.WxRefundControl{}.Refund(payment.OutTradeNo, payment.Order, totaFee, returnFee, orderForm.Seller.CancelReason)
							if nil == err {
								util.Glog.Debugf("微信退款申请-申请微信退款成功-订单号%s", orderForm.ID.Hex())
							} else {
								util.Glog.Errorf("微信退款申请-申请微信退款失败-订单号%s-支付记录%v", orderForm.ID.Hex(), payment)
							}
						} else {
							// err = &util.GError{Code: -1, Err: fmt.Sprintf("微信退款申请-退款金额错误-订单金额%d-退还金额%d", payment.TradeAmount, payment.PayAmount)}
							util.Glog.Errorf("微信退款申请-退款金额错误-订单金额%f-退还金额%f", payment.TradeAmount, payment.PayAmount)
						}
					} else {
						util.Glog.Errorf("微信退款申请-未找到该订单的支付记录-订单号%s", orderForm.ID.Hex())
					}

					// }()

				} else {
					err = &util.GError{Code: -1, Err: "取消原因不能为空"}
				}
			} else {
				err = &util.GError{Code: -1, Err: "只能关闭[待付款]和[待购买]的订单"}
			}
		case "0": //修改价格
			if dborder.State == "0" {
				if orderForm.Charge >= 0 { //修改价格
					productAMount := dborder.StrikePrice - dborder.Charge
					update = bson.M{"$set": bson.M{"charge": orderForm.Charge, "strikePrice": productAMount + orderForm.Charge}}
				} else {
					err = &util.GError{Code: -1, Err: "代购费需大于等于0"}
				}
			} else {
				err = &util.GError{Code: -1, Err: "只能修改[待付款]的价格"}
			}
		case "1": //更新购买
			if "" != orderForm.Seller.BuyTicket {
				if currentUser == orderForm.Seller.ID {
					update = bson.M{"$set": bson.M{"ticket": orderForm.ID.Hex(), "state": "1"}}
				} else {
					util.Glog.Warnf("更新订单-操作人%s,所属人%s-非本人操作", currentUser, dborder.Seller.ID)
					err = &util.GError{Code: 0, Err: "非法操作已被系统记录"}
				}
			} else {
				err = &util.GError{Code: 0, Err: "上传凭证不能为空"}
			}
		case "2": //代发货
			update = bson.M{"$set": bson.M{"state": "3", "express.name": orderForm.Buyer.Express.Name, "express.number": orderForm.Buyer.Express.Number, "express.createAt": time.Now(), "express.state": "已寄出"}}
		case "3": //待收货
			//确认收货

		case "4": //已完成
			//评价买家
			update = bson.M{"$set": bson.M{"seller.reviews": orderForm.Seller.Reviews}}
		case "500": //拒绝退换货
			update = bson.M{"$set": bson.M{"state": "500"}}
		case "501": //退款
			update = bson.M{"$set": bson.M{"state": "500"}}
		case "510": //确认退货
			update = bson.M{"$set": bson.M{"state": "510", "seller.returnAddress": orderForm.Seller.ReturnAddress}}
		default:
			err = &util.GError{Code: -1, Err: "未找到对应的订单操作"}
		}
	}
	return update, err
}
