package controlers

import (
	"lb-api/models"
	"lb-api/util"
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

func (UserControl) NewTipOffs(c *gin.Context) {
	var tips = new(models.TipOff)
	var err error
	if err = c.Bind(tips); nil == err {
		if len(tips.Content) > 10 && len(tips.Content) < 140 {
			tips.ID = bson.NewObjectId()
			tips.CreateAt = time.Now()
			err = tips.Insert()
		} else {
			err = &util.GError{Code: 0, Err: "举报内容需在10～140字之间哦"}
		}

	}
	util.JSON(c, util.ResponseMesage{Message: "举报", Data: nil, Error: err})

}
