package purchase

import (
	"fmt"
	"lb-api/models"
	"lb-api/util"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
	"labix.org/v2/mgo/bson"
)

type DayTime struct {
	time.Time
}

//报价单
type QuotationOrder struct {
	ID         bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	PurchaseID string        `json:"purchaseID" form:"purchaseID" query:"purchaseID" bson:"purchaseID" binding:"required"` //代购单号
	Amount     float64       `json:"amount" form:"amount" query:"amount" bson:"amount" binding:"-"`                        //总金额
	Products   []Product     `json:"products" form:"products[]" query:"products[]" bson:"products" binding:"checkProducts"`
	Charge     float64       `json:"charge" form:"charge" query:"charge" bson:"charge" binding:"required"`        //服务费
	BuyByID    string        `json:"buyByID" form:"buyByID" query:"buyByID" bson:"buyByID" binding:"-"`           //报价人ID
	BuyByName  string        `json:"buyByName" form:"buyByName" query:"buyByName" bson:"buyByName" binding:"-"`   //报价人昵称
	CreateAt   time.Time     `json:"createAt" form:"-" query:"createAt" bson:"createAt" binding:"-"`              //报价时间
	State      string        `json:"state" form:"state" query:"state" bson:"state" binding:"-"`                   //报价单状态
	ExpiryTime DayTime       `json:"expiryTime" form:"-" query:"expiryTime" bson:"expiryTime" binding:"required"` //失效时间
}

func (t *DayTime) UnmarshalJSON(data []byte) (err error) {
	now, err := getTime(string(data[1 : len(data)-1]))
	*t = DayTime{now}
	return
}

const (
	quotationOrderCN = "quotationOrder"
)

func checkProducts(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	fmt.Println("checkProducts\n\n\\", field.Interface())
	var flag = true
	if products, ok := field.Interface().([]Product); ok {
		for _, p := range products {
			if "" == p.ID.Hex() || p.Price <= 0 {
				flag = false
			}
		}
	}
	return flag
}

//新增文章
func (qo *QuotationOrder) Insert() error {
	util.Glog.Debugf("新增代购报价单%v", qo.ExpiryTime)
	return models.DB.C(quotationOrderCN).Insert(qo)
}

func (QuotationOrder) Find(sort string, limit int, selectM bson.M, condition bson.M) ([]QuotationOrder, error) {
	var qos = make([]QuotationOrder, 0)
	query := models.DB.C(quotationOrderCN).Find(condition)
	if "" != sort {
		query = query.Sort(sort)
	}
	if 0 != limit {
		query = query.Limit(limit)
	}
	if len(selectM) > 0 {
		query = query.Select(selectM)
	}
	err := query.All(&qos)
	return qos, err
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("checkProducts", checkProducts)
	}
}

func getTime(src string) (time.Time, error) {
	fmt.Println("srcsrcsrc", src)
	if "" == src {
		return time.Time{}, &util.GError{Code: 0, Err: "截止时间不能为空"}
	}

	n := time.Now()
	year := n.Year()
	month := n.Month()
	day := n.Day()
	hour, _ := strconv.Atoi(strings.Split(src, ":")[0])
	min, _ := strconv.Atoi(strings.Split(src, ":")[1])
	fmt.Println("\n\n\n", hour, min, strings.Split(src, ":")[0], strings.Split(src, ":")[1])
	t := time.Date(year, month, day, hour, min, 0, 0, time.Local)
	fmt.Println("t", t, "n", n)
	if t.Before(n) {
		return time.Time{}, &util.GError{Code: 0, Err: "截止时间不能小于当前时间"}
	}
	return t, nil
}
