package models

import (
	"labix.org/v2/mgo"
)

const (
	DbName         = "luobo"
	UserCollection = "user"
)

var DB *mgo.Database

// var Session *mgo.Session.DB()

func init() {
	session, err := mgo.Dial("106.12.10.77:27017")
	DB = session.DB(DbName)
	//session, err := mgo.Dial("localhost")
	if nil != err {
		panic("127.0.0.1:27017 数据库链接失败")
	}
}
