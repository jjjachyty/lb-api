package controlers

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/util"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

type UserControl struct {
}

var avatarURL = "http://pa7c4jxbs.bkt.clouddn.com/"

type Register struct {
	Phone         string `json:"phone" form:"phone" query:"phone"`
	Passwd        string `json:"passwd" form:"passwd" query:"passwd"`
	ConfirmPasswd string `json:"confirmPasswd" form:"confirmPasswd" query:"confirmPasswd"`
	SmsCode       string `json:"smsCode" form:"smsCode" query:"smsCode"`
}

// func (userCtl UserControl) LoginIn(c *gin.Context) {
// 	var err error
// 	var token string
// 	user := new(models.User)
// 	if err = c.Bind(user); err == nil {
// 		fmt.Println("err", err)
// 		if err = user.ValidUser(); nil == err {
// 			token, err = util.TokenGen(c, map[string]interface{}{"id": user.ID.Hex(), "phone": user.Phone, "forgotPasswd": false})

// 		} else {
// 			err = &util.GError{Code: 2001, Err: "用户名或密码错误"}
// 		}
// 	}

// 	util.JSON(c, util.ResponseMesage{Message: "用户登陆", Data: echo.Map{"token": token, "user": echo.Map{"id": user.ID.Hex(), "anNickName": user.AnNickName}}, Error: err})
// }

func (userCtl UserControl) Register(c *gin.Context) {
	var err error
	var token string
	var user *models.User
	register := new(Register)

	if err = c.Bind(register); nil == err {
		util.Glog.Debugf("用户注册-手机号%s", register.Phone)

		user = &models.User{IDCard: new(models.IDCard), IDCardValid: false, Phone: register.Phone, Passwd: register.Passwd, Avatar: avatarURL + "1000", NickName: "萝卜" + register.Phone[7:], AnNickName: "1000"}

		if CaptchaVaild(register.Phone, register.SmsCode) { //验证码通过
			user.State = "1" //生效
			user.ID = bson.NewObjectId()
			// fmt.Println("Object iD", u.ID.String())
			user.CreateAt = time.Now()
			// u.ValidCode = util.UUID()
			user.Passwd = util.MD5(user.Passwd)
			err = user.Create()
			if nil == err {
				user.Passwd = ""
				token, err = util.TokenGen(c, map[string]interface{}{"id": user.ID.Hex(), "phone": user.Phone, "forgotPasswd": false})
			} else {
				err = &util.GError{Code: 11000, Err: "该手机号已被注册"}
			}
		} else {
			err = &util.GError{Code: 1001, Err: "短信验证码输入错误"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "用户注册", Data: map[string]interface{}{"token": token, "user": map[string]interface{}{"id": user.ID.Hex(), "anNickName": user.AnNickName}}, Error: err})

}

func (userCtl UserControl) GetUserInfo(c *gin.Context) {
	var err error
	var user = new(models.User)
	// var result map[string]string
	getType := c.Query("type")
	user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
	user.IDCard = new(models.IDCard)
	err = user.GetInfoByID()
	switch getType {
	case "profile": //简要信息
		user = &models.User{ID: user.ID, AnNickName: user.AnNickName, NickName: user.NickName, IDCardValid: user.IDCardValid, Avatar: user.Avatar}
	default:
		//user.Passwd = ""

	}

	util.JSON(c, util.ResponseMesage{Message: "获取用户信息", Data: user, Error: err})
}

func (userCtl UserControl) GetUserProfile(c *gin.Context) {
	var err error
	var user = new(models.User)

	user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
	user.IDCard = new(models.IDCard)
	err = user.GetInfoByID()

	util.JSON(c, util.ResponseMesage{Message: "获取用户信息", Data: user, Error: err})
}

//UpdateByID 更新用户信息
func (userCtl UserControl) UpdateByID(c *gin.Context) {
	var err error
	var user = new(models.User)

	user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
	if err = c.Bind(user); nil == err {
		if "" != user.NickName && "" != user.AnNickName {
			user.Avatar = avatarURL + user.ID.Hex()
			err = user.UpdateByID()
		} else {
			err = &util.GError{Code: 0000, Err: "匿名昵称和自定义昵称不能为空"}
		}

	}
	fmt.Println("err", err)
	util.JSON(c, util.ResponseMesage{Message: "更新用户信息", Data: user, Error: err})

}

func (userCtl UserControl) PhoneLogin(c *gin.Context) {
	var err error
	var token string
	var expire time.Time
	var user = new(models.User)
	var resultUsers []models.User
	register := new(Register)

	if err = c.Bind(register); nil == err {
		util.Glog.Debugf("短信登陆-手机号%s", register.Phone)

		if CaptchaVaild(register.Phone, register.SmsCode) { //验证码通过
			resultUsers, err = models.User{}.FindAllByCondition(bson.M{"phone": register.Phone})
			if nil == err && 1 == len(resultUsers) {
				user = &resultUsers[0]
				data := make(map[string]interface{}, 0)
				data["id"] = user.ID.Hex()
				data["phone"] = user.Phone
				data["forgotPasswd"] = true
				token, expire, err = middlewares.GetJWTToken(user.Phone, data)
			}

		} else {
			err = &util.GError{Code: 1001, Err: "短信验证码验证错误"}
		}

	}

	var result = map[string]interface{}{"token": token, "expire": expire, "user": map[string]interface{}{"id": user.ID.Hex(), "anNickName": user.AnNickName}}

	util.JSON(c, util.ResponseMesage{Message: "短信登陆", Data: result, Error: err})

}

func (userCtl UserControl) ModifyPasswd(c *gin.Context) {
	var user = new(models.User)
	var err error

	orgPasswd := c.PostForm("orgPasswd")
	newPasswd := c.PostForm("newPasswd")
	phone := c.PostForm("phone")
	//验证数据有效性
	if orgPasswd != newPasswd && regexp.MustCompile(`^[A-Za-z\d]{8,16}$`).MatchString(newPasswd) {
		//user.Phone = tokenInfo.Phone
		user.Passwd = orgPasswd
		user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
		if "" != phone { //忘记密码，从手机短信登录，无需验证原密码
			user.Passwd = util.MD5(newPasswd)
			err = user.UpdateByID()
		} else { //需要验证原密码
			if err = user.ValidUser(); nil == err {
				fmt.Println("user--\n\n", user.ID.Hex())

				user.Passwd = util.MD5(newPasswd)
				err = user.UpdateByID()
			} else {
				fmt.Println("\n\n", err)
				err = &util.GError{Code: 5000, Err: "原密码输入错误,请重新输入"}
			}

		}
	} else {
		err = &util.GError{Code: 5001, Err: "新密码与旧密码一致或密码格式不正确"}
	}
	util.JSON(c, util.ResponseMesage{Message: "密码修改", Data: nil, Error: err})

}

func (userCtl UserControl) UpdateByIDAddress(c *gin.Context) {
	var address = new(models.Address)
	var err error

	if err = c.ShouldBind(address); nil == err {
		models.User{}.Update(bson.M{"address.id":address.ID,bson.M{"$set":bson.M{"province":,"city":,"county":,"street"}}})
	}
	util.JSON(c, util.ResponseMesage{Message: "密码修改", Data: nil, Error: err})

}
