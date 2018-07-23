package wx

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lb-api/models/order"
	"lb-api/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
)

type WxRefundControl struct{}

//Refund func
//outTradeNo 订单号 outRefundNo系统内部订单号orderid,totalFee订单金额 refundFee退款金额
func (WxRefundControl) Refund(outTradeNo, outRefundNo string, totalFee, refundFee int64, refundDesc string) error {
	var refundRt = new(RefundReturn)
	var returnAmount float64
	var params = RefundParams{Appid: config.AppID, MchID: config.MchID, NonceStr: util.GetRandomString(6), NotifyURL: config.RefundNotifyURL, OutTradeNo: outTradeNo, OutRefundNo: outRefundNo, TotalFee: totalFee, RefundFee: refundFee, RefundDesc: refundDesc}
	err := requestWxRefund(ReFundURL, &params, refundRt)
	if nil == err { //退款申请成功
		returnAmount, err = strconv.ParseFloat(strconv.FormatInt(params.RefundFee, 10), 64)
		if nil == err {
			payment := order.Payment{ID: bson.NewObjectId(), OutTradeNo: params.OutTradeNo, Order: params.OutRefundNo, PayType: "return", CreateAt: time.Now(), TradeAmount: returnAmount}
			err = payment.Insert()
			if nil != err {
				util.Glog.Errorf("人工补偿-微信退款-新增系统退款记录-失败-信息%s-记录%v", err.Error(), payment)
			}
		} else {
			err = &util.GError{Code: -1, Err: "退款金额错误"}
		}
	}
	return err
}

//微信退款通知
func (WxRefundControl) Notify(c *gin.Context) {
	var err error
	var notifyReturn NotifyReturn
	// var payment = new(order.Payment)
	notify := new(RefundNotify)

	b, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("微信退款返回数据", string(b))

	// if err = c.ShouldBind(notify); nil == err {
	if err = xml.Unmarshal(b, notify); nil == err {
		if notify.ReturnCode == "SUCCESS" { //成功
			encodeString, err := base64.StdEncoding.DecodeString(notify.ReqInfo)
			fmt.Println("微信退款BASE64解密后数据", "原来", notify.ReqInfo, "BASE64", string(encodeString))
			md5MchKey := util.MD5(config.ApiKey)
			reqInfo, err := util.AESDecrypt(encodeString, []byte(md5MchKey))
			if nil == err {
				fmt.Println("微信退款解密后数据", string(reqInfo))
			} else {
				fmt.Println("微信退款解密后数据", string(reqInfo), "错误", err.Error())
			}

		}
		if nil == err {
			notifyReturn = NotifyReturn{ReturnCode: "SUCCESS", ReturnMsg: "OK"}
		} else {
			notifyReturn = NotifyReturn{ReturnCode: "FAIL", ReturnMsg: err.Error()}
		}
	}
	c.XML(200, notifyReturn)
}
