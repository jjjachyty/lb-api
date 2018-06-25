package controlers

import (
	"bytes"
	"html/template"
	"lb-api/util"

	gomail "gopkg.in/gomail.v2"
)

type EmailValidStrc struct {
	UserID     string
	EmailValid string
}

var XXX_MAIL_TEMPLATE = `
    <div><a style='font-size: 20px;padding-left:200px' href='http://127.0.0.1:8080/emailvalid/{{.UserID}}/{{.EmailValid}}'>欢迎注册,点击我去激活吧</a></div>`

func SendMail(userid string, email string, emailValid string) {

	MAIL_TEMPLATE := XXX_MAIL_TEMPLATE
	m := gomail.NewMessage()
	m.SetHeader("From", "332642088@qq.com")
	m.SetHeader("To", email) //send email to multipul persons
	m.SetHeader("Subject", "[萝卜网]注册验证邮箱")
	t, err := template.New("mail template").Parse(MAIL_TEMPLATE)
	if err != nil {
		util.Glog.Errorf("邮箱发送失败-HTML模版解析错误,%v", err)
	}
	buffer := new(bytes.Buffer)
	t.Execute(buffer, EmailValidStrc{UserID: userid, EmailValid: emailValid})
	m.SetBody("text/html", buffer.String())

	d := gomail.Dialer{Host: "smtp.qq.com", Port: 465, Username: "332642088@qq.com", Password: "zfigkyvbymrcbjhf", SSL: true}
	if err := d.DialAndSend(m); err != nil {

		util.Glog.Errorf("邮箱发送失败,%v", err)
	}
}

// //验证邮箱
// func ValidEmail(c *gin.Context) error {
// 	var err error
// 	id := c.FormValue("id")
// 	code := c.FormValue("code")
// 	fmt.Println("id", id, "code", code)
// 	user := new(models.User)
// 	// user.ValidCode = code
// 	if bson.IsObjectIdHex(id) {
// 		user.ID = bson.ObjectIdHex(id)
// 		// err = user.EmailValid()
// 	} else {
// 		err = &util.GError{Code: 3001, Err: "邮箱验证错误,数据格式不正确,请重新认证"}
// 	}

// 	return util.JSON(c, util.ResponseMesage{Message: "邮箱验证", Data: nil, Error: err})
// }
