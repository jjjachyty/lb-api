package controlers

import (
	"lb-api/models"
	"lb-api/util"

	"github.com/gin-gonic/gin"
)

func AvatarList(c *gin.Context) {
	result, err := models.FindAllAvatar()
	util.JSON(c, util.ResponseMesage{Message: "获取头像数据", Data: result, Error: err})
}
