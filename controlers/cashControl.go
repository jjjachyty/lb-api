package controlers

// func (UserControl) AllApplyCash(c *gin.Context) {
// 	var wc = &models.WithdrawCash{}
// 	orders, err := wc.FindAllByCondition(bson.M{"userID": middlewares.GetUserIDFromToken(c)})
// 	util.JSON(c, util.ResponseMesage{Message: "获取提现申请单", Data: orders, Error: err})

// }

// func (UserControl) NewApplyCash(c *gin.Context) {
// 	var wc = new(models.WithdrawCash)
// 	var err error
// 	if err = c.Bind(wc); nil == err {
// 		if "" != wc.BankName && "" != wc.CardNumber && 0 < wc.Amount && 0 == int(wc.Amount)%100 {
// 			err = wc.Insert()
// 		} else {
// 			err = &util.GError{Code: 4004, Err: "提现表单数据不完整，请检查银行名/卡号/金额"}
// 		}

// 	}
// 	util.JSON(c, util.ResponseMesage{Message: "申请提现", Data: nil, Error: err})

// }
