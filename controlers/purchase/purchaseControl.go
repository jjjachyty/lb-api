package purchase

import (
	"lb-api/models/purchase"
	"lb-api/util"

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
