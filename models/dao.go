package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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

func Find(CN string, results interface{}, sort string, limit int, selectM bson.M, condition bson.M) error {
	// var results = make([]interface{}, 0)
	query := DB.C(CN).Find(condition)
	if "" != sort {
		query = query.Sort(sort)
	}
	if 0 != limit {
		query = query.Limit(limit)
	}
	if len(selectM) > 0 {
		query = query.Select(selectM)
	}
	err := query.All(results)
	return err
}

// Remove func 删除
func Remove(cn string, selector bson.M) error {
	return DB.C(cn).Remove(selector)
}
