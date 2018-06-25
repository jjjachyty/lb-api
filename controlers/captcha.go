package controlers

import (
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

var CaptchaStore captcha.Store

func init() {
	CaptchaStore = captcha.NewMemoryStore(100000, time.Minute)
	captcha.SetCustomStore(CaptchaStore)

}
func Captcha(c *gin.Context) {
	id := c.Param("id")
	number := captcha.RandomDigits(4)
	image := captcha.NewImage(id, number, 100, 30)
	CaptchaStore.Set(id, number)

	image.WriteTo(c.Writer)

}

func CaptchaVaild(id, code string) bool {
	return captcha.VerifyString(id, code)
}
