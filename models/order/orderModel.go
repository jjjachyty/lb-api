package order

import (
	"lb-api/models/purchase"
	"time"

	"labix.org/v2/mgo/bson"
)

type Express struct {
	Name          string    `json:"name" form:"name" query:"name" bson:"name" binding:"-"`                                     //快递公司
	Number        string    `json:"number" form:"number" query:"number" bson:"number" binding:"-"`                             //快递单号
	ArrivedAt     string    `json:"arrivedAt" form:"arrivedAt" query:"arrivedAt" bson:"arrivedAt" binding:"-"`                 //当前到达点
	Courier       string    `json:"courier" form:"courier" query:"courier" bson:"courier" binding:"-"`                         //派送人
	ContactNumber string    `json:"contactNumber" form:"contactNumber" query:"contactNumber" bson:"contactNumber" binding:"-"` //联系电话
	CreateAt      time.Time `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`                     //创建时间
	UpdateAt      time.Time `json:"updateAt" form:"updateAt" query:"updateAt" bson:"updateAt" binding:"-"`                     //更新时间
	State         string    `json:"state" form:"state" query:"state" bson:"state" binding:"-"`                                 //快递状态
}

// Order struct 订单实体
type Order struct {
	ID bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id,omitempty" binding:"-"`
	// BuyAmount        float64       `json:"buyAmount" form:"buyAmount" query:"buyAmount" bson:"buyAmount" binding:"required"`
	BuyBy string `json:"buyBy" form:"buyBy" query:"buyBy" bson:"buyBy" binding:"required"`

	SellBy string `json:"sellBy" form:"sellBy" query:"sellBy" bson:"sellBy" binding:"required"`
	Seller string `json:"seller" form:"seller" query:"seller" bson:"seller" binding:"-"`
	// SellAmount       float64       `json:"sellAmount" form:"sellAmount" query:"sellAmount" bson:"sellAmount" binding:"required"`
	OriginalLink  string             `json:"originalLink" form:"originalLink" query:"originalLink" bson:"originalLink" binding:"required"`
	Products      []purchase.Product `json:"products" form:"products" query:"products" bson:"products" binding:"required"`
	Type          string             `json:"type" form:"type" query:"type" bson:"type" binding:"required"` // 订单类型、代购、转卖、旅拍
	StrikePrice   float64            `json:"strikePrice" form:"strikePrice" query:"strikePrice" bson:"strikePrice" binding:"required"`
	Charge        float64            `json:"charge" form:"charge" query:"charge" bson:"charge" binding:"required"` //服务费
	Ticket        string             `json:"ticket" form:"ticket" query:"ticket" bson:"ticket" binding:"-"`
	TicketExplain string             `json:"ticketExplain" form:"ticketExplain" query:"ticketExplain" bson:"ticketExplain" binding:"-"`
	Express       Express            `json:"express" form:"express" query:"express" bson:"express" binding:"-"`
	Reviews       string             `json:"reviews" form:"reviews" query:"reviews" bson:"reviews" binding:"-"`
	ReturnReason  string             `json:"returnReason" form:"returnReason" query:"returnReason" bson:"returnReason" binding:"-"` //退货原因
	CancelReason  string             `json:"cancelReason" form:"cancelReason" query:"cancelReason" bson:"cancelReason" binding:"-"` //取消订单原因
	CreateAt      time.Time          `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
	State         string             `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
}

const orderCN = "order"
