package purchase

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/purchase"
	"lb-api/util"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

type PurchaseControl struct{}

func (PurchaseControl) List(c *gin.Context) {
	sort := c.DefaultQuery("sort", "-updateAt")
	keyWords := c.Query("keyWords")
	var cond bson.M
	var appCond bson.M
	if "" != keyWords {
		appCond = bson.M{"$or": []bson.M{bson.M{"content": bson.M{"$regex": keyWords}}, bson.M{"products.name": bson.M{"$regex": keyWords}}, bson.M{"products.describe": bson.M{"$regex": keyWords}}, bson.M{"location": bson.M{"$regex": keyWords}}}}

	}
	cond = bson.M{"$and": []bson.M{bson.M{"state": "0"}, appCond}}

	result, err := purchase.Purchase{}.Find(sort, 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}
func (PurchaseControl) UserList(c *gin.Context) {
	var cond bson.M
	userID := middlewares.GetUserIDFromToken(c)
	cond = bson.M{"createBy": userID}
	result, err := purchase.Purchase{}.Find("-createAt", 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取我的代购单列表", Data: result, Error: err})

}

func (PurchaseControl) Get(c *gin.Context) {
	id := c.Query("id")
	var cond bson.M
	var result purchase.Purchase
	var results []purchase.Purchase

	var qos []purchase.QuotationOrder
	var err error
	if "" != id {
		cond = bson.M{"_id": bson.ObjectIdHex(id)}

		results, err = purchase.Purchase{}.Find("_id", 10, bson.M{}, cond)
		if len(results) > 0 {
			result = results[0]
			//查询报价单
			qos, err = purchase.QuotationOrder{}.Find("-createAt", 0, bson.M{}, bson.M{"purchaseID": id})
			if nil == err {
				result.QuotationOrders = qos
			}

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
		purchases, err = purchase.Purchase{}.Find("_id", 1, bson.M{}, bson.M{"_id": purchaseObj.ID})
		if len(purchases) == 1 {
			dbPurchase := purchases[0]
			if dbPurchase.CreateBy == middlewares.GetUserIDFromToken(c) {
				if "0" == dbPurchase.State {
					var amount = 0.0
					for _, p := range purchaseObj.Products {
						quantity, _ := strconv.ParseFloat(strconv.Itoa(p.Quantity), 64)
						amount += p.Price * quantity
					}
					err = purchase.Purchase{}.Update(bson.M{"_id": purchaseObj.ID}, bson.M{"$set": bson.M{"amount": amount, "targetLocation": purchaseObj.TargetLocation, "address": purchaseObj.Address, "content": purchaseObj.Content, "products": purchaseObj.Products, "updateAt": purchaseObj.UpdateAt}})
				} else {
					err = &util.GError{Code: 0, Err: "只能更新状态为[待报价]的订购单"}
				}
			} else {
				err = &util.GError{Code: 0, Err: "不能操作他人代购单"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "该报价单不存在"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "更新物流代购", Data: purchaseObj, Error: err})

}

// Delete func 删除
func (PurchaseControl) Remove(c *gin.Context) {
	var id = c.Query("id")
	var err error
	var purchases []purchase.Purchase
	if bson.IsObjectIdHex(id) {
		purchases, err = purchase.Purchase{}.Find("_id", 1, bson.M{}, bson.M{"_id": bson.ObjectIdHex(id)})
		if len(purchases) == 1 {
			dbPurchase := purchases[0]
			if "0" == dbPurchase.State { //待报价状态可删除

				if dbPurchase.CreateBy == middlewares.GetUserIDFromToken(c) {
					err = models.Remove("purchase", bson.M{"_id": bson.ObjectIdHex(id)})
					util.Glog.Debugf("删除代购单-操作人%s-原数据%v-状态%v", dbPurchase.CreateBy, dbPurchase, err)
					if nil == err { //删除成功后更新报价单
						go purchase.QuotationOrder{}.Update(bson.M{"purchaseID": id}, bson.M{"state": "-1", "refuseReason": "该报价单已被删除"})
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
