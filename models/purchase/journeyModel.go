package purchase

import (
	"lb-api/models"
	"time"

	"labix.org/v2/mgo/bson"
)

type Journey struct {
	ID          bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id,omitempty" binding:"-"`
	StartDate   time.Time     `json:"startDate" form:"startDate" query:"startDate" bson:"startDate" binding:"required" time_format:"2006-01-02" time_utc:"1" `
	EndDate     time.Time     `json:"endDate" form:"endDate" query:"endDate" bson:"endDate" binding:"required" time_format:"2006-01-02" time_utc:"1" `
	Destination string        `json:"destination" form:"destination" query:"destination" bson:"destination" binding:"required"`
	Products    []string      `json:"products" form:"products[]" query:"products" bson:"products" binding:"required"`
	Remarks     string        `json:"remarks" form:"remarks" query:"remarks" bson:"remarks" binding:"-"`
	ChargeType  string        `json:"chargeType" form:"chargeType" query:"chargeType" bson:"chargeType" binding:"-"`
	ChargeValue float64       `json:"chargeValue" form:"chargeValue" query:"chargeValue" bson:"chargeValue" binding:"-"`
	CreateBy    string        `json:"createBy" form:"createBy" query:"createBy" bson:"createBy" binding:"-"`
	CreateAt    time.Time     `json:"createAt" form:"-" query:"createAt" bson:"createAt" binding:"-"`
	UpdateAt    time.Time     `json:"updateAt" form:"-" query:"updateAt" bson:"updateAt" binding:"-"`
	State       string        `json:"state" form:"state" query:"state" bson:"state" binding:"-"`
}

const (
	journeyCN = "journey"
)

func (Journey) Find(sort string, limit int, selectM bson.M, condition bson.M) ([]Journey, error) {
	var journey = make([]Journey, 0)
	query := models.DB.C(journeyCN).Find(condition)
	if "" != sort {
		query = query.Sort(sort)
	}
	if 0 != limit {
		query = query.Limit(limit)
	}
	if len(selectM) > 0 {
		query = query.Select(selectM)
	}
	err := query.All(&journey)
	return journey, err
}

//Update 更新代购单
func (Journey) Update(selector bson.M, update bson.M) error {
	return models.DB.C(journeyCN).Update(selector, update)
}

//Update 更新代购单
func (p *Journey) Insert() error {
	return models.DB.C(journeyCN).Insert(p)
}
