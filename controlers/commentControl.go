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

type CommentControl struct{}

func (CommentControl) Add(c *gin.Context) {
	var l sync.Mutex
	var err error
	var ct = new(models.Comment)
	var cts []models.Comment
	var expArts []models.ExposureArticle
	if err = c.ShouldBind(ct); nil == err {
		if len(ct.Content) > 10 && len(ct.Content) < 140 {
			l.Lock()
			//查询文章
			expArts, err = models.ExposureArticle{}.Find("_id", 1, bson.M{"comments": 1, "_id": 1, "createUser": 1, "title": 1}, bson.M{"_id": bson.ObjectIdHex(ct.ArticleID)})
			if nil == err && len(expArts) > 0 {
				ct.ID = bson.NewObjectId()
				ct.CreateAt = time.Now()
				ct.By = middlewares.GetUserIDFromToken(c)
				ct.IP = strings.Split(c.Request.RemoteAddr, ":")[0] //

				if "1" == ct.Type { //文章

					expArts[0].Comments++
					err := models.ExposureArticle{ID: expArts[0].ID}.UpdateCond(bson.M{"$set": bson.M{"comments": expArts[0].Comments}})
					util.Glog.Debugf("新增文章[%s]评论-[%v]-原评论数[%d]", expArts[0].ID, err, expArts[0].Comments)
					go func() {
						models.Message{ID: bson.NewObjectId(), Type: "1", To: expArts[0].CreateBy, From: ct.By, Content: "您的文章<<a href='/article/" + expArts[0].ID.Hex() + "'>" + expArts[0].Title + "</a>>收到1个评论:<br/>< <small class='grey--text'>" + ct.Content + "</small>>", CreateAt: time.Now(), State: "1"}.Insert()
					}()

				} else { //评论

					cts, err = models.Comment{}.Find("_id", 1, bson.M{}, bson.M{"_id": bson.ObjectIdHex(ct.ReplyCommentID)})
					if len(cts) > 0 {
						go func() { //更新评论回复消息
							models.Message{ID: bson.NewObjectId(), Type: "2", To: cts[0].By, From: ct.By, Content: "您在文章<<a href='/article/" + expArts[0].ID.Hex() + "'>" + expArts[0].Title + "</a>>的评论收到1个回复:<br/>< <small class='grey--text'>" + ct.Content + "</small>>", CreateAt: time.Now(), State: "1"}.Insert()
						}()
					} else {
						err = &util.GError{Code: 0, Err: "评论不存在"}
					}
				}
				err = ct.Insert()
				l.Unlock()

			} else {
				err = &util.GError{Code: 0, Err: "文章不存在"}
			}
		} else {
			err = &util.GError{Code: 0, Err: "评论内容需在10～140字之间"}
		}
	}
	util.JSON(c, util.ResponseMesage{Message: "新增最新评论", Data: ct, Error: err})

}
func (ct CommentControl) HotList(c *gin.Context) {
	var err error
	var cts []models.Comment
	var articleID = c.Query("articleID")
	//根据时间查处前10个
	cts, err = models.Comment{}.Find("-thumbsUps", 5, bson.M{"ip": 0}, bson.M{"$and": []bson.M{bson.M{"articleID": articleID}, bson.M{"thumbsUps": bson.M{"$gt": 0}}}})

	util.JSON(c, util.ResponseMesage{Message: "获取最热评论", Data: ct.hander(cts), Error: err})

}

func (CommentControl) NewList(c *gin.Context) {
	var err error
	var cts []models.Comment
	var articleID = c.Query("articleID")
	var lastID = c.Query("lastID")
	fmt.Println("articleID", articleID, "lastID", lastID)
	if "" != articleID { //
		if "" != lastID {
			cts, err = models.Comment{}.Find("-createAt", 10, bson.M{"ip": 0}, bson.M{"$and": []bson.M{bson.M{"articleID": articleID}, bson.M{"_id": bson.M{"$lt": bson.ObjectIdHex(lastID)}}}})
		} else {

			cts, err = models.Comment{}.Find("-createAt", 10, bson.M{"ip": 0}, bson.M{"articleID": articleID})
		}

	} else {
		err = &util.GError{Code: 0, Err: "文章号不能为空"}
	}
	for k, v := range cts {
		fmt.Println(k, v.ID.Hex())
	}
	util.JSON(c, util.ResponseMesage{Message: "获取最新评论", Data: cts, Error: err})

}

func (CommentControl) hander(comments []models.Comment) []models.Comment {
	var repliyIDs []string
	if 0 < len(comments) {

		for i, k := range comments {
			if k.Anonymous {
				comments[i].By = ""
			}
			repliyIDs = append(repliyIDs, k.ReplyCommentID)
		}
		//查询所有的回复
		cts, err := models.Comment{}.Find("_id", 10, bson.M{"ip": 0}, bson.M{"$and": []bson.M{bson.M{"articleID": comments[0].ArticleID}, bson.M{"replyCommentID": bson.M{"$in": repliyIDs}}}})
		fmt.Println("cts", cts, comments[0].ArticleID, repliyIDs)
		//处理回复
		if nil == err {
			for _, n := range cts {
				for o, p := range comments {
					fmt.Println("p.ID.Hex() ", p.ID.Hex(), n.ReplyCommentID)
					if p.ReplyCommentID == n.ID.Hex() {
						comments[o].Reply = models.Reply{n.By, n.Content, n.AnNickName, n.Anonymous}
					}
				}
			}
		}
	}
	return comments
}
