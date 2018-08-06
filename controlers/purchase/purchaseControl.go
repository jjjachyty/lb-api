package purchase

import (
	"fmt"
	"lb-api/controlers/pay/wx"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/models/purchase"
	"lb-api/util"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

type PurchaseControl struct{}

func (PurchaseControl) List(c *gin.Context) {
	sort := []string{"state", "-updateAt"}
	keyWords := c.Query("keyWords")
	var cond bson.M
	var appCond bson.M
	if "" != keyWords {
		appCond = bson.M{"$or": []bson.M{bson.M{"destination": bson.M{"$regex": keyWords}}, bson.M{"content": bson.M{"$regex": keyWords}}, bson.M{"products.name": bson.M{"$regex": keyWords}}, bson.M{"products.describe": bson.M{"$regex": keyWords}}, bson.M{"location": bson.M{"$regex": keyWords}}}}

	}
	cond = bson.M{"$and": []bson.M{bson.M{"$or": []bson.M{bson.M{"state": "0"}, bson.M{"state": "1"}}}, appCond}}

	result, err := purchase.Purchase{}.Find(sort, 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}

func (PurchaseControl) DestinationList(c *gin.Context) {
	sort := c.DefaultQuery("sort", "endDate")
	destination := c.Query("destination")
	id := c.Query("id")
	user := c.Query("user")
	var cond bson.M
	if "" != destination {
		cond = bson.M{"destination": destination, "state": "1"}
		if "" != user {
			cond["createBy"] = bson.M{"$ne": user}
		}
		if "" != user {
			cond["_id"] = bson.M{"$ne": bson.ObjectIdHex(id)}
		}
	}
	result, err := purchase.Purchase{}.Find([]string{sort}, 10, bson.M{"createBy": 1, "amount": 1, "products.name": 1}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取代购推荐", Data: result, Error: err})
}

// Invitation  func 邀请代购报价
func (PurchaseControl) Invitation(c *gin.Context) {
	var err error
	purchaseID := c.PostForm("purchaseID")
	beInviter := c.PostForm("beInviter")
	destination := c.PostForm("destination")
	inviter := c.PostForm("inviter")
	if "" != purchaseID && "" != destination && "" != inviter && "" != beInviter {
		if beInviter == middlewares.GetUserIDFromToken(c) {
			err = models.Message{ID: bson.NewObjectId(), Type: "代购邀请", From: inviter, To: beInviter, Content: fmt.Sprintf("您有一个%s<a href='#/purchase/%s'>代购</a>邀请", destination, purchaseID), CreateAt: time.Now(), State: "1"}.Insert()
			if nil == err {
				err = purchase.Purchase{}.Update(bson.M{"_id": bson.ObjectIdHex(purchaseID)}, bson.M{"$push": bson.M{"inviters": inviter}})
			}
		} else {
			err = &util.GError{Code: 0, Err: "非法操作他人代购单"}
		}

	} else {
		err = &util.GError{Code: 0, Err: "数据完整性错误"}
	}
	util.JSON(c, util.ResponseMesage{Message: "邀请代购", Data: nil, Error: err})

}

func (PurchaseControl) UserList(c *gin.Context) {
	var cond bson.M
	userID := middlewares.GetUserIDFromToken(c)
	cond = bson.M{"createBy": userID}
	result, err := purchase.Purchase{}.Find([]string{"-createAt"}, 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取我的代购单列表", Data: result, Error: err})

}

func (PurchaseControl) Get(c *gin.Context) {
	id := c.Query("id")
	var cond bson.M
	var result purchase.Purchase
	var results []purchase.Purchase

	// var qos []purchase.QuotationOrder
	var err error
	if "" != id {
		cond = bson.M{"_id": bson.ObjectIdHex(id)}

		results, err = purchase.Purchase{}.Find([]string{"_id"}, 10, bson.M{}, cond)
		if len(results) > 0 {
			result = results[0]

		}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}

func (PurchaseControl) Add(c *gin.Context) {
	var err error
	var purchase = new(purchase.Purchase)
	if err = c.ShouldBindJSON(purchase); nil == err {
		purchase.ID = bson.NewObjectId()
		purchase.CreateAt = time.Now()
		purchase.UpdateAt = purchase.CreateAt
		purchase.State = "0"
		purchase.CreateBy = middlewares.GetUserIDFromToken(c)
		err = purchase.Insert()
	}
	fmt.Println("err", err)
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: purchase, Error: err})

}

func (PurchaseControl) Update(c *gin.Context) {
	var err error
	var purchases []purchase.Purchase
	var purchaseObj = new(purchase.Purchase)

	if err = c.ShouldBindJSON(purchaseObj); nil == err {
		purchases, err = purchase.Purchase{}.Find([]string{"_id"}, 1, bson.M{}, bson.M{"_id": purchaseObj.ID})
		if len(purchases) == 1 {
			dbPurchase := purchases[0]
			if dbPurchase.CreateBy == middlewares.GetUserIDFromToken(c) {
				if "0" == dbPurchase.State || "-1" == dbPurchase.State {
					var amount = 0.0
					for _, p := range purchaseObj.Products {
						quantity, _ := strconv.ParseFloat(strconv.FormatInt(p.Quantity, 20), 64)
						amount += p.Price * quantity
					}
					err = purchase.Purchase{}.Update(bson.M{"_id": purchaseObj.ID}, bson.M{"$set": bson.M{"state": "0", "amount": amount, "destination": purchaseObj.Destination, "address": purchaseObj.Address, "content": purchaseObj.Content, "products": purchaseObj.Products, "updateAt": time.Now()}})
				} else {
					err = &util.GError{Code: 0, Err: "只能更新状态为[待报价][逾期未处理]的订购单"}
				}
			} else {
				err = &util.GError{Code: 0, Err: "不能操作他人代购单"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "该报价单不存在"}
		}

	}
	fmt.Println("更新报价单", err)
	util.JSON(c, util.ResponseMesage{Message: "更新物流代购", Data: purchaseObj, Error: err})

}

// Delete func 删除
func (PurchaseControl) Remove(c *gin.Context) {
	var id = c.Query("id")
	var err error
	var purchases []purchase.Purchase
	if bson.IsObjectIdHex(id) {
		purchases, err = purchase.Purchase{}.Find([]string{"_id"}, 1, bson.M{}, bson.M{"_id": bson.ObjectIdHex(id)})
		if len(purchases) == 1 {
			dbPurchase := purchases[0]
			if "0" == dbPurchase.State { //待报价状态可删除

				if dbPurchase.CreateBy == middlewares.GetUserIDFromToken(c) {
					err = models.Remove("purchase", bson.M{"_id": bson.ObjectIdHex(id)})
					util.Glog.Debugf("删除代购单-操作人%s-原数据%v-状态%v", dbPurchase.CreateBy, dbPurchase, err)
					if nil == err { //删除成功后更新报价单
						go purchase.QuotationOrder{}.Update(bson.M{"purchaseID": id}, bson.M{"state": "-1", "refuseReason": "该报价单已被删除"})
						//删除7牛云的图片
						go func(purchase.Purchase) {
							var keys = make([]string, 0)
							for _, pd := range dbPurchase.Products {
								keys = append(keys, pd.Images)
							}
							middlewares.DeleteFiles("4t-purchase", keys...)
						}(dbPurchase)
					}
				} else {
					err = &util.GError{Code: 0, Err: "不能操作他人代购单"}
				}
			}
		} else {
			err = &util.GError{Code: 0, Err: "该代购单不存在"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "删除我的代购单", Data: nil, Error: err})
}

//选中代购单
func (PurchaseControl) Confirm(c *gin.Context) {
	var err error
	var qo = new(purchase.QuotationOrder)
	var purch = new(purchase.Purchase)
	var userID = middlewares.GetUserIDFromToken(c)
	var returnMsg wx.ReturnMsg
	var orderObj = new(order.Order)
	quotationid := c.PostForm("quotationid")
	if bson.IsObjectIdHex(quotationid) {
		models.One("quotation", bson.ObjectIdHex(quotationid), qo)
		if qo.ID.Valid() { //查到报价单
			//处理总费用
			var totalAmount = 0.0
			for _, pd := range qo.Products {
				qt, _ := strconv.ParseFloat(strconv.FormatInt(pd.Quantity, 10), 64)
				totalAmount += (pd.Price * qt)
			}
			totalAmount += qo.Charge
			//查询代购单
			models.FindOne("purchase", bson.M{"_id": bson.ObjectIdHex(qo.PurchaseID)}, purch)
			if purch.ID.Valid() {
				if purch.CreateBy == userID { //确保是本人操作
					//更新代购单
					err = models.Update("purchase", bson.M{"_id": purch.ID}, bson.M{"$set": bson.M{"state": "2", "quotationID": quotationid}})
					err = models.Update("quotation", bson.M{"_id": qo.ID}, bson.M{"$set": bson.M{"state": "2"}})
					if nil == err { //更新代购成功后，新增定价单
						//判断是否已经存在定价单
						models.FindOne("order", bson.M{"originalID": purch.ID.Hex()}, orderObj)
						if !orderObj.ID.Valid() { //不存在则新增
							var buyer = order.Buyer{ID: purch.CreateBy, Name: purch.Creator, IP: c.ClientIP(), Express: order.Express{ReceivingAddress: purch.Address}}
							var seller = order.Seller{ID: qo.CreateBy, Name: qo.Creator}
							orderObj = &order.Order{Buyer: buyer, Seller: seller, OriginalID: purch.ID.Hex(), ID: bson.NewObjectId(), StrikePrice: totalAmount, Charge: qo.Charge, OriginalLink: "/purchase/" + qo.PurchaseID, Type: "1", State: "0", Products: qo.Products}
							err = models.Insert("order", orderObj)
						}

						if nil == err { //调用微信支付
							returnMsg, err = wx.WxPayControl{}.GetWxPay(orderObj, userID)
						}
					}
				} else {
					util.Glog.Errorf("非法操作,用户%s操作被操作用户%s的代购单-IP地址", userID, qo.CreateBy, c.ClientIP())
					err = &util.GError{Code: -1, Err: "非法操作,只能操作自己的订单,已被系统记录"}
				}
			}
		}
	}
	util.JSON(c, util.ResponseMesage{Message: "选定代购人,并支付代购单", Data: returnMsg, Error: err})

}
