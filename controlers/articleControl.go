package controlers

import (
	"fmt"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

func (UserControl) NewExposureArticle(c *gin.Context) {
	var ea = new(models.ExposureArticle)
	var user = new(models.User)
	var err error

	if err = c.ShouldBind(ea); nil == err {
		ea.State = "1" //审批中
		user.ID = bson.ObjectIdHex(middlewares.GetUserIDFromToken(c))
		err = user.GetInfoByID()
		if nil == err { //查到文章发表人
			ea.CreateBy = user.ID.Hex()
			if ea.NickNamePublish { //匿名发表
				ea.NickName = user.AnNickName
			} else { //实名发表//自定义昵称
				ea.NickName = user.NickName
			}
			ea.ID = bson.NewObjectId()
			fmt.Println("xxxxx", ea)
			ea.CreateAt = time.Now()
			err = ea.Insert()
		}

	} else {
		fmt.Println("\n\n\n新增文章", err)
		err = &util.GError{Code: 0000, Err: "表单数据错误,请检查"}
	}
	fmt.Println("\n\n\n新增文章", err)

	util.JSON(c, util.ResponseMesage{Message: "新增曝光文章", Data: nil, Error: err})
}

func (UserControl) GetUserExposureArticles(c *gin.Context) {
	var ea = new(models.ExposureArticle)
	var cond = bson.M{}
	var err error
	var timeCond = bson.M{}
	var beginDate time.Time
	var endDate time.Time
	var title = c.Query("title")
	var place = c.Query("place")
	var taget = c.Query("taget")
	var beginDateStr = c.Query("beginDate")
	var endDateStr = c.Query("endDate")
	if "" != title {
		cond["title"] = bson.M{"$regex": title}
	}
	if "" != place {
		cond["place"] = bson.M{"$regex": place}

	}
	if "" != taget {
		cond["taget"] = bson.M{"$regex": taget}

	}
	//处理开始时间
	if "" != beginDateStr {
		beginDate, err = time.Parse("2006-01-02", beginDateStr)
		if nil == err {
			timeCond["$gte"] = beginDate

		} else {
			err = &util.GError{Code: 0000, Err: "开始时间格式错误"}
		}
	}
	//处理结束时间
	if "" != endDateStr {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if nil == err {
			timeCond["$lte"] = endDate
		} else {
			err = &util.GError{Code: 0000, Err: "结束时间格式错误"}
		}
	}
	if nil == err && ("" != beginDateStr || "" != endDateStr) {
		cond["createAt"] = timeCond

	}

	fmt.Println("\n\n", title, place, taget, beginDateStr, endDateStr, cond)

	cond["createBy"] = middlewares.GetUserIDFromToken(c)
	result, err := ea.Find("-createAt", 10, bson.M{}, cond)
	util.JSON(c, util.ResponseMesage{Message: "获取我的文章", Data: result, Error: err})
}

func (UserControl) GetUserExposureArticle(c *gin.Context) {
	var err error
	var article = new(models.ExposureArticle)
	var userid = middlewares.GetUserIDFromToken(c)
	id := c.Param("id")
	if "" != id {
		models.FindOne(models.ArticleCN, bson.M{"_id": bson.ObjectIdHex(id), "createBy": userid}, article)
		if !article.ID.Valid() {
			err = &util.GError{Code: 0, Err: "文章不存在"}
		}
	} else {
		err = &util.GError{Code: 0, Err: "文章ID号不能为空"}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取我的文章", Data: article, Error: err})

}

func (UserControl) DeleteExposureArticles(c *gin.Context) {
	var err error
	var articleId = c.Query("articleid")
	var ea = new(models.ExposureArticle)
	createUser := middlewares.GetUserIDFromToken(c)
	fmt.Println("delete", articleId, createUser, "" != articleId && "" != createUser)
	if "" != articleId && "" != createUser {
		err = ea.Remove(bson.M{"_id": bson.ObjectIdHex(articleId), "createBy": createUser})
	} else {
		err = &util.GError{Code: 0000, Err: "缺少删除信息"}
	}
	fmt.Println("delete,err", err)
	util.JSON(c, util.ResponseMesage{Message: "删除我的文章", Data: nil, Error: err})
}

func (UserControl) UpdateExposureArticles(c *gin.Context) {
	var err error
	var ea = new(models.ExposureArticle)
	var id = c.Param("id")
	if err = c.ShouldBind(ea); nil == err {
		ea.ID = bson.ObjectIdHex(id)
		ea.State = "0" //重新审核
		err = ea.Update()
	} else {

		err = &util.GError{Code: 0000, Err: "更新数据格式错误"}
	}

	util.JSON(c, util.ResponseMesage{Message: "更新我的文章", Data: nil, Error: err})
}

//用户浏览自己的文章
func (UserControl) GetUserExpArtById(c *gin.Context) {
	var err error
	var ea = new(models.ExposureArticle)
	var expAts []models.ExposureArticle
	var id = c.Query("id")
	if "" != id {
		expAts, err = ea.Find("_id", 10, bson.M{}, bson.M{"_id": bson.ObjectIdHex(id), "createBy": middlewares.GetUserIDFromToken(c)})
		if len(expAts) > 0 {
			*ea = expAts[0]
		} else {
			err = &util.GError{Code: 0000, Err: "文章不存在"}
		}
	} else {
		err = &util.GError{Code: 0000, Err: "文章ID号不能为空"}
	}

	util.JSON(c, util.ResponseMesage{Message: "更新我的文章", Data: ea, Error: err})
}

//获取所以可浏览的文章
func (UserControl) GetExpArt(c *gin.Context) {
	var err error
	var ea = new(models.ExposureArticle)
	var expAts []models.ExposureArticle

	expAts, err = ea.Find("createAt", 10, bson.M{}, bson.M{"state": "1"})

	util.JSON(c, util.ResponseMesage{Message: "获取所有文章", Data: expAts, Error: err})
}

//获取所以可浏览的文章
func (UserControl) GetExpArtById(c *gin.Context) {
	var err error
	var ea = new(models.ExposureArticle)
	var expAts []models.ExposureArticle
	var id = c.Param("id")
	if "" != id && bson.IsObjectIdHex(id) {
		expAts, err = ea.Find("createAt", 10, bson.M{}, bson.M{"state": "1", "_id": bson.ObjectIdHex(id)})
		if len(expAts) > 0 {
			*ea = expAts[0]
			ea.Views++
			// if ea.CreateUser != middlewares.GetUserIDFromToken(c) {
			//非自己查看
			go func() {
				err := ea.UpdateCond(bson.M{"$set": bson.M{"views": ea.Views}})
				util.Glog.Debugf("更新文章浏览数-[%v]原%d,新%d", err, expAts[0].Views, ea.Views)
			}()
			// }

		}
	} else {
		err = &util.GError{Code: 0000, Err: "文章ID号获取错误"}
	}

	util.JSON(c, util.ResponseMesage{Message: "获取文章", Data: ea, Error: err})
}

func GetTopExpArts(c *gin.Context) {
	var err error
	var sort = c.Query("sort")
	var limitstr = c.Query("limit")
	var limit int
	var result []models.ExposureArticle
	if "" != sort && "" != limitstr {
		limit, err = strconv.Atoi(limitstr)

		if err == nil && limit < 11 {
			result, err = models.ExposureArticle{}.Find(sort, limit, bson.M{"auditor": 0, "auditOpinion": 0, "createBy": 0}, bson.M{"$and": []bson.M{bson.M{"state": "1"}, bson.M{"thumbsUps": bson.M{"$gt": 10}}}})
		}
	}
	util.JSON(c, util.ResponseMesage{Message: "获取文章排序", Data: result, Error: err})
}

func SerchExpArts(c *gin.Context) {
	var err error
	var keywords = c.Query("keywords")
	var result []models.ExposureArticle
	// if "" != keywords {
	result, err = models.ExposureArticle{}.Find("createAt", 10, bson.M{"auditor": 0, "auditOpinion": 0, "createBy": 0}, bson.M{"$and": []bson.M{bson.M{"state": "1"}, bson.M{"$or": []bson.M{bson.M{"title": bson.M{"$regex": keywords}}, bson.M{"content": bson.M{"$regex": keywords}}, bson.M{"location": bson.M{"$regex": keywords}}, bson.M{"taget": bson.M{"$regex": keywords}}, bson.M{"tags": bson.M{"$regex": keywords}}}}}})

	// }
	util.JSON(c, util.ResponseMesage{Message: "搜索文章", Data: result, Error: err})
}
