package controlers

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/util"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type ThumbsUpControl struct{}

func (ThumbsUpControl) UP(c *gin.Context) {
	var err error
	var l sync.Mutex
	var tb = new(models.Thumb)
	var tbs []models.Thumb
	var expArts []models.ExposureArticle
	if err = c.ShouldBind(tb); nil == err {
		tb.ID = bson.NewObjectId()
		tb.By = middlewares.GetUserIDFromToken(c)
		tb.CreateAt = time.Now()
		tb.IP = strings.Split(c.Request.RemoteAddr, ":")[0]
		l.Lock()
		expArts, err = models.ExposureArticle{}.Find("_id", 1, bson.M{"thumbsUps": 1, "_id": 1}, bson.M{"_id": bson.ObjectIdHex(tb.ArticleID)})
		if nil == err && 0 < len(expArts) {
			//先查询是否已经点过赞
			fmt.Println("articleID", tb.ArticleID, "by", tb.By)
			var cond = bson.M{"articleID": tb.ArticleID, "by": tb.By}
			if "2" == tb.Type && "" != tb.CommentID {
				cond["commentID"] = tb.CommentID
			}
			tbs, err = models.Thumb{}.Find(cond)
			fmt.Println("tbs", len(tbs), tbs)
			if 1 > len(tbs) {

				//没有点赞记录
				err = tb.Insert()
				if "1" == tb.Type && nil == err { //更新文章点赞数
					models.ExposureArticle{ID: expArts[0].ID}.UpdateCond(bson.M{"$set": bson.M{"thumbsUps": expArts[0].ThumbsUps + 1}})
				} else { //点赞评论

					cts, err := models.Comment{}.Find("_id", 1, bson.M{"_id": 1, "thumbsUps": 1}, bson.M{"_id": bson.ObjectIdHex(tb.CommentID)})
					fmt.Println("CommentID", tb.CommentID, err, cts)

					if nil == err && 0 < len(cts) {
						count := cts[0].ThumbsUps + 1
						models.Comment{ID: cts[0].ID}.Update(bson.M{"$set": bson.M{"thumbsUps": count}})
					}
				}
			} else {

				err = &util.GError{Code: 0, Err: "您已点过赞啦"}
			}
		}
		l.Unlock()
	}

	util.JSON(c, util.ResponseMesage{Message: "点赞", Data: tb, Error: err})
}

func (ThumbsUpControl) List(c *gin.Context) {
	var err error
	var tb = new(models.Thumb)
	var ths []models.Thumb
	if err = c.ShouldBind(tb); nil == err {
		// tb.By = middlewares.GetUserIDFromToken(c)

		ths, err = models.Thumb{}.Find(bson.M{"articleID": tb.ArticleID})

	}
	util.JSON(c, util.ResponseMesage{Message: "获取点赞", Data: ths, Error: err})
}
