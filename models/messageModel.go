package models

import (
	"lb-api/util"
	"time"

	"labix.org/v2/mgo/bson"
)

type Message struct {
	ID       bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	From     string        `json:"from" form:"from" query:"from" bson:"from" binding:"required"`
	Type     string        `json:"type" form:"type" query:"type" bson:"type" binding:"required"`
	To       string        `json:"to" form:"to" query:"to" bson:"to" binding:"required"`
	Content  string        `json:"content" form:"content" query:"content" bson:"content" binding:"required"`
	CreateAt time.Time     `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
	State    string        `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
}

const (
	messageCN = "message"
)

func (m Message) FindMessage(condition bson.M) ([]Message, error) {
	var result = make([]Message, 0)
	err := DB.C(messageCN).Find(condition).Limit(5).Sort("-createAt", "-state").All(&result)
	return result, err
}
func (m Message) GetNewMessageCount() (int, error) {
	return DB.C(messageCN).Find(bson.M{"to": m.To, "state": "1"}).Count()
}

func (m Message) Insert() error {
	util.Glog.Debugf("新增通知消息-发送人:%s 接收者:%s 内容:%s", m.From, m.To, m.Content)
	return DB.C(messageCN).Insert(m)
}

func (m Message) Remove(condition bson.M) error {
	_, err := DB.C(messageCN).RemoveAll(condition)
	util.Glog.Debugf("删除通知消息-操作人[%s]", m.To)

	return err
}

func (m Message) Update(selector bson.M, update bson.M) error {
	_, err := DB.C(messageCN).UpdateAll(selector, update)
	util.Glog.Debugf("更新通知消息-更新字段%v  更新条件 %v 操作人[%s]", update, selector, m.To)

	return err
}
