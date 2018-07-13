package order

import (
	"lb-api/models/purchase"
	"time"

	"labix.org/v2/mgo/bson"
)

// Order struct 订单实体
type Order struct {
	ID bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id,omitempty" binding:"-"`
	// BuyAmount        float64       `json:"buyAmount" form:"buyAmount" query:"buyAmount" bson:"buyAmount" binding:"required"`
	BuyBy  string `json:"buyBy" form:"buyBy" query:"buyBy" bson:"buyBy" binding:"required"`
	SellBy string `json:"sellBy" form:"sellBy" query:"sellBy" bson:"sellBy" binding:"required"`
	// SellAmount       float64       `json:"sellAmount" form:"sellAmount" query:"sellAmount" bson:"sellAmount" binding:"required"`
	Products         []purchase.Product `json:"products" form:"products" query:"products" bson:"products" binding:"required"`
	Type             string             `json:"type" form:"type" query:"type" bson:"type" binding:"required"` // 订单类型、代购、转卖、旅拍
	StrikePrice      float64            `json:"strikePrice" form:"strikePrice" query:"strikePrice" bson:"strikePrice" binding:"required"`
	Ticket           string             `json:"ticket" form:"ticket" query:"ticket" bson:"ticket" binding:"-"`
	ExpressName      string             `json:"expressName" form:"expressName" query:"expressName" bson:"expressName" binding:"-"`
	ExpressNumber    string             `json:"expressNumber" form:"expressNumber" query:"expressNumber" bson:"expressNumber" binding:"-"`
	ExpressArrivedAt string             `json:"expressArrivedAt" form:"expressArrivedAt" query:"expressArrivedAt" bson:"expressArrivedAt" binding:"-"`
	ExpressState     string             `json:"expressState" form:"expressState" query:"expressState" bson:"expressState" binding:"-"`
	CreateAt         time.Time          `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
	State            string             `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
}

const orderCN = "order"
