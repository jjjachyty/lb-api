package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin/binding"

	validator "gopkg.in/go-playground/validator.v8"

	"labix.org/v2/mgo/bson"
)

const (
	exposureArticleCN = "exposureArticle"
)

type Date time.Time

type ExposureArticle struct {
	ID              bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	Title           string        `json:"title" form:"title" query:"title" bson:"title" binding:"required"`                                                                           //文章标题
	Tags            []string      `json:"tags" form:"tags[]" query:"tags"  bson:"tags" binding:"checkTags"`                                                                           //标签
	OccurrenceDate  time.Time     `json:"occurrenceDate" form:"occurrenceDate" query:"occurrenceDate" bson:"occurrenceDate" time_format:"2006-01-02" time_utc:"1" binding:"required"` //发生时间
	Location        string        `json:"location" form:"location" query:"location" bson:"location" binding:"required"`                                                               //发生地
	Domain          string        `json:"domain" form:"taget" query:"domain" bson:"domain" binding:"required"`                                                                        //涉事对象
	Wastage         string        `json:"wastage" form:"wastage" query:"wastage" bson:"wastage" binding:"required"`                                                                   //损失
	Content         string        `json:"content" form:"content" query:"content" bson:"content" binding:"required"`                                                                   //文章内容
	State           string        `json:"state" form:"state" query:"state" bson:"state"`                                                                                              //文章状态
	CreateAt        time.Time     `json:"createAt" form:"-" query:"createAt" bson:"createAt" binding:"-"`
	UpdateAt        time.Time     `json:"updateAt" form:"-" query:"updateAt" bson:"updateAt" binding:"-"`                                         //创建时间
	NickNamePublish bool          `json:"nickNamePublish" form:"nickNamePublish" query:"nickNamePublish" bson:"nickNamePublish" binding:"exists"` //匿名/实名
	CreateUser      string        `json:"createUser" form:"createUser" query:"createUser" bson:"createUser"`                                      //创建人
	NickName        string        `json:"nickName" form:"nickName" query:"nickName" bson:"nickName"`
	Auditor         string        `json:"auditor" form:"auditor" query:"auditor" bson:"auditor"`                     //审核人
	AuditOpinion    string        `json:"auditOpinion" form:"auditOpinion" query:"auditOpinion" bson:"auditOpinion"` //审核意见
	AuditTime       time.Time     `json:"auditTime" form:"-" query:"auditTime" bson:"auditTime" binding:"-"`         //审核时间
	Comments        int           `json:"comments" form:"-" query:"auditTime" bson:"comments" binding:"-"`
	Views           int           `json:"views" form:"-" query:"views" bson:"views" binding:"-"`
	ThumbsUps       int           `json:"thumbsUps" form:"-" query:"thumbsUps" bson:"thumbsUps" binding:"-"`
	ThumbsDowns     int           `json:"thumbsDowns" form:"-" query:"thumbsDowns" bson:"thumbsDowns" binding:"-"`
}

func checkTags(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	fmt.Println("checkTags\n\n\\", field.Interface())
	if tags, ok := field.Interface().([]string); ok {
		if len(tags) > 0 {
			return true
		}
	}
	return false
}

//新增文章
func (ea *ExposureArticle) Insert() error {
	// ea.OccurrenceDate = Date(time.Now())
	fmt.Println("\n\nInsert", ea.CreateAt, ea.OccurrenceDate, ea.Tags)
	return DB.C(exposureArticleCN).Insert(ea)
}

//新增文章
// func (ea *ExposureArticle) FindUserArticle(condition bson.M) ([]ExposureArticle, error) {
// 	var articles = make([]ExposureArticle, 0)
// 	//.Select(bson.M{"title": 1, "tags": 1, "occurrenceDate": 1, "location": 1, "taget": 1, "wastage": 1, "nickNamePublish": 1, "content": 1, "state": 1, "auditOpinion": 1})
// 	err := DB.C(exposureArticleCN).Find(condition).Sort("-createAt").Limit(10).All(&articles)
// 	return articles, err
// }

//获取排行榜
func (ea ExposureArticle) Find(sort string, limit int, selectM bson.M, condition bson.M) ([]ExposureArticle, error) {
	var articles = make([]ExposureArticle, 0)
	// Select(bson.M{"title": 1, "tags": 1, "occurrenceDate": 1, "location": 1, "taget": 1, "wastage": 1, "nickNamePublish": 1, "content": 1, "state": 1, "auditOpinion": 1})
	err := DB.C(exposureArticleCN).Find(condition).Select(selectM).Sort(sort).Limit(limit).All(&articles)
	return articles, err
}

//新增文章
func (ea *ExposureArticle) Remove(condition bson.M) error {
	_, err := DB.C(exposureArticleCN).RemoveAll(condition)
	return err
}

//更新文章
func (ea *ExposureArticle) Update() error {
	return DB.C(exposureArticleCN).UpdateId(ea.ID, bson.M{"$set": bson.M{"title": ea.Title, "tags": ea.Tags, "occurrenceDate": ea.OccurrenceDate, "location": ea.Location, "taget": ea.Domain, "wastage": ea.Wastage, "content": ea.Content, "state": "0", "updateAt": time.Now()}})
}

//更新文章
func (ea ExposureArticle) UpdateCond(cond bson.M) error {
	return DB.C(exposureArticleCN).UpdateId(ea.ID, cond)
}
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		fmt.Println("RegisterValidation")
		v.RegisterValidation("checkTags", checkTags)
	}
}
