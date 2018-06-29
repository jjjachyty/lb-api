package controlers

import (
	"lb-api/models"
	"lb-api/util"
	"strconv"

	"labix.org/v2/mgo/bson"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

//web验证验证码
func SendSMS(c *gin.Context) {
	var err error
	var users []models.User
	phone := c.PostForm("phone")
	captchaID := c.PostForm("captchaId")
	captchaSolution := c.PostForm("captchaSolution")
	smsType := c.PostForm("smsType")

	util.Glog.Debugf("短信验证-电话%s,验证码ID%s,验证码%s", phone, captchaID, captchaSolution)
	if phone != "" && captchaID != "" && captchaSolution != "" {
		if CaptchaVaild(captchaID, captchaSolution) { //验证码通过后发送短信
			//检查是否已经注册
			if "1" == smsType { //注册
				users, err = models.User{}.FindAllByCondition(bson.M{"phone": phone})
				if len(users) > 1 { //可以注册
					err = &util.GError{Code: 0, Err: "该手机号已存在，请直接登陆"}
				}
			}
			if nil == err { //可以注册
				var codeStr string
				code := captcha.RandomDigits(4)
				CaptchaStore.Set(phone, code)
				for _, v := range code {
					codeStr += strconv.Itoa(int(v))
				}
				err = util.SendSMS(phone, "84934", codeStr)
			}

		} else {
			err = &util.GError{Code: 1001, Err: "验证码输入错误"}
		}
	} else {
		err = &util.GError{Code: 1002, Err: "数据完整性错误"}

	}
	util.JSON(c, util.ResponseMesage{Message: "短信发送", Data: nil, Error: err})

}
