package models

import (
	"fmt"
	"lb-api/util"
	"reflect"
	"time"

	"labix.org/v2/mgo/bson"
)

const (
	collectionName = "user"
)

type IDCard struct {
	Name         string `json:"name" form:"name" query:"name"`                                                 //名字
	Gender       string `json:"gender" form:"gender" query:"gender"`                                           //性别
	IdCardNumber string `json:"id_card_number" form:"id_card_number" query:"idCardNumber" bson:"idCardNumber"` //身份证号
	Birthday     string `json:"birthday" form:"birthday" query:"birthday"`                                     //生日
	Race         string `json:"race" form:"race" query:"race"`                                                 //名族
	Address      string `json:"address" form:"address" query:"address"`                                        //地址
	// BeginDate time.Time `json:"beginDate" form:"beginDate" query:"beginDate"`
	// EndDate   time.Time `json:"endDate" form:"endDate" query:"endDate"`    //有效日期
	ValidDate string `json:"valid_date" form:"valid_date" query:"validDate" bson:"validDate"` //有效日期
	IssuedBy  string `json:"issued_by" form:"issued_by" query:"issuedBy" bson:"issuedBy"`     //签发机关
}
type User struct {
	ID            bson.ObjectId `json:"id" form:"id" query:"id" bson:"_id"`
	UserName      string        `json:"userName" form:"userName" query:"userName"`
	AnNickName    string        `json:"anNickName" form:"anNickName" query:"anNickName"`
	NickName      string        `json:"nickName" form:"nickName" query:"nickName"`
	Email         string        `json:"email" form:"email" query:"email"`
	Passwd        string        `json:"-" form:"passwd" query:"passwd"`
	Phone         string        `json:"phone" form:"phone" query:"phone"`
	IDCard        *IDCard       `json:"idCard" form:"idCard" query:"idCard" bson:"idCard"`
	IDCardValid   bool          `json:"idCardValid" form:"idCardValid" query:"idCardValid"`
	Avatar        string        `json:"avatar" form:"avatar" query:"avatar"`
	Location      string        `json:"location" form:"location" query:"location"`
	Address       string        `json:"address" form:"address" query:"address"`
	State         string        `json:"state" form:"state" query:"state"`                                              //用户状态
	Bond          float64       `json:"bond" form:"bond" query:"bond"`                                                 //保证金
	AvailableBond float64       `json:"availableBond" form:"availableBond" query:"availableBond" bson:"availableBond"` //可用保证金
	// ValidCode string        `json:"validCode" form:"validCode" query:"validCode"`
	CreateAt time.Time `json:"createAt" form:"createAt" query:"createAt"`
}

func (u *User) ValidUser() error {
	return DB.C(collectionName).Find(bson.M{"phone": u.Phone, "passwd": util.MD5(u.Passwd), "state": bson.M{"$ne": "-1"}}).Select(bson.M{"_id": 1, "annickname": 1}).One(u)
}
func (u *User) GetInfoByID() error {
	return DB.C(collectionName).FindId(u.ID).One(&u)
}
func (u *User) Create() error {
	return DB.C(collectionName).Insert(u)
}

func (u *User) UpdateIdCard() error {
	var idCard bson.M
	if "" == (u.IDCard.IssuedBy) {
		idCard = bson.M{"idCard.name": u.IDCard.Name, "idCard.gender": u.IDCard.Gender, "idCard.idCardNumber": u.IDCard.IdCardNumber, "idCard.birthday": u.IDCard.Birthday, "idCard.race": u.IDCard.Race, "idCard.address": u.IDCard.Address}
	} else {
		idCard = bson.M{"idCard.validDate": u.IDCard.ValidDate, "idCard.issuedBy": u.IDCard.IssuedBy}
	}

	return DB.C(collectionName).UpdateId(u.ID, bson.M{"$set": idCard})
}

// Update 修改用户信息
func (u *User) Update() error {
	var updateValue bson.M
	if "" != u.NickName { //更新基本信息
		updateValue = bson.M{"annickname": u.AnNickName, "nickname": u.NickName, "address": u.Address, "avatar": u.Avatar}
	} else if 0 != u.AvailableBond { //更新可用保证金
		updateValue = bson.M{"availableBond": u.AvailableBond}
	} else if "" != u.Passwd { //修改密码
		updateValue = bson.M{"passwd": u.Passwd}
	}

	return DB.C(collectionName).UpdateId(u.ID, bson.M{"$set": updateValue})
}

// //验证用户邮箱
// func (u *User) EmailValid() error {
// 	var err error
// 	err = DB.C(collectionName).Find(bson.M{"_id": u.ID, "validcode": u.ValidCode}).One(&u)
// 	fmt.Println("1-----", u, err)
// 	if nil == err {
// 		err = DB.C(collectionName).UpdateId(u.ID, bson.M{"$set": bson.M{"validcode": ""}})
// 	} else {
// 		return &util.GError{Code: 3001, Err: "用户不存在或验证码不存在，请重新认证"}
// 	}
// 	return err
// }

func (u *User) VaildIDCard() error {
	var dbUser = new(User)
	err := DB.C(collectionName).FindId(u.ID).One(dbUser)
	fmt.Println("\n\n", dbUser.IDCard, u.IDCard)
	if nil == err {
		if reflect.DeepEqual(u.IDCard, dbUser.IDCard) {
			util.Glog.Debugf("身份验证-更新%s验证状态", u.ID.Hex())
			err = DB.C(collectionName).UpdateId(u.ID, bson.M{"$set": bson.M{"idcardvalid": true}})
		} else {
			err = &util.GError{Code: 3003, Err: "验证数据有误，请重新验证"}
		}
	}
	return err
}

func (u User) FindAllByCondition(condition bson.M) ([]User, error) {
	var users = make([]User, 0)
	// rech.ID = bson.NewObjectId()
	err := DB.C(collectionName).Find(condition).All(&users)
	return users, err
}