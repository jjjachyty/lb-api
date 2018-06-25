package models

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type Product struct {
	ID       bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`                        //商品编号
	Name     string        `json:"name" form:"name" query:"name" bson:"name" binding:"-"`                 //商品名称
	Price    float64       `json:"price" form:"price" query:"price" bson:"price" binding:"-"`             //商品参考价格
	Describe string        `json:"describe" form:"describe" query:"describe" bson:"describe" binding:"-"` //商品描述
	Images   string        `json:"images" form:"images" query:"images" bson:"images" binding:"-"`         //商品图片
	ShopName string        `json:"shopName" form:"shopName" query:"shopName" bson:"shopName" binding:"-"` //购买平台
	Quantity int           `json:"quantity" form:"quantity" query:"quantity" bson:"quantity" binding:"-"`
	// ShopLocation string        `json:"name" form:"name" query:"name" bson:"name" binding:"-"`                 //购买国家
}

//报价单
type QuotationOrder struct {
	ID         bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	PurchaseID bson.ObjectId //代购单号
	Amount     float64       //总金额
	Products   []Product     `json:"products" form:"products[]" query:"products" bson:"products" binding:"required"`
	Charge     float64       //服务费
	BuyByID    string        //报价人ID
	BuyByName  string        //报价人昵称
	CreateAt   time.Time     //报价时间
	State      string        //报价单状态
	TimeOut    time.Time     //失效时间
}

//代购发起人信息
type Purchase struct {
	ID          bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	Content     string        `json:"content" form:"content" query:"content" bson:"content" binding:"required"` //内容描述
	Amount      float64       `json:"amount" form:"amount" query:"amount" bson:"amount" binding:"required"`     //内容描述
	Products    []Product     `json:"products" form:"products[]" query:"products" bson:"products" binding:"required"`
	Address     string        `json:"address" form:"address" query:"address" bson:"address" binding:"required"`
	Location    string        `json:"location" form:"location" query:"location" bson:"location" binding:"required"` //目的地
	CreateBy    string        `json:"createBy" form:"createBy" query:"createBy" bson:"createBy" binding:"required"` //创建人
	Creator     string        `json:"creator" form:"creator" query:"creator" bson:"creator" binding:"required"`
	CreatAt     time.Time     `json:"creatAt" form:"creatAt" query:"creatAt" bson:"creatAt" binding:"required"`     //创建时间
	UpdateAt    time.Time     `json:"updateAt" form:"updateAt" query:"updateAt" bson:"updateAt" binding:"required"` //更新时间
	State       string        `json:"state" form:"state" query:"state" bson:"state" binding:"required"`
	QuotationID bson.ObjectId `json:"quotationID" form:"quotationID" query:"quotationID" bson:"quotationID" binding:"-"`
}

const (
	purchaseCN = "purchase"
)

func (Purchase) Find(sort string, limit int, selectM bson.M, condition bson.M) ([]Purchase, error) {
	var purchase = make([]Purchase, 0)
	query := DB.C(purchaseCN).Find(condition)
	if "" != sort {
		query = query.Sort(sort)
	}
	if 0 != limit {
		query = query.Limit(limit)
	}
	if len(selectM) > 0 {
		query = query.Select(selectM)
	}
	err := query.All(&purchase)
	return purchase, err
}
