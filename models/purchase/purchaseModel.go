package purchase

import (
	"lb-api/models"
	"time"

	"labix.org/v2/mgo/bson"
)

type Product struct {
	ID       bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id,omitempty" binding:"-"`              //商品编号
	Name     string        `json:"name" form:"name" query:"name" bson:"name" binding:"-"`                 //商品名称
	Price    float64       `json:"price" form:"price" query:"price" bson:"price" binding:"-"`             //商品参考价格
	Describe string        `json:"describe" form:"describe" query:"describe" bson:"describe" binding:"-"` //商品描述
	Images   string        `json:"images" form:"images" query:"images" bson:"images" binding:"-"`         //商品图片
	ShopName string        `json:"shopName" form:"shopName" query:"shopName" bson:"shopName" binding:"-"` //购买平台
	Quantity int64         `json:"quantity" form:"quantity" query:"quantity" bson:"quantity" binding:"-"`
	// ShopLocation string        `json:"name" form:"name" query:"name" bson:"name" binding:"-"`                 //购买国家
}

//代购发起人信息
type Purchase struct {
	ID          bson.ObjectId  `json:"id" form:"id" query:"id" bson:"_id,omitempty" binding:"-"`
	Content     string         `json:"content" form:"content" query:"content" bson:"content" binding:"required"` //内容描述
	Amount      float64        `json:"amount" form:"amount" query:"amount" bson:"amount" binding:"exists"`       //内容描述
	Products    []Product      `json:"products" form:"products[]" query:"products" bson:"products" binding:"required"`
	Address     models.Address `json:"address" form:"address" query:"address" bson:"address" binding:"required"`
	Destination string         `json:"destination" form:"destination" query:"destination" bson:"destination" binding:"required"` //目的地
	CreateBy    string         `json:"createBy" form:"createBy" query:"createBy" bson:"createBy" binding:"-"`                    //创建人
	Creator     string         `json:"creator" form:"creator" query:"creator" bson:"creator" binding:"-"`
	CreateAt    time.Time      `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"` //创建时间
	UpdateAt    time.Time      `json:"updateAt" form:"updateAt" query:"updateAt" bson:"updateAt" binding:"-"` //更新时间
	State       string         `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
	Views       int64          `json:"views" form:"views" query:"views" bson:"views" binding:"-"`
	QuotationID string         `json:"quotationID" form:"quotationID" query:"quotationID" bson:"quotationID" binding:"-"`
	// QuotationOrders []QuotationOrder `json:"quotationOrders" form:"-" query:"-" bson:"-" binding:"-"`
	Inviters []string `json:"inviters" form:"-" query:"-" bson:"inviters" binding:"-"`
}

const (
	purchaseCN = "purchase"
)

func (Purchase) Find(sort []string, limit int, selectM bson.M, condition bson.M) ([]Purchase, error) {
	var purchase = make([]Purchase, 0)
	query := models.DB.C(purchaseCN).Find(condition)
	if len(sort) > 0 {
		query = query.Sort(sort...)
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

//Update 更新代购单
func (Purchase) Update(selector bson.M, update bson.M) error {
	return models.DB.C(purchaseCN).Update(selector, update)
}

//Update 更新代购单
func (p *Purchase) Insert() error {
	return models.DB.C(purchaseCN).Insert(p)
}
