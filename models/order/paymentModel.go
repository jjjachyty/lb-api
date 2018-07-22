package order

import (
	"lb-api/models"
	"time"

	"labix.org/v2/mgo/bson"
)

type Payer struct {
	ID            string
	PayType       string //微信、支付宝、银行卡
	WeiXinAcct    string
	ALiPayAcct    string
	UserName      string
	BankName      string
	BankCard      string
	CreateType    string
	ValidDate     string
	CVV2          string
	ReservedPhone string
}

type Payment struct {
	ID          bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id"`
	OutTradeNo  string        `bson:"outTradeNo"`  //微信订单号
	TradeType   string        `bson:"tradeType"`   //订单类型/1.代购，2.转卖，3为旅拍
	TradeAmount float64       `bson:"tradeAmount"` //订单金额
	PayAmount   float64       `bson:"payAmount"`   //支付金额
	// Payer      Payer
	// Payee      Payer
	PayType  string    `bson:"payType"` //支付、提现、转账
	Order    string    `bson:"order"`   //系统订单号
	CreateAt time.Time `bson:"createAt"`
	PayAt    time.Time `bson:"payAt"`
	WxPayURL string    `bson:"wxPayURL"`
	State    string    `bson:"state"`
}

const paymentCN = "payment"

//新增支付订单
func (pm *Payment) Insert() error {
	// ea.OccurrenceDate = Date(time.Now())
	return models.DB.C(paymentCN).Insert(pm)
}

//查询订单
// func (pm *Payment) One() error {
// 	// ea.OccurrenceDate = Date(time.Now())
// 	return models.DB.C(paymentCN).FindId(pm.ID).One(pm)
// }

//查询订单
func (pm *Payment) One(query interface{}) error {
	// ea.OccurrenceDate = Date(time.Now())
	return models.DB.C(paymentCN).Find(query).One(pm)
}

func (pm *Payment) Remove() error {
	// ea.OccurrenceDate = Date(time.Now())
	return models.DB.C(paymentCN).RemoveId(pm.ID)
}

func (pm *Payment) UpdateOne(update interface{}) error {
	// ea.OccurrenceDate = Date(time.Now())
	return models.DB.C(paymentCN).UpdateId(pm.ID, update)
}
