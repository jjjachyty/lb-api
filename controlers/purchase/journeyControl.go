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

type JourneyControl struct{}

func (JourneyControl) List(c *gin.Context) {
	sort := c.DefaultQuery("sort", "-updateAt")
	keyWords := c.Query("keyWords")
	var cond bson.M
	var appCond bson.M
	if "" != keyWords {
		appCond = bson.M{"$or": []bson.M{bson.M{"content": bson.M{"$regex": keyWords}}, bson.M{"products.name": bson.M{"$regex": keyWords}}, bson.M{"products.describe": bson.M{"$regex": keyWords}}, bson.M{"location": bson.M{"$regex": keyWords}}}}

	}
	cond = bson.M{"$and": []bson.M{bson.M{"state": "0"}, appCond}}

	result, err := purchase.Journey{}.Find(sort, 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})
}

func (JourneyControl) DestinationList(c *gin.Context) {
	sort := c.DefaultQuery("sort", "endDate")
	destination := c.Query("destination")
	user := c.Query("user")
	var cond bson.M
	if "" != destination {
		cond = bson.M{"destination": destination, "state": "1"}
		if "" != user {
			cond["createBy"] = bson.M{"$ne": user}
		}

	}
	// cond = bson.M{"$and": []bson.M{bson.M{"state": "1"}, appCond}}

	result, err := purchase.Journey{}.Find(sort, 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取可代购推荐", Data: result, Error: err})
}

func (JourneyControl) UserList(c *gin.Context) {
	var cond bson.M
	userID := middlewares.GetUserIDFromToken(c)
	cond = bson.M{"createBy": userID}
	result, err := purchase.Journey{}.Find("-createAt", 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取我的旅程", Data: result, Error: err})

}

func (JourneyControl) Get(c *gin.Context) {
	id := c.Query("id")
	var cond bson.M
	var result purchase.Journey
	var results []purchase.Journey

	var err error
	if "" != id {
		cond = bson.M{"_id": bson.ObjectIdHex(id)}

		results, err = purchase.Journey{}.Find("_id", 10, bson.M{}, cond)
		if len(results) > 0 {
			result = results[0]
		}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: result, Error: err})

}

func (JourneyControl) Add(c *gin.Context) {
	var err error
	var journey = new(purchase.Journey)
	var journeys []purchase.Journey
	if err = c.ShouldBind(journey); nil == err {
		journey.ID = bson.NewObjectId()
		journey.CreateAt = time.Now()
		journey.UpdateAt = journey.CreateAt
		journey.State = "1" //新增状态
		journey.CreateBy = middlewares.GetUserIDFromToken(c)
		if journey.EndDate.Before(journey.StartDate) { //检查开始结束时间
			err = &util.GError{Code: 0, Err: "结束时间必须大于开始时间"}
		} else {
			//检查是否有重复的行程
			journeys, err = purchase.Journey{}.Find("_id", 10, bson.M{}, bson.M{"endDate": bson.M{"$gt": journey.StartDate}, "state": "1"})
			if len(journeys) == 0 {
				err = journey.Insert()
			} else {
				err = &util.GError{Code: 0, Err: "与去往" + journeys[0].Destination + "的行程时间有重叠"}
			}

		}
	}
	fmt.Println("err", err)
	util.JSON(c, util.ResponseMesage{Message: "获取物流代购列表", Data: journey, Error: err})

}

func (JourneyControl) Update(c *gin.Context) {
	var err error
	var journeys []purchase.Journey
	var journeyObj = new(purchase.Journey)
	var id = c.PostForm("id")
	if err = c.ShouldBind(journeyObj); nil == err && "" != id {

		journeys, err = purchase.Journey{}.Find("_id", 1, bson.M{}, bson.M{"_id": bson.ObjectIdHex(id)})
		if len(journeys) == 1 {
			dbJourney := journeys[0]
			if dbJourney.CreateBy == middlewares.GetUserIDFromToken(c) {
				err = purchase.Journey{}.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"chargeType": journeyObj.ChargeType, "chargeValue": journeyObj.ChargeValue, "startDate": journeyObj.StartDate, "endDate": journeyObj.EndDate, "destination": journeyObj.Destination, "remarks": journeyObj.Remarks, "products": journeyObj.Products, "updateAt": time.Now()}})
				util.Glog.Debugf("更新行程单-原数据%v-新数据-%v", dbJourney, journeyObj)
			} else {
				err = &util.GError{Code: 0, Err: "不能操作他人行程单"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "该行程不存在"}
		}

	} else {
		err = &util.GError{Code: 0, Err: "表单数据不完整"}
	}
	fmt.Println("errrrrrr", err, c.PostForm("id"))
	util.JSON(c, util.ResponseMesage{Message: "更新我的行程", Data: journeyObj, Error: err})
}

// Delete func 删除
func (JourneyControl) Remove(c *gin.Context) {
	var id = c.Query("id")
	var err error
	var journeys []purchase.Journey
	if bson.IsObjectIdHex(id) {
		journeys, err = purchase.Journey{}.Find("_id", 1, bson.M{}, bson.M{"_id": bson.ObjectIdHex(id)})
		if len(journeys) == 1 {
			dbJourney := journeys[0]
			if dbJourney.CreateBy == middlewares.GetUserIDFromToken(c) {
				err = purchase.Journey{}.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
				util.Glog.Debugf("删除行程单-操作人%s-原数据%v-状态%v", dbJourney.CreateBy, dbJourney, err)
			} else {
				err = &util.GError{Code: 0, Err: "不能操作他人行程单"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "该行程不存在"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "删除我的行程", Data: nil, Error: err})
}
