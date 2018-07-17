package purchase

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models/purchase"
	"lb-api/util"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

type QuotationOrderControl struct{}

func (QuotationOrderControl) UserQuotation(c *gin.Context) {
	var err error
	var qos []purchase.QuotationOrder
	qos, err = purchase.QuotationOrder{}.Find("_id", 10, bson.M{}, bson.M{"createBy": middlewares.GetUserIDFromToken(c)})
	util.JSON(c, util.ResponseMesage{Message: "获取我的报价单", Data: qos, Error: err})

}

func (qoc QuotationOrderControl) NewQuotationOrder(c *gin.Context) {
	var err error
	qo := new(purchase.QuotationOrder)
	//处理失效时间
	// expiryTime, err = getTime(c.PostForm("expiryTime"))
	// if nil == err {

	if err = c.ShouldBindJSON(qo); nil == err {

		//处理总金额
		for _, p := range qo.Products {
			qo.Amount += p.Price
		}

		jwtData := middlewares.GetPalyloadFromToken(c)
		qo.CreateBy = jwtData["id"].(string)
		qo.State = "1" //报价
		qo.ID = bson.NewObjectId()
		qo.Creator = jwtData["nickName"].(string)
		qo.CreateAt = time.Now().In(util.GetLocation())
		qo.Amount += qo.Charge
		fmt.Println("qo.CreateAt", qo.CreateAt)
		// qo.ExpiryTime = expiryTime
		//检查能否新增
		err = qoc.canNew(*qo)
		if nil == err {

			err = qo.Insert()
			if err == nil {
				//更新代购单为报价反馈中
				err = purchase.Purchase{}.Update(bson.M{"_id": bson.ObjectIdHex(qo.PurchaseID)}, bson.M{"$set": bson.M{"state": "1"}})
				if nil != err { //报价失败
					err = purchase.QuotationOrder{ID: qo.ID}.Delete()
				}
				util.Glog.Debugf("删除-报价单%v-状态%v", qo, err)

			}
		}
	}
	// }
	util.JSON(c, util.ResponseMesage{Message: "新增报价单", Data: qo, Error: err})

}

func (qoc QuotationOrderControl) UpdateQuotationOrder(c *gin.Context) {
	var err error
	qo := new(purchase.QuotationOrder)
	//处理失效时间
	// expiryTime, err = getTime(c.PostForm("expiryTime"))
	// if nil == err {

	if err = c.ShouldBindJSON(qo); nil == err {

		jwtData := middlewares.GetPalyloadFromToken(c)

		if qo.CreateBy == jwtData["id"].(string) { //本人操作本人的报价单
			//处理总金额
			qo.Amount = 0
			for _, p := range qo.Products {
				qo.Amount += p.Price
			}
			qo.Amount += qo.Charge

			err = qoc.canUpdate(*qo)
			if nil == err {

				err = qo.Update(bson.M{"_id": qo.ID}, bson.M{"$set": bson.M{"state": "1", "refuseReason": "", "products": qo.Products, "amount": qo.Amount, "charge": qo.Charge, "expiryTime": qo.ExpiryTime, "deliveryTime": qo.DeliveryTime}})
				if err == nil {
					//更新代购单为报价反馈中
					err = purchase.Purchase{}.Update(bson.M{"_id": bson.ObjectIdHex(qo.PurchaseID)}, bson.M{"$set": bson.M{"state": "1"}})
					if nil != err { //报价失败
						err = purchase.QuotationOrder{ID: qo.ID}.Delete()
					}
					util.Glog.Debugf("删除-报价单%v-状态%v", qo, err)

				}
			}
		} else {
			err = &util.GError{Code: 0, Err: "非本人报价单,非法操作"}
		}
	}
	// }
	util.JSON(c, util.ResponseMesage{Message: "更新报价单", Data: nil, Error: err})
}

//RefuseQuotationOrder 拒绝报价单
func (QuotationOrderControl) RefuseQuotationOrder(c *gin.Context) {
	var purchs []purchase.Purchase
	var allowRepeatBool bool
	var err error
	quotationID := c.PostForm("quotationID")
	purchaseID := c.PostForm("purchaseID")
	reasonType := c.PostForm("reasonType")

	refuseReason := c.PostForm("reason")
	allowRepeat := c.PostForm("allowRepeat")
	if "" != quotationID && bson.IsObjectIdHex(purchaseID) && "" != refuseReason && "" != allowRepeat && "" != reasonType {
		//验证是否是本人操作自己的代购单
		purchs, err = purchase.Purchase{}.Find([]string{"_id"}, 0, bson.M{}, bson.M{"_id": bson.ObjectIdHex(purchaseID), "createBy": middlewares.GetUserIDFromToken(c)})
		if len(purchs) == 1 { //是本人操作
			//更新代购单和报价单
			util.Glog.Debugf("代购单%s拒绝报价单%s,拒绝理由%s,运行再次报价%s", purchaseID, quotationID, refuseReason, allowRepeat)
			allowRepeatBool, err = strconv.ParseBool(allowRepeat)
			if nil == err {
				err = purchase.QuotationOrder{}.Update(bson.M{"_id": bson.ObjectIdHex(quotationID)}, bson.M{"$set": bson.M{"state": "0", "reasonType": reasonType, "refuseReason": refuseReason, "allowRepeat": allowRepeatBool}})
				if nil == err {
					err = purchase.Purchase{}.Update(bson.M{"_id": bson.ObjectIdHex(purchaseID)}, bson.M{"$set": bson.M{"state": "0"}})
				}
			}

		} else {
			err = &util.GError{Code: 0, Err: "非本人代购单,非法操作"}
		}
	} else {
		err = &util.GError{Code: 0, Err: "数据完整性错误"}
	}
	util.JSON(c, util.ResponseMesage{Message: "拒绝报价单", Data: nil, Error: err})

}

func (QuotationOrderControl) canNew(qo purchase.QuotationOrder) error {
	var err error
	var purchases []purchase.Purchase
	var quotations []purchase.QuotationOrder
	//检查报价单是否能报价
	//检查是否已经报过价
	purchases, err = purchase.Purchase{}.Find([]string{"_id"}, 0, bson.M{}, bson.M{"_id": bson.ObjectIdHex(qo.PurchaseID)})
	if len(purchases) == 1 {
		if qo.CreateBy != purchases[0].CreateBy {
			if "0" == purchases[0].State { //代价单为待报价状态
				//检查是否重复报价
				fmt.Println("purchaseIDs", qo.PurchaseID, "createBy", qo.CreateBy)
				quotations, err = purchase.QuotationOrder{}.Find("_id", 0, bson.M{}, bson.M{"purchaseID": qo.PurchaseID, "buyByID": qo.CreateBy})
				if len(quotations) > 0 {
					err = &util.GError{Code: 0, Err: "已报价不能重复报价"}
				}

			} else {
				err = &util.GError{Code: 0, Err: "该报价单非[待报价]状态"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "不能自己代购自己的"}
		}
	} else {
		err = &util.GError{Code: 0, Err: "该报价单不存在"}

	}
	return err
}

func (QuotationOrderControl) canUpdate(qo purchase.QuotationOrder) error {
	var err error
	var purchases []purchase.Purchase
	//检查报价单是否能报价
	//检查是否已经报过价
	purchases, err = purchase.Purchase{}.Find([]string{"_id"}, 0, bson.M{}, bson.M{"_id": bson.ObjectIdHex(qo.PurchaseID)})

	if len(purchases) == 1 {
		if qo.CreateBy != purchases[0].CreateBy {
			if "0" != purchases[0].State { //代价单为待报价状态
				err = &util.GError{Code: 0, Err: "该报价单非[待报价]状态"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "不能自己代购自己的"}
		}
	} else {
		err = &util.GError{Code: 0, Err: "该报价单不存在"}

	}
	return err
}
