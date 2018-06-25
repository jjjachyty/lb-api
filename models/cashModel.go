package models

import (
	"lb-api/util"
	"time"

	"labix.org/v2/mgo/bson"
)

const (
	withdrawCashCN = "recharge"
)

type WithdrawCash struct {
	ID         bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id"`
	Amount     float64       `json:"amount" form:"amount" query:"amount" bson:"amount"` //转账金额
	BankName   string        `json:"bankName" form:"bankName" query:"bankName" bson:"bankName"`
	BankOrgin  string        `json:"bankOrgin" form:"bankOrgin" query:"bankOrgin" bson:"bankOrgin"`
	CardNumber string        `json:"cardNumber" form:"cardNumber" query:"cardNumber" bson:"cardNumber"`
	ApplyTime  time.Time     `json:"applyTime" form:"applyTime" query:"applyTime" bson:"applyTime"`
	UserID     string        `json:"userID" form:"userID" query:"userID" bson:"userID"`         //用户号
	UserName   string        `json:"userName" form:"userName" query:"userName" bson:"userName"` //用户
	Operator   string        `json:"operator" form:"operator" query:"operator" bson:"operator"` //操作人
	OperatorAt time.Time     `json:"operatorAt" form:"operatorAt" query:"operatorAt" bson:"operatorAt"`
	State      string        `json:"state" form:"state" query:"state" bson:"state"` //订单状态
}

func (wc *WithdrawCash) Insert() error {
	var result []WithdrawCash
	var err error
	wc.ID = bson.NewObjectId()
	result, err = wc.FindAllByCondition(bson.M{"userID": wc.UserID, "state": "0"})
	if err == nil {
		if len(result) < 3 {
			user := &User{ID: bson.ObjectIdHex(wc.UserID)}
			user.GetInfoByID()
			if user.AvailableBond > wc.Amount {
				//手动实现事物的提交和回滚
				err = DB.C(withdrawCashCN).Insert(wc)
				if nil == err {
					//用户表减去可用保证金
					user.AvailableBond -= wc.Amount
					//更新用户表的可用保证金
					err = user.Update()
					if nil != err {
						//如果更新用户表失败,那么回滚申请提现
						err = wc.RemoveByID()
						if nil != err { //事物回滚失败，日志记录
							util.Glog.Fatalf("事物回滚失败-recharge表_id%s", wc.ID)
						}
						err = &util.GError{Code: 4002, Err: "提现申请失败,请稍后再试"}
					}
				}

			}
		} else {
			err = &util.GError{Code: 4001, Err: "待取现订单过多,请先等待已申请的取现订单"}

		}
	}

	return err
}

// func (wc *WithdrawCash) FindAllByID() ([]WithdrawCash, error) {
// 	var withdrawCashs = make([]WithdrawCash, 0)
// 	err := DB.C(withdrawCashCN).Find(bson.M{"userID": wc.UserID}).All(&withdrawCashs)
// 	return withdrawCashs, err
// }
func (wc *WithdrawCash) RemoveByID() error {
	return DB.C(withdrawCashCN).RemoveId(wc.ID)
}

func (wc *WithdrawCash) FindAllByCondition(condition bson.M) ([]WithdrawCash, error) {
	var withdrawCashs = make([]WithdrawCash, 0)
	err := DB.C(withdrawCashCN).Find(condition).All(&withdrawCashs)
	return withdrawCashs, err
}
