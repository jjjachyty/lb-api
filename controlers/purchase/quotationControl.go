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

		jwtData := middlewares.GetPalyloadFromToken(c)
		qo.BuyByID = jwtData["id"].(string)
		qo.State = "1" //报价
		qo.ID = bson.NewObjectId()
		qo.BuyByName = jwtData["nickName"].(string)
		qo.CreateAt = time.Now().In(util.GetLocation())
		qo.Amount += qo.Charge
		fmt.Println("qo.CreateAt", qo.CreateAt)
		// qo.ExpiryTime = expiryTime
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
	// }
	fmt.Println("NewQuotationOrder", err, expiryTime)
	util.JSON(c, util.ResponseMesage{Message: "新增报价单", Data: nil, Error: err})

}

func (QuotationOrderControl) UpdateQuotationOrder(c *gin.Context) {

}

func (QuotationOrderControl) RefuseQuotationOrder(c *gin.Context) {
	var purchs []purchase.Purchase
	var err error
	quotationID := c.PostForm("quotationID")
	purchaseID := c.PostForm("purchaseID")
	refuseReason := c.PostForm("reason")
	if "" != quotationID && bson.IsObjectIdHex(purchaseID) && "" != refuseReason {
		//验证是否是本人操作自己的代购单
		purchs, err = purchase.Purchase{}.Find("_id", 0, bson.M{}, bson.M{"_id": bson.ObjectIdHex(purchaseID), "createBy": middlewares.GetUserIDFromToken(c)})
		if len(purchs) == 1 { //是本人操作
			//更新代购单和报价单
			util.Glog.Debugf("代购单%s拒绝报价单%s,拒绝理由%s", purchaseID, quotationID, refuseReason)
			err = purchase.QuotationOrder{}.Update(bson.M{"_id": bson.ObjectIdHex(quotationID)}, bson.M{"$set": bson.M{"state": "0", "refuseReason": refuseReason}})
			if nil == err {
				err = purchase.Purchase{}.Update(bson.M{"_id": bson.ObjectIdHex(purchaseID)}, bson.M{"$set": bson.M{"state": "0"}})
			}
		}
	} else {
		err = &util.GError{Code: 0, Err: "数据完整性错误"}
	}
	util.JSON(c, util.ResponseMesage{Message: "拒绝报价单", Data: nil, Error: err})

}

func (QuotationOrderControl) canNew(qo purchase.QuotationOrder) {
	var err error
	var purchases []purchase.QuotationOrder
	var quotations []purchase.QuotationOrder
	//检查报价单是否能报价
	//检查是否已经报过价
	purchases, err = purchase.Purchase{}.Find("_id", 0, bson.M{}, bson.M{"_id": bson.ObjectIdHex(qo.PurchaseID)})
	if (purchase) == 1 {
		if "0" == purchases[0].State { //代价单为待报价状态
			//检查是否重复报价
			quotations, err = purchase.QuotationOrder{}.Find("_id", 0, bson.M{}, bson.M{"purchaseID": qo.PurchaseID, "createBy": qo.BuyByID})
		} else {
			err = util.GError{Code: 0, Err: "该报价单非[待报价]状态"}
		}
	} else {
		err = util.GError{Code: 0, Err: "该报价单不存在"}

	}
	return err
}
