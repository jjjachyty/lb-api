package order

import (
	"lb-api/models"
	"time"

	"labix.org/v2/mgo/bson"
)

//Transaction struct交易表
type Transaction struct {
	ID              bson.ObjectId `json:"id" form:"id" query:"id" bson:"id" binding:"-"`
	OrderID         string        `json:"orderID" form:"orderID" query:"orderID" bson:"orderID" binding:"required"`
	Buyer           string        `json:"buyer" form:"buyer" query:"buyer" bson:"buyer" binding:"required"`
	BuyerPreAmount  float64       `json:"buyerBeforeAmount" form:"buyerBeforeAmount" query:"buyerBeforeAmount" bson:"buyerBeforeAmount" binding:"required"`
	Seller          string        `json:"seller" form:"seller" query:"seller" bson:"seller" binding:"required"`
	SellerPreAmount float64       `json:"sellerBeforeAmount" form:"sellerBeforeAmount" query:"sellerBeforeAmount" bson:"sellerBeforeAmount" binding:"required"`
	CreateAt        time.Time     `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"required"`
	Amount          float64       `json:"amount" form:"amount" query:"amount" bson:"amount" binding:"required"`
	State           string        `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
	StateExplain    string        `json:"stateExplain" form:"stateExplain" query:"stateExplain" bson:"stateExplain" binding:"-"`
}

//Transaction struct交易表
type ApplyCash struct {
	ID           bson.ObjectId   `json:"id" form:"id" query:"id" bson:"id" binding:"-"`
	PreAmount    float64         `json:"preAmount" form:"preAmount" query:"preAmount" bson:"preAmount" binding:"-"` //体现前金额
	Amount       float64         `json:"amount" form:"amount" query:"amount" bson:"amount" binding:"required"`      //提现金额
	UserName     string          `json:"userName" form:"userName" query:"userName" bson:"userName" binding:"-"`     //提现金额
	Phone        string          `json:"phone" form:"phone" query:"phone" bson:"phone" binding:"-"`
	BankCard     models.BankCard `json:"bankCard" form:"bankCard" query:"bankCard" bson:"bankCard" binding:"required"` //体现银行卡
	CreateBy     string          `json:"createBy" form:"createBy" query:"createBy" bson:"createBy" binding:"-"`        //申请人
	CreateAt     time.Time       `json:"createAt" form:"createAt" query:"createAt" bson:"createAt" binding:"-"`
	UpdateBy     string          `json:"updateBy" form:"updateBy" query:"updateBy" bson:"updateBy" binding:"-"`
	UpdateAt     time.Time       `json:"updateAt" form:"updateAt" query:"updateAt" bson:"updateAt" binding:"-"`
	State        string          `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
	StateExplain string          `json:"stateExplain" form:"stateExplain" query:"stateExplain" bson:"stateExplain" binding:"-"`
}
