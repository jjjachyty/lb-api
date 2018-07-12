package middlewares

import (
	"fmt"
	"lb-api/util"

	"github.com/gin-gonic/gin"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	ACCESS_KEY           = "TvtvJ6rK7JfRmuDQvRfsZZ9zm1JvD_MotrxgT4aN"
	SECRET_KEY           = "noyP3xm9Ibj_Al1CqCHdS7H6aWP4lsg15ps2IZXy"
	AVATAR_BUCKET_NAME   = "luobo"
	ARTICLE_BUCKET_NAME  = "luobo-article"
	PURCHASE_BUCKET_NAME = "4t-purchase"
)

//获取骑牛云上传key
func GetQnToken(c *gin.Context) {
	var err error
	var upToken string
	var scope string
	tokenType := c.Query("type")

	if "" != tokenType {

		switch tokenType {
		case "1":
			scope = AVATAR_BUCKET_NAME + ":" + GetUserIDFromToken(c)
		case "2":
			scope = ARTICLE_BUCKET_NAME
		case "3":
			scope = PURCHASE_BUCKET_NAME
		}

		putPolicy := storage.PutPolicy{
			Scope:      scope,
			InsertOnly: 0,
		}
		mac := qbox.NewMac(ACCESS_KEY, SECRET_KEY)

		upToken = putPolicy.UploadToken(mac)
	} else {
		err = &util.GError{Code: 0, Err: "缺少图片上传方式"}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取上传Key", Data: upToken, Error: err})
}

func DeleteFile(c *gin.Context) {
	var err error
	bucket := c.Query("bucket")
	key := c.Query("key")
	userid := GetUserIDFromToken(c)
	if "" != userid && "" != bucket && "" != key {
		mac := qbox.NewMac(ACCESS_KEY, SECRET_KEY)
		cfg := storage.Config{
			// 是否使用https域名进行资源管理
			UseHTTPS: false,
		}
		// 指定空间所在的区域，如果不指定将自动探测
		// 如果没有特殊需求，默认不需要指定
		//cfg.Zone=&storage.ZoneHuabei
		key = userid[18:] + key
		bucketManager := storage.NewBucketManager(mac, &cfg)
		err = bucketManager.Delete(bucket, key)
	}
	fmt.Println("DeleteFile", bucket, key, userid, err)
	util.JSON(c, util.ResponseMesage{Message: "删除文件", Data: nil, Error: err})

}

func DeleteFiles(bucket string, keys ...string) error {
	var err error

	if "" != bucket && len(keys) > 0 {
		mac := qbox.NewMac(ACCESS_KEY, SECRET_KEY)
		cfg := storage.Config{
			// 是否使用https域名进行资源管理
			UseHTTPS: false,
		}
		// 指定空间所在的区域，如果不指定将自动探测
		// 如果没有特殊需求，默认不需要指定
		//cfg.Zone=&storage.ZoneHuabei

		bucketManager := storage.NewBucketManager(mac, &cfg)
		for _, key := range keys {
			err = bucketManager.Delete(bucket, key)
		}
		util.Glog.Debugf("批量删除图片-bucket%s-keys%v-状态%v", bucket, keys, err)

	}
	return err

}
