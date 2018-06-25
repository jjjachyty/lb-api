package util

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Id           string `json:"id"`
	Phone        string `json:"phone"`
	ForgotPasswd bool   `json:"forgotPasswd"`
	jwt.StandardClaims
}

//GetToken 根据用户名称生成Token
func TokenGen(c *gin.Context, data map[string]interface{}) (string, error) {

	// Set custom claims
	claims := &JwtCustomClaims{
		data["id"].(string),
		data["phone"].(string),
		data["forgotPasswd"].(bool),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("goFuckYouerself!"))
	// if nil == err {
	// 	c.Response().Header().Set("token", t)
	// }
	Glog.Debugf("生成Token-[%v]生成Token%s", data, t)
	return t, err
}
