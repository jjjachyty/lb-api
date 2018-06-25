package models

import (
	"time"

	"labix.org/v2/mgo/bson"
)

const (
	tipOffCN = "tipoff"
)

type TipOff struct {
	ID           bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	ArticleID    string        `json:"articleID" form:"articleID" query:"articleID" bson:"articleID" binding:"required"`
	CommentID    string        `json:"commentID" form:"commentID" query:"commentID" bson:"commentID"`
	Category     string        `json:"category" form:"category" query:"category" bson:"category" binding:"required"`
	Content      string        `json:"content" form:"content" query:"content" bson:"content" binding:"required"`
	CreateAt     time.Time     `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
	CreatBy      string        `json:"creatBy" form:"creatBy" query:"creatBy" bson:"creatBy" binding:"-"`
	Auditor      string        `json:"auditor" form:"auditor" query:"auditor" bson:"auditor" binding:"-"`
	AuditOpinion string        `json:"auditOpinion" form:"auditOpinion" query:"auditOpinion" bson:"auditOpinion" binding:"-"`
	AuditAt      time.Time     `json:"auditAt" form:"auditAt" query:"auditAt" bson:"auditAt" binding:"-"`
}

func (tip *TipOff) Insert() error {
	return DB.C(tipOffCN).Insert(tip)
}
