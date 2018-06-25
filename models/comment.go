package models

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Comment struct {
	ID             bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id"`
	Type           string        `json:"type" form:"type" query:"type" bson:"type" binding:"required"`
	ArticleID      string        `json:"articleID" form:"articleID" query:"articleID" bson:"articleID" binding:"required"`
	Anonymous      bool          `json:"anonymous" form:"anonymous" query:"anonymous" bson:"anonymous" binding:"exists"`
	Content        string        `json:"content" form:"content" query:"content" bson:"content" binding:"required"`
	IP             string        `json:"ip" form:"ip" query:"ip" bson:"ip" binding:"-"`
	CreateAt       time.Time     `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
	By             string        `json:"by" form:"by" query:"by" bson:"by" binding:"-"`
	AnNickName     string        `json:"anNickName" form:"anNickName" query:"anNickName" bson:"anNickName" binding:"required"`
	ThumbsUps      int           `json:"thumbsUps" form:"thumbsUps" query:"thumbsUps" bson:"thumbsUps" binding:"-"`
	ReplyCommentID string        `json:"replyCommentID" form:"replyCommentID" query:"replyCommentID" bson:"replyCommentID" binding:"-"`
	Reply          Reply         `json:"reply" form:"reply" query:"reply" bson:"-" binding:"-"`
}
type Reply struct {
	UserID     string `json:"userID" form:"userID" query:"userID" bson:"userID" binding:"-"`
	Content    string `json:"content" form:"content" query:"content" bson:"content" binding:"-"`
	AnNickName string `json:"anNickName" form:"anNickName" query:"anNickName" bson:"anNickName" binding:"-"`
	Anonymous  bool   `json:"anonymous" form:"anonymous" query:"anonymous" bson:"anonymous" binding:"-"`
}

const (
	commentCN = "comment"
)

//新增文章
func (ct *Comment) Insert() error {
	return DB.C(commentCN).Insert(ct)
}

//获取排行榜
func (ct Comment) Find(sort string, limit int, selectM bson.M, condition bson.M) ([]Comment, error) {
	var cts = make([]Comment, 0)
	// Select(bson.M{"title": 1, "tags": 1, "occurrenceDate": 1, "location": 1, "taget": 1, "wastage": 1, "nickNamePublish": 1, "content": 1, "state": 1, "auditOpinion": 1})
	err := DB.C(commentCN).Find(condition).Select(selectM).Sort(sort).Limit(limit).All(&cts)
	return cts, err
}

func (ct Comment) Update(condition bson.M) error {
	return DB.C(commentCN).UpdateId(ct.ID, condition)
}

//删除评论
func (ea *Comment) Remove(condition bson.M) error {
	_, err := DB.C(commentCN).RemoveAll(condition)
	return err
}
