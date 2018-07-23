package wx

import (
	"encoding/xml"
	"io/ioutil"
	"lb-api/middlewares"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type WxPayControl struct{}

var l = new(sync.Mutex)

func (WxPayControl) Notify(c *gin.Context) {
	var err error
	var notifyReturn NotifyReturn
	var payment = new(order.Payment)
	notify := new(PayNotify)

	b, _ := ioutil.ReadAll(c.Request.Body)
	util.Glog.Debugln("微信支付回调返回数据", string(b))

	// if err = c.ShouldBind(notify); nil == err {
	if err = xml.Unmarshal(b, notify); nil == err {
		if notify.ReturnCode == "SUCCESS" && notify.ResultCode == "SUCCESS" { //成功

			valid := checkSign(getSign(notify), notify.Sign)
			if valid {
				l.Lock()
				defer l.Unlock()
				//检查订单状态
				err = payment.One(bson.M{"order": notify.Attach, "payType": "pay"}) //查询支付订单
				if payment.State == "0" {                                           //待支付，更新
					// strconv.ParseFloat(strconv.FormatInt(notify.CashFee,10),64)
					if payment.TradeAmount == notify.TotalFee/100 { //支付金额等于订单金额
						err = payment.UpdateOne(bson.M{"$set": bson.M{"state": "1", "payAt": time.Now(), "payAmount": notify.TotalFee / 100}})

						if nil != err {
							util.Glog.Errorf("人工补偿-微信支付回调-1更新系统支付payment-失败-微信返回数据%v", notify)
						} else { //更新订单状态
							err = models.Update("order", bson.M{"_id": bson.ObjectIdHex(payment.Order)}, bson.M{"$set": bson.M{"state": "1"}})
							if nil != err {
								//回退支付状态
								err = payment.UpdateOne(bson.M{"$set": bson.M{"state": "0", "payAt": time.Now(), "payAmount": notify.TotalFee / 100}})
								if nil != err {
									util.Glog.Errorf("人工补偿-微信支付回调-2更新订单order-失败-微信返回数据%v", notify)
								}
							}
						}
					} else {
						err = &util.GError{Code: -1, Err: "订单金额与支付金额不一致"}
						util.Glog.Errorf("微信支付回调-订单金额与支付金额不匹配-失败-订单%d-微信数据%d", payment.TradeAmount, notify.TotalFee)
					}
				}
			} else {
				err = &util.GError{Code: -1, Err: "数据校验失败"}
				util.Glog.Errorf("微信支付回调-数据校验-失败-数据%v", notify)

			}
		} else {
			util.Glog.Errorf("微信支付回调-返回业务-失败-%v", notify)

		}

	} else {

		util.Glog.Errorf("微信支付回调-解析返回数据-失败-%v", c.Request.PostForm)

	}

	if nil != err {
		notifyReturn = NotifyReturn{ReturnCode: "FAIL", ReturnMsg: err.Error()}
	} else {
		notifyReturn = NotifyReturn{ReturnCode: "SUCCESS", ReturnMsg: "OK"}

	}
	c.XML(200, notifyReturn)
}

//获取微信支付号
func (WxPayControl) Pay(c *gin.Context) {
	var err error
	var payment = new(order.Payment)
	var dbOrder = new(order.Order)
	var returnMsg = new(ReturnMsg)
	var totoleFree int64
	var orderID = c.Param("id")
	// var tradeType = c.DefaultQuery("tradetype", "NATIVE")
	if "" != orderID {
		userid := middlewares.GetUserIDFromToken(c)
		//查询订单
		dbOrder.ID = bson.ObjectIdHex(orderID)
		dbOrder.One(bson.M{"_id": bson.ObjectIdHex(orderID)})
		ip := c.ClientIP()
		//判断该订单是否是操作人的
		if dbOrder.Buyer.ID == userid {
			//查询支付表中该订单是否已支付
			payment.Order = orderID
			err = payment.One(bson.M{"order": orderID})
			orderSeq := strings.Split(time.Now().Format("20060102150405.000000"), ".")
			outTradeNo := orderSeq[0] + orderSeq[1]
			if payment.ID.Valid() { //查询到该订单已存在
				if payment.State == "0" { //未支付
					if time.Now().Unix()-payment.CreateAt.Unix() <= 2*60*60 { // 两个小时内重复支付，则直接返回
						returnMsg = &ReturnMsg{ResultCode: "SUCCESS", ReturnCode: "SUCCESS", CodeURL: payment.WxPayURL}
						util.Glog.Debugf("微信支付二维码未过期,还剩%vs", 2*60*60-(time.Now().Unix()-payment.CreateAt.Unix()))

					} else { //两个小时外，已失效,重新发起支付请求
						amount := dbOrder.StrikePrice*100 + systemCharge*100
						totoleFree, err = strconv.ParseInt(strconv.FormatFloat(amount, 'f', 0, 64), 10, 64)
						if nil == err {

							postData := NativePayParams{Appid: config.AppID, Attach: orderID, MchID: config.MchID, NonceStr: util.GetRandomString(6), Body: "4T fortravel.cn 订单支付", OutTradeNO: outTradeNo, TotalFee: totoleFree, SpbillCreateIP: ip, NotifyURL: config.PayNotifyURL, TradeType: "NATIVE"}

							err = requestWx(UnifiedOrderURL, &postData, returnMsg)
							if nil == err { //调用微信支付成功
								if "SUCCESS" == returnMsg.ReturnCode && "SUCCESS" == returnMsg.ResultCode {
									err = payment.UpdateOne(bson.M{"$set": bson.M{"wxPayURL": returnMsg.CodeURL, "createAt": time.Now()}})
									util.Glog.Debugf("微信支付二维码已过期重新生成%v", returnMsg)
								} else {
									err = &util.GError{Code: -1, Err: returnMsg.ReturnMsg + returnMsg.ErrCodeDes}
								}
							}
						} else {
							err = &util.GError{Code: -1, Err: "微信支付支付金额转换失败"}
						}
					}
				} else { //已支付过
					err = &util.GError{Code: -1, Err: "已支付，不要重复支付"}
				}
			} else { //订单不存在，未支付过,新增订单

				amount := dbOrder.StrikePrice*100 + systemCharge*100
				totoleFree, err = strconv.ParseInt(strconv.FormatFloat(amount, 'f', 0, 64), 10, 64)
				if nil == err {
					postData := NativePayParams{Appid: config.AppID, Attach: orderID, MchID: config.MchID, NonceStr: util.GetRandomString(6), Body: "4T fortravel.cn 订单支付", OutTradeNO: outTradeNo, TotalFee: totoleFree, SpbillCreateIP: ip, NotifyURL: config.PayNotifyURL, TradeType: "NATIVE"}
					err = requestWx(UnifiedOrderURL, &postData, returnMsg)
					if nil == err { //调取微信支付成功
						if "SUCCESS" == returnMsg.ReturnCode && "SUCCESS" == returnMsg.ResultCode {

							//新增支付记录
							payment = &order.Payment{ID: bson.NewObjectId(), Order: orderID, PayType: "pay", OutTradeNo: outTradeNo, TradeType: "purchase", CreateAt: time.Now(), TradeAmount: amount / 100, WxPayURL: returnMsg.CodeURL, State: "0"}
							payment.Insert()
						} else {
							err = &util.GError{Code: -1, Err: returnMsg.ReturnMsg + returnMsg.ErrCodeDes}
						}
					}
				} else {
					err = &util.GError{Code: -1, Err: "微信支付支付金额转换失败"}
				}
			}
		} else {
			err = &util.GError{Code: -1, Err: "非法操作,不能操作他人订单,系统已记录此次操作"}
		}
	} else {
		err = &util.GError{Code: -1, Err: "数据完整性错误"}
	}

	util.JSON(c, util.ResponseMesage{Message: "获取微信支付", Data: returnMsg, Error: err})
}
