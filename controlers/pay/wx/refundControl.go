package wx

import (
	"encoding/base64"
	"encoding/xml"
	"lb-api/models"
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
	var err error
	var returnAmount float64
	var params = RefundParams{Appid: config.AppID, MchID: config.MchID, NonceStr: util.GetRandomString(6), NotifyURL: config.RefundNotifyURL, OutTradeNo: outTradeNo, OutRefundNo: outRefundNo, TotalFee: totalFee, RefundFee: refundFee, RefundDesc: refundDesc}

	returnAmount, err = strconv.ParseFloat(strconv.FormatInt(params.RefundFee, 10), 64)
	if nil == err {
		payment := order.Payment{ID: bson.NewObjectId(), OutTradeNo: params.OutTradeNo, Order: params.OutRefundNo, PayType: "return", CreateAt: time.Now(), TradeAmount: returnAmount, State: "-1"}

		err = requestWxRefund(ReFundURL, &params, refundRt)
		if "SUCCESS" == refundRt.ReturnCode && "SUCCESS" == refundRt.ResultCode { //退款申请成功
			util.Glog.Debugf("微信退款申请-申请微信退款-成功-订单号%s", payment.Order)

			//修改系统订单为申请成功
			err = models.Update("payment", bson.M{"order": outRefundNo, "payType": "return"}, bson.M{"$set": bson.M{"state": "0"}})
			if nil != err {
				util.Glog.Errorf("人工补偿-微信退款申请-更新支付状态为[申请中]-失败-%s", err.Error())
			}
		} else { //退款申请失败
			err = &util.GError{Code: -1, Err: refundRt.ReturnMsg + refundRt.ErrCodeDes}
			util.Glog.Errorf("人工补偿-微信退款申请-申请微信退款失败-失败-%s", refundRt.ReturnMsg+refundRt.ErrCodeDes)
		}

	} else {
		err = &util.GError{Code: -1, Err: "退款金额错误"}
	}

	return err
}

//微信退款通知
func (WxRefundControl) Notify(c *gin.Context) {
	var err error
	var notifyReturn NotifyReturn
	// var payment = new(order.Payment)
	notify := new(RefundNotify)
	notifyDecrypt := new(RefundNotifyDecrypt)
	// b, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println("微信退款返回数据", string(b))

	if err = c.ShouldBind(notify); nil == err {
		// if err = xml.Unmarshal(b, notify); nil == err {
		if notify.ReturnCode == "SUCCESS" { //成功
			encodeString, err := base64.StdEncoding.DecodeString(notify.ReqInfo)
			md5MchKey := util.MD5(config.ApiKey)
			reqInfo, err := util.AESDecrypt(encodeString, []byte(md5MchKey))
			if nil == err {
				err = xml.Unmarshal(reqInfo, notifyDecrypt)
				if nil == err {
					if "SUCCESS" == notifyDecrypt.RefundStatus {
						err = models.Update("payment", bson.M{"order": notifyDecrypt.OutRefundNO, "payType": "return"}, bson.M{"$set": bson.M{"state": "1"}})
						if nil == err { //更新订单状态为退款成功
							err = models.Update("order", bson.M{"_id": bson.ObjectIdHex(notifyDecrypt.OutRefundNO)}, bson.M{"$set": bson.M{"state": "-11"}})
						}
					}
				} else {
					util.Glog.Errorf("微信解密后转换RefundNotifyDecrypt错误-%s", string(reqInfo))
				}
			} else {
				util.Glog.Errorf("微信解密错误-错误信息%s-%s", err.Error(), encodeString)
			}

		}

	}
	if nil == err {
		notifyReturn = NotifyReturn{ReturnCode: "SUCCESS", ReturnMsg: "OK"}
	} else {
		notifyReturn = NotifyReturn{ReturnCode: "FAIL", ReturnMsg: err.Error()}
	}
	c.XML(200, notifyReturn)
}
