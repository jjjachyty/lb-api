package controlers

import (
	"fmt"
	"lb-api/models"
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
	cond = bson.M{"$and": []bson.M{bson.M{"state": "1"}, appCond}}

	fmt.Println("List--cond", cond)
	result, err := models.Purchase{}.Find(sort, 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}
func (PurchaseControl) Get(c *gin.Context) {
	id := c.Query("id")
	var cond bson.M
	var result models.Purchase
	var results []models.Purchase
	var err error
	if "" != id {
		cond = bson.M{"_id": bson.ObjectIdHex(id), "state": "1"}

		results, err = models.Purchase{}.Find("_id", 10, bson.M{}, cond)
		if len(results) > 0 {
			result = results[0]
		}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}
