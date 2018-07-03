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
type Date struct {
	time.Time
}

//报价单
type QuotationOrder struct {
	ID           bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id" binding:"-"`
	PurchaseID   string        `json:"purchaseID" form:"purchaseID" query:"purchaseID" bson:"purchaseID" binding:"required"` //代购单号
	Amount       float64       `json:"amount" form:"amount" query:"amount" bson:"amount" binding:"-"`                        //总金额
	Products     []Product     `json:"products" form:"products[]" query:"products[]" bson:"products" binding:"checkProducts"`
	Charge       float64       `json:"charge" form:"charge" query:"charge" bson:"charge" binding:"required"`                  //服务费
	BuyByID      string        `json:"buyByID" form:"buyByID" query:"buyByID" bson:"buyByID" binding:"-"`                     //报价人ID
	BuyByName    string        `json:"buyByName" form:"buyByName" query:"buyByName" bson:"buyByName" binding:"-"`             //报价人昵称
	CreateAt     time.Time     `json:"createAt" form:"-" query:"createAt" bson:"createAt" binding:"-"`                        //报价时间
	State        string        `json:"state" form:"state" query:"state" bson:"state" binding:"-"`                             //报价单状态
	ReasonType   string         `json:"reasonType" form:"reasonType" query:"reasonType" bson:"reasonType" binding:"-"`                             //拒绝原因
	RefuseReason string        `json:"refuseReason" form:"refuseReason" query:"refuseReason" bson:"refuseReason" binding:"-"` //拒绝理由
	ExpiryTime   DayTime       `json:"expiryTime" form:"-" query:"expiryTime" bson:"expiryTime" binding:"required"`           //失效时间
	DeliveryTime Date          `json:"deliveryTime" form:"-" query:"deliveryTime" bson:"deliveryTime" binding:"required"`     //失效时间
	AllowRepeat  bool          `json:"allowRepeat" form:"allowRepeat" query:"allowRepeat" bson:"allowRepeat"`                 //是否允许再次报价
}

func (t *DayTime) UnmarshalJSON(data []byte) (err error) {
	now, err := getTime(string(data[1 : len(data)-1]))
	*t = DayTime{now}
	return
}
func (t *Date) UnmarshalJSON(data []byte) (err error) {
	fmt.Println("string(data)", string(data))
	now, err := time.Parse("2006-01-02", string(data[1:len(data)-1]))
	*t = Date{now}
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

//新增文章
func (qo QuotationOrder) Delete() error {
	util.Glog.Debugf("删除报价单%v", qo.ExpiryTime)
	return models.DB.C(quotationOrderCN).RemoveId(qo.ID)
}

//Update 更新代购单
func (QuotationOrder) Update(selector bson.M, update bson.M) error {
	return models.DB.C(quotationOrderCN).Update(selector, update)
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("checkProducts", checkProducts)
	}
}

func getTime(src string) (time.Time, error) {
	fmt.Println("srcsrcsrc", src)
	if "" == src {
		return time.Time{}, &util.GError{Code: 0, Err: "时间不能为空"}
	}

	n := time.Now()
	year := n.Year()
	month := n.Month()
	day := n.Day()
	hour, _ := strconv.Atoi(strings.Split(src, ":")[0])
	min, _ := strconv.Atoi(strings.Split(src, ":")[1])

	t := time.Date(year, month, day, hour, min, 0, 0, time.Local)
	//默认预留10分钟响应
	fn := n.Add(10 * time.Minute)

	if t.Before(fn) {
		return time.Time{}, &util.GError{Code: 0, Err: "截止时间必须大于" + fn.Format("2006-01-02 15:04")}
	}
	return t, nil
}
