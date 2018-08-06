package order

import (
	"lb-api/controlers/pay/wx"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type OrderControl struct{}

const orderCN = "order"

var l = new(sync.Mutex)

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
					update, err = buyerUpdate(currentUser, orderForm, &dbOrder)
				case dbOrder.Seller.ID: //卖家身份
					update, err = sellerUpdate(currentUser, orderForm, &dbOrder)
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
func buyerUpdate(currentUser string, orderForm *order.Order, dborder *order.Order) (update bson.M, err error) {
	switch orderForm.Type {
	case "1": //代购订单
		switch orderForm.State {
		case "-1": //关闭订单

			if dborder.State == "0" {
				update = bson.M{"$set": bson.M{"state": "-1", "buyer.cancelReason": orderForm.Buyer.CancelReason}}
			} else if dborder.State == "1" { //已付款，取消订单
				if "" != orderForm.Buyer.CancelReason {
					update = bson.M{"$set": bson.M{"state": "-1", "buyer.cancelReason": orderForm.Buyer.CancelReason}}

					//退款流程

					err = refund(orderForm, dborder)
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
			//评价买家
			var state = "4"
			if "" != dborder.Seller.Evaluate { //买家也评论了
				state = "43"
			} else { //买家未评论
				state = "41"
			}
			update = bson.M{"$set": bson.M{"buyer.evaluate": orderForm.Buyer.Evaluate, "buyer.evaluateRate": orderForm.Buyer.EvaluateRate, "state": state}}
		case "50": //退换款
			if "" != orderForm.Buyer.ReturnTicket {
				update = bson.M{"$set": bson.M{"state": "50", "buyer.returnTicket": orderForm.Buyer.ReturnTicket, "buyer.returnReason": orderForm.Buyer.ReturnReason}}
			} else {
				err = &util.GError{Code: -1, Err: "凭证不能为空"}

			}
		case "51": //申请退货
			if "" != orderForm.Buyer.ReturnTicket {

				update = bson.M{"$set": bson.M{"state": "51", "buyer.returnTicket": orderForm.Buyer.ReturnTicket, "buyer.returnReason": orderForm.Buyer.ReturnReason}}
			} else {
				err = &util.GError{Code: -1, Err: "凭证不能为空"}

			}
		}
	}
	return update, err
}

//buyerUpdate func买家身份更新
func sellerUpdate(currentUser string, orderForm *order.Order, dborder *order.Order) (update bson.M, err error) {
	switch orderForm.Type {
	case "1": //代购订单
		switch orderForm.State {
		case "-1": //关闭订单

			if dborder.State == "0" {
				update = bson.M{"$set": bson.M{"state": "-1", "seller.cancelReason": "[卖家关闭订单]", "buyer.cancelReason": ""}}
			} else if dborder.State == "1" { //已付款，取消订单
				if "" != orderForm.Seller.CancelReason {
					update = bson.M{"$set": bson.M{"state": "-10", "seller.cancelReason": orderForm.Seller.CancelReason}}
					//退款流程
					err = refund(orderForm, dborder)
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
				if currentUser == dborder.Seller.ID {
					update = bson.M{"$set": bson.M{"seller.buyTicket": orderForm.Seller.BuyTicket, "state": "2"}}
				} else {
					util.Glog.Errorf("更新订单-操作人%s,所属人%s-非本人操作", currentUser, dborder.Seller.ID)
					err = &util.GError{Code: 0, Err: "非法操作已被系统记录"}
				}
			} else {
				err = &util.GError{Code: 0, Err: "上传凭证不能为空"}
			}
		case "2": //代发货
			update = bson.M{"$set": bson.M{"state": "3", "buyer.express.name": orderForm.Buyer.Express.Name, "buyer.express.number": orderForm.Buyer.Express.Number, "buyer.express.createAt": time.Now(), "buyer.express.state": "0"}}
		case "3": //待收货
			//确认收货

		case "4": //已完成
			var state = "4"
			//评价买家
			if "" != dborder.Buyer.Evaluate { //买家也评论了
				state = "43"
			} else { //买家未评论
				state = "42"
			}
			update = bson.M{"$set": bson.M{"seller.evaluate": orderForm.Seller.Evaluate, "seller.evaluateRate": orderForm.Seller.EvaluateRate, "state": state}}
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

func refund(orderForm *order.Order, dborder *order.Order) error {
	var err error
	//退款给买家
	//查询支付信息
	var payPayment = new(order.Payment)
	var returnPayment = new(order.Payment)
	var cancelReason = ""
	if cancelReason = orderForm.Buyer.CancelReason; "" == cancelReason {
		cancelReason = orderForm.Seller.CancelReason
	}
	models.FindOne("payment", bson.M{"order": dborder.ID.Hex(), "payType": "pay", "state": "1"}, payPayment)
	if payPayment.ID.Valid() { //该订单已支付过
		models.FindOne("payment", bson.M{"order": dborder.ID.Hex(), "payType": "return"}, returnPayment)
		if !returnPayment.ID.Valid() { //未收到退款记录
			//新增退款记录
			payment := order.Payment{ID: bson.NewObjectId(), OutTradeNo: payPayment.OutTradeNo, Order: payPayment.Order, PayType: "return", CreateAt: time.Now(), TradeAmount: payPayment.TradeAmount, PayAmount: payPayment.PayAmount, State: "-1"}
			err = payment.Insert()
			if nil == err { //请求支付微信退款
				go func() {
					totaFee, err1 := strconv.ParseInt(strconv.FormatFloat(payment.TradeAmount*100, 'f', 0, 64), 10, 64)
					returnFee, err1 := strconv.ParseInt(strconv.FormatFloat(payment.PayAmount*100, 'f', 0, 64), 10, 64)
					if nil == err1 {
						//申请微信退款
						err = wx.WxRefundControl{}.Refund(payPayment.OutTradeNo, payment.Order, totaFee, returnFee, cancelReason)
						if nil == err {
							util.Glog.Debugf("微信退款申请-申请微信退款成功-订单号%s", orderForm.ID.Hex())
						} else {
							util.Glog.Errorf("微信退款申请-申请微信退款失败[%s]-订单号%s-支付记录%v", err.Error(), orderForm.ID.Hex(), payment)
						}
					} else {
						// err = &util.GError{Code: -1, Err: fmt.Sprintf("微信退款申请-退款金额错误-订单金额%d-退还金额%d", payment.TradeAmount, payment.PayAmount)}
						util.Glog.Errorf("微信退款申请-退款金额错误-订单金额%f-退还金额%f", payment.TradeAmount, payment.PayAmount)
					}

				}()
			} else {
				util.Glog.Errorf("微信退款申请-申请退款-失败-订单号%s", orderForm.ID.Hex())

			}
		} else { //已收到退款申请
			err = &util.GError{Code: -1, Err: "该订单已在退款中,请勿重复退款"}

		}

	} else {
		util.Glog.Errorf("微信退款申请-未找到该订单的支付记录-订单号%s", orderForm.ID.Hex())
	}
	return err
}
func (OrderControl) CheckPay(c *gin.Context) {
	var err error
	var order = new(order.Order)
	var originalID = c.Param("id")
	var payState = false
	if "" != originalID {
		models.FindOne("order", bson.M{"originalID": originalID}, order)
		if order.ID.Valid() {
			if "0" < order.State {
				payState = true
			}
		}
	} else {
		err = &util.GError{Code: -1, Err: "原始单号不能为空"}
	}

	util.JSON(c, util.ResponseMesage{Message: "检查订单支付状态", Data: payState, Error: err})

}

//确认收货
func (OrderControl) Received(c *gin.Context) {
	var err error
	var orderObj = new(order.Order)
	var orderID = c.Param("id")
	var transaction = new(order.Transaction)
	var seller = new(models.User)
	var buyer = new(models.User)
	var userID = middlewares.GetUserIDFromToken(c)
	if "" != orderID {
		models.FindOne("order", bson.M{"_id": bson.ObjectIdHex(orderID)}, orderObj)
		if orderObj.ID.Valid() {
			if orderObj.Buyer.ID == userID {

				err = models.Update(orderCN, bson.M{"_id": bson.ObjectIdHex(orderID)}, bson.M{"$set": bson.M{"state": "4", "buyer.express.state": "1"}})
				if nil == err {
					//转账给卖家
					l.Lock()
					defer l.Unlock()
					err = models.FindOne("user", bson.M{"_id": bson.ObjectIdHex(orderObj.Seller.ID)}, seller)
					err = models.FindOne("user", bson.M{"_id": bson.ObjectIdHex(orderObj.Buyer.ID)}, buyer)

					transaction = &order.Transaction{ID: bson.NewObjectId(), OrderID: orderID, Seller: orderObj.Seller.ID, SellerPreAmount: seller.Wallet.TotalAmount, Buyer: orderObj.Buyer.ID, BuyerPreAmount: buyer.Wallet.TotalAmount, CreateAt: time.Now(), State: "1", Amount: orderObj.StrikePrice}

					err = models.Insert("transaction", transaction)

					if nil == err {
						//更新卖家金额
						err = models.Update("user", bson.M{"_id": bson.ObjectIdHex(orderObj.Seller.ID)}, bson.M{"$set": bson.M{"wallet.totalAmount": transaction.SellerPreAmount + transaction.Amount}})
						// err = models.Update("user", bson.M{"_id": bson.ObjectIdHex(orderObj.Buyer.ID)}, bson.M{"$set": bson.M{"wallet.totalAmount": transaction.BuyerPreAmount + transaction.Amount}})
						if nil != err {
							util.Glog.Errorf("人工补偿-转账失败，往卖家%s账户新增金额错误-错误信息%s-账户之前金额%f-转账金额%f", orderObj.Seller.ID, err, transaction.SellerPreAmount, transaction.Amount)
						}
					}
				}
			} else {
				util.Glog.Errorf("确认收货-用户%s非法操作用户%s的订单%sIP%s", userID, orderObj.Buyer.ID, orderID, c.ClientIP())
				err = &util.GError{Code: -1, Err: "非法操作他人订单,已被系统记录"}
			}
		} else {
			err = &util.GError{Code: -1, Err: "订单号不存在"}
		}
	} else {
		err = &util.GError{Code: -1, Err: "订单号不能为空"}
	}
	util.JSON(c, util.ResponseMesage{Message: "确认收货", Data: nil, Error: err})

}
