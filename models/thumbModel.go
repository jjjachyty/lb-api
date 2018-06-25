package models

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Thumb struct {
	ID        bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	Type      string        `json:"type" form:"type" query:"type" bson:"type" binding:"-"` //文章 评论
	ArticleID string        `json:"articleID" form:"articleID" query:"articleID" bson:"articleID" binding:"required"`
	CommentID string        `json:"commentID" form:"commentID" query:"commentID" bson:"commentID" binding:"-"`
	By        string        `json:"by" form:"by" query:"by" bson:"by" binding:"-"`
	IP        string        `json:"ip" form:"ip" query:"ip" bson:"ip" binding:"-"`
	CreateAt  time.Time     `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
}

const (
	thumbCN = "thumb"
)

//新增文章
func (tb *Thumb) Insert() error {
	return DB.C(thumbCN).Insert(tb)
}

func (th Thumb) Find(condition bson.M) ([]Thumb, error) {
	var ths = make([]Thumb, 0)
	err := DB.C(thumbCN).Find(condition).All(&ths)
	return ths, err
}
func (th Thumb) Update(condition bson.M) error {
	return DB.C(thumbCN).UpdateId(th.ID, condition)

}
