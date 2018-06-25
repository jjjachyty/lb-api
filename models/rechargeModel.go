package models

import (
	"fmt"
	"lb-api/util"
	"time"

	"labix.org/v2/mgo/bson"
)

const (
	rechargeCN = "recharge"
)

type Recharge struct {
	ID          bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id"`
	TradeNumber string        `json:"tradeNumber" form:"tradeNumber" query:"tradeNumber" bson:"tradeNumber"` //交易号
	Type        int           `json:"type" form:"type" query:"type" bson:"type"`                             //转账类型
	Amount      float64       `json:"amount" form:"amount" query:"amount" bson:"amount"`                     //转账金额
	UserID      string        `json:"userID" form:"userID" query:"userID" bson:"userID"`                     //用户号
	Operator    string        `json:"operator" form:"operator" query:"operator" bson:"operator"`             //操作人
	OperatorAt  time.Time     `json:"operatorAt" form:"operatorAt" query:"operatorAt" bson:"operatorAt"`
	State       string        `json:"state" form:"state" query:"state" bson:"state"` //订单状态
	CreateAt    time.Time     `json:"createAt" form:"createAt" query:"createAt" bson:"createAt"`
}

func (rech *Recharge) Insert() error {
	var result []Recharge
	var err error
	rech.ID = bson.NewObjectId()
	result, err = rech.FindAllByCondition(bson.M{"userID": rech.UserID, "state": "0"})
	fmt.Println("result", result)
	if err == nil {
		if len(result) < 3 {
			err = DB.C(rechargeCN).Insert(rech)
		} else {
			err = &util.GError{Code: 4001, Err: "待充值订单过多,请先完成待充值订单"}

		}
	}

	return err
}

// func (rech *Recharge) FindAllByID() ([]Recharge, error) {
// 	var recharges = make([]Recharge, 0)
// 	rech.ID = bson.NewObjectId()
// 	err := DB.C(rechargeCN).Find(bson.M{"userID": rech.UserID}).All(&recharges)
// 	return recharges, err
// }

func (rech *Recharge) FindAllByCondition(condition bson.M) ([]Recharge, error) {
	var recharges = make([]Recharge, 0)
	// rech.ID = bson.NewObjectId()
	err := DB.C(rechargeCN).Find(condition).All(&recharges)
	return recharges, err
}
