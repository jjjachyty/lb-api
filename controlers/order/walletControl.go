package order

import (
	"fmt"
	"io/ioutil"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"labix.org/v2/mgo/bson"
)

type WalletControl struct {
}

const (
	applyCashCN = "applycash"
)

func (WalletControl) GetBankCode(c *gin.Context) {
	var bankCode string
	var bodyByts []byte
	var err error
	var resp *http.Response
	number := c.Param("number")
	if 16 <= len(number) {
		resp, err = http.Get("https://ccdcapi.alipay.com/validateAndCacheCardInfo.json?_input_charset=GBK&cardNo=" + number + "&cardBinCheck=true")
		if nil == err {
			defer resp.Body.Close()
			bodyByts, err = ioutil.ReadAll(resp.Body)
			if nil == err {
				bodyStr := string(bodyByts)
				fmt.Println(gjson.Get(bodyStr, "validated").Bool())
				if gjson.Get(bodyStr, "validated").Bool() {
					bankCode = gjson.Get(bodyStr, "bank").String()
				} else {
					err = &util.GError{Code: -1, Err: "银行卡识别错误"}
				}

			}
		}
	} else {
		err = &util.GError{Code: -1, Err: "银行卡号错误"}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取银行编码", Data: bankCode, Error: err})

}

func (WalletControl) ApplyCash(c *gin.Context) {
	var err error
	var appCash = new(order.ApplyCash)
	var appUser = new(models.User)
	var userid = middlewares.GetUserIDFromToken(c)
	if err = c.ShouldBind(appCash); nil == err {
		err = models.One("user", bson.ObjectIdHex(userid), appUser)
		if "" != appUser.IDCard.Name {
			appCash.ID = bson.NewObjectId()
			appCash.UserName = appUser.IDCard.Name
			appCash.Phone = appUser.Phone
			err = models.Insert(applyCashCN, appCash)
		} else {
			err = &util.GError{Code: -1, Err: "该用户不存在或未实名认证"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "申请提现", Data: nil, Error: err})

}
