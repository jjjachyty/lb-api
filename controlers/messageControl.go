package controlers

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/util"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

type MessageControl struct{}

func (MessageControl) GetUserMessage(c *gin.Context) {
	var mgs []models.Message
	var err error
	var cond = bson.M{}
	userid := middlewares.GetUserIDFromToken(c)

	var condAnd = []bson.M{bson.M{"to": userid}}
	var lastID = c.Query("lastID")
	if "" != userid {
		if "" != lastID { //分页查询
			condAnd = append(condAnd, bson.M{"_id": bson.M{"$lt": bson.ObjectIdHex(lastID)}})
		}
		cond["$and"] = condAnd
		fmt.Println("cond", cond)
		mgs, err = models.Message{To: userid}.FindMessage(cond)

	}
	util.JSON(c, util.ResponseMesage{Message: "获取我的消息", Data: mgs, Error: err})

}

//GetNewMessageCount 获取用户新消息
func (MessageControl) GetNewMessageCount(c *gin.Context) {
	userid := middlewares.GetUserIDFromToken(c)

	count, err := models.Message{To: userid}.GetNewMessageCount()

	util.JSON(c, util.ResponseMesage{Message: "获取我的新消息个数", Data: count, Error: err})

}

//GetNewMessageCount 获取用户新消息
func (MessageControl) Remove(c *gin.Context) {
	var cond bson.M
	var err error

	userid := middlewares.GetUserIDFromToken(c)

	removeType := c.Query("removeType")
	if "" != removeType {

		if "1" == removeType { //删除所有
			cond = bson.M{"to": userid}
		} else {
			cond = bson.M{"to": userid, "state": "0"}
		}
		err = models.Message{}.Remove(cond)
	} else {
		err = &util.GError{Code: 0, Err: "删除类型不能为空"}
	}
	util.JSON(c, util.ResponseMesage{Message: "清除消息", Data: nil, Error: err})

}

//GetNewMessageCount 获取用户新消息
func (MessageControl) Update(c *gin.Context) {
	var cond bson.M
	var err error
	var updateIDs []bson.ObjectId
	userid := middlewares.GetUserIDFromToken(c)

	updateType := c.PostForm("updateType")
	messageids := c.PostFormArray("messageids[]")

	if len(messageids) > 0 {
		for _, msg := range messageids {
			updateIDs = append(updateIDs, bson.ObjectIdHex(msg))
		}
	}
	if "" != updateType {

		if "1" == updateType { //所有标记为已读
			cond = bson.M{"to": userid}
		} else { //单个标记
			cond = bson.M{"_id": bson.M{"$in": updateIDs}}
		}
		fmt.Println("更新未读消息", messageids, "更新人", userid, cond)

		err = models.Message{To: userid}.Update(cond, bson.M{"$set": bson.M{"state": "0"}})
	} else {
		err = &util.GError{Code: 0, Err: "更新类型不能为空"}
	}
	util.JSON(c, util.ResponseMesage{Message: "更新消息", Data: nil, Error: err})

}
