package controlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/util"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

//web验证验证码
func (userCtl UserControl) IDCardOCR(c *gin.Context) {
	var err error
	var user = new(models.User)
	var idCard = new(models.IDCard)
	var idCardText gjson.Result
	var ddd []byte
	cardBase64 := strings.Split(c.PostForm("img"), ",")[1]
	userId := middlewares.GetUserIDFromToken(c)
	//如果Token中检测到用户ID
	if "" != userId {
		user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
		idCardText, err = util.GetIDCardText(cardBase64)
		if nil == err {
			//图片识别成功,保存文件
			ddd, err = base64.StdEncoding.DecodeString(cardBase64) //成图片文件并把文件写入到buffer
			if nil == err {
				err = ioutil.WriteFile("./assets/idcard/"+userId+"_"+idCardText.Get("side").String()+".jpg", ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件

			}
			json.Unmarshal([]byte(idCardText.String()), &idCard)
			user.IDCardValid = false
			user.IDCard = idCard
			//如果为身份证背面，则检查生效日期
			if "back" == idCardText.Get("side").String() {
				endDate, _ := time.Parse("2006.01.02", strings.Split(idCard.ValidDate, "-")[1])
				fmt.Println("失效日期", endDate, time.Now())
				if time.Now().Before(endDate) { // 身份证有效
					user.IDCardValid = true
				} else {
					err = &util.GError{Code: 3002, Err: "生份证已失效,请重新上传身份证"}
				}
			}

			//如果没有错误则更新身份信息
			if nil == err {
				err = user.UpdateIdCard()
			}

		}

	}

	util.JSON(c, util.ResponseMesage{Message: "身份证信息识别", Data: user, Error: err})
}

//ValidIDCard 更新身份验证信息
func (userCtl UserControl) ValidIDCard(c *gin.Context) {
	var err error
	user := new(models.User)
	idCard := new(models.IDCard)

	if err = c.Bind(idCard); nil == err {
		user.IDCard = idCard
		user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
		util.Glog.Debug("身份认证-信息%v", user.IDCard)
		fmt.Println("身份认证-", user.IDCard)
		if "" != user.IDCard.IdCardNumber && "" != user.IDCard.ValidDate {
			err = user.VaildIDCard()
		} else {
			err = &util.GError{Code: 3004, Err: "身份认证信息不完整，请重新认证"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "身份验证成功", Data: user, Error: err})

}
