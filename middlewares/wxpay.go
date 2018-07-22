package middlewares

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lb-api/models"
	"lb-api/models/order"
	"lb-api/util"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/gin-gonic/gin"
)

var UnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

//系统收取手续费
const systemCharge = 0

var config = struct {
	Appid      string
	Mch_id     string
	Apikey     string
	Notify_url string
}{
	Appid:      "wxb5f20ab8f1933772",
	Mch_id:     "1501291661",
	Apikey:     "jjjachyty929133jjjachyty92913364",
	Notify_url: "https://www.fortravel.cn/api/v1/pay/wx/notifiy",
}

type ReturnMsg struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

type UnifiedOrderParams struct {
	Appid            string `xml:"appid"` //公众账号ID
	Attach           string `xml:"attach"`
	Mch_id           string `xml:"mch_id"` //商户号
	device_ifno      string //设备号,非必填
	Nonce_str        string `xml:"nonce_str"` //随机字符串
	Sign             string `xml:"sign"`      //签名
	sign_type        string //签名类型，默认为MD5，支持HMAC-SHA256和MD5。非必填
	Body             string `xml:"body"` //商品描述
	detail           string //商品详情 ,非必填
	attach           string //附加数据,非必填
	Out_trade_no     string `xml:"out_trade_no"` //商户订单号
	fee_type         string //标价币种,非必填 默认CNy
	Total_fee        int64  `xml:"total_fee"`        //标价金额
	Spbill_create_ip string `xml:"spbill_create_ip"` //终端IP
	time_start       string //交易起始时间,非必填
	time_expire      string //交易结束时间,非必填
	goods_tag        string //订单优惠标记,非必填
	Notify_url       string `xml:"notify_url"` //通知地址
	Trade_type       string `xml:"trade_type"` //交易类型
	product_id       string //非必填,trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	limit_pay        string //非必填,上传此参数no_credit--可限制用户不能使用信用卡支付
	openid           string //用户标识	trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识

}

type Notify struct {
	Appid         string  `xml:"appid"`
	Attach        string  `xml:"attach"`    //
	BankType      string  `xml:"bank_type"` //付款银行
	CashFee       float64 `xml:"cash_fee"`  //现金支付金额
	FeeType       string  `xml:"fee_type"`  //现金支付金额
	IsSubscribe   string  `xml:"is_subscribe"`
	MchID         string  `xml:"mch_id"`
	NonceStr      string  `xml:"nonce_str"`
	OpenID        string  `xml:"openid"` //用户标识
	OutTradeNo    string  `xml:"out_trade_no"`
	ResultCode    string  `xml:"result_code"`
	ReturnCode    string  `xml:"return_code"`
	Sign          string  `xml:"sign"`
	TimeEnd       string  `xml:"time_end"`
	TotalFee      float64 `xml:"total_fee"`
	TradeType     string  `xml:"trade_type"`
	TransactionID string  `xml:"transaction_id"` //微信支付订单号

	// ReturnMsg  string `xml:"return_msg"`
	// ErrCode    string `xml:"err_code"`
	// ErrCodeDes string `xml:"err_code_des"`
}

type NotifyReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

func NotifyCallBack(c *gin.Context) {
	var err error
	var notifyReturn NotifyReturn
	var payment = new(order.Payment)
	notify := new(Notify)

	if err = c.ShouldBind(notify); nil == err {

		if notify.ReturnCode == "SUCCESS" && notify.ResultCode == "SUCCESS" { //成功

			valid := checkDataVaild(notify)
			if valid {
				//检查订单状态
				err = payment.One(bson.M{"order": notify.Attach})
				if payment.State == "0" { //待支付，更新
					// strconv.ParseFloat(strconv.FormatInt(notify.CashFee,10),64)
					if payment.TradeAmount == notify.TotalFee/100 { //支付金额等于订单金额
						err = payment.UpdateOne(bson.M{"$set": bson.M{"state": "1", "payAt": time.Now(), "payAmount": notify.TotalFee / 100}})

						if nil != err {
							util.Glog.Errorf("人工补偿-微信回调-1更新系统支付-失败-微信返回数据%v", notify)
						} else { //更新订单状态
							err = models.Update("order", bson.M{"_id": bson.ObjectIdHex(payment.Order)}, bson.M{"$set": bson.M{"state": "1"}})
							if nil != err {
								//回退支付状态
								err = payment.UpdateOne(bson.M{"$set": bson.M{"state": "0", "payAt": time.Now(), "payAmount": notify.TotalFee / 100}})
								if nil != err {
									util.Glog.Errorf("人工补偿-微信回调-2更新订单-失败-微信返回数据%v", notify)
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
func checkDataVaild(nt *Notify) bool {
	var checkStr = ""
	rtType := reflect.TypeOf(nt)
	rtValue := reflect.ValueOf(nt)

	for i := 0; i < rtType.Elem().NumField(); i++ {
		xmltage := rtType.Elem().Field(i).Tag.Get("xml")
		if "" != xmltage {
			var fieldValueStr = ""
			var rtFieldValue = rtValue.Elem().Field(i).Interface()

			switch rtFieldValue.(type) {
			case float64:
				fieldValueStr = strconv.FormatFloat(rtFieldValue.(float64), 'f', 0, 64)
			case string:
				fieldValueStr = rtFieldValue.(string)
			}
			if "sgin" != xmltage { //签名不用

				if i > 0 {
					checkStr += "&" + xmltage + "=" + fieldValueStr
				} else {
					checkStr += xmltage + "=" + fieldValueStr
				}
			}
		}
	}
	checkStr += "&key=" + config.Apikey
	fmt.Println("MD5加密串", checkStr)
	if strings.ToUpper(util.MD5(checkStr)) == nt.Sign {
		return true
	}
	return false
}

//获取微信支付号
func WxPay(c *gin.Context) {
	var err error
	var payment = new(order.Payment)
	var dbOrder = new(order.Order)
	var returnMsg = new(ReturnMsg)
	var totoleFree int64
	var orderID = c.Param("id")
	var tradeType = c.Query("tradetype")
	if "" != orderID {
		userid := GetUserIDFromToken(c)
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
					if payment.CreateAt.Unix()-time.Now().Unix() > 2*60*60 { // 两个小时内重复支付，则直接返回
						returnMsg = &ReturnMsg{ResultCode: "SUCCESS", ReturnCode: "SUCCESS", CodeURL: payment.WxPayURL}
						util.Glog.Debugf("微信支付二维码未过期,还剩%vs", payment.CreateAt.Unix()-time.Now().Unix()-2*60*60)

					} else { //两个小时外，已失效,重新发起支付请求
						amount := dbOrder.StrikePrice*100 + systemCharge*100
						totoleFree, err = strconv.ParseInt(strconv.FormatFloat(amount, 'f', 0, 64), 10, 64)
						if nil == err {
							returnMsg, err = requestWxPay(payment.Order, outTradeNo, "4T_fortravel.cn-订单支付", totoleFree, ip, tradeType)
							if nil == err { //调用微信支付成功
								err = payment.UpdateOne(bson.M{"$set": bson.M{"wxPayURL": returnMsg.CodeURL, "createAt": time.Now()}})
								util.Glog.Debugf("微信支付二维码已过期重新生成%v", returnMsg)

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
					returnMsg, err = requestWxPay(orderID, outTradeNo, "4T_fortravel.cn-订单支付", totoleFree, ip, tradeType)
					if nil == err { //调取微信支付成功
						//新增支付记录
						payment = &order.Payment{ID: bson.NewObjectId(), Order: orderID, PayType: "pay", OutTradeNo: outTradeNo, TradeType: "purchase", CreateAt: time.Now(), TradeAmount: amount / 100, WxPayURL: returnMsg.CodeURL, State: "0"}
						payment.Insert()
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

func requestWxPay(attach string, out_trade_no, body string, total_fee int64, spbill_create_ip string, trade_type string) (returnMsg *ReturnMsg, err error) {

	// amount := order.StrikePrice * 100 + systemCharge*100
	// totoleFree, err := strconv.ParseInt(strconv.FormatFloat(amount, 'f', 0, 64), 10, 64)
	// if nil == err {
	// 	//timeseq := strings.Split(time.Now().Format("20060102150405.000000"), ".")
	//fmt.Println("订单号", timeseq)
	sendData := getxml(attach, out_trade_no, body, total_fee, spbill_create_ip, trade_type)
	resp, err := http.Post(UnifiedOrderURL, "application/xml", strings.NewReader(sendData))
	if nil == err {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if nil == err {
			err = xml.Unmarshal(body, &returnMsg)
			util.Glog.Debugf("微信支付返回信息%s", string(body))
			if nil != err {
				util.Glog.Errorf("微信支付-解析返回信息-失败%v", err)
				err = &util.GError{Code: -1, Err: "解析微信返回信息失败"}
			}
		} else {
			util.Glog.Errorf("微信支付-读取请求返回Body-失败%v", err)
			err = &util.GError{Code: -1, Err: "读取请求返回Body失败"}
		}
	} else {
		util.Glog.Errorf("微信支付-HTTPPOST-失败%v", err)
		err = &util.GError{Code: -1, Err: "请求微信支付网络连接失败"}
	}
	// }
	//  else {
	// 	err = &util.GError{Code: -1, Err: "转换支付金额错误" + err1.Error()}
	// }
	return returnMsg, err
}

//获取微信支付发送数据XML
func getxml(attach string, out_trade_no, body string, total_fee int64, spbill_create_ip string, trade_type string) string {
	if "" == trade_type {
		trade_type = "NATIVE"
	}

	var params = UnifiedOrderParams{
		Appid:            config.Appid,
		Attach:           attach,
		Mch_id:           config.Mch_id,
		Nonce_str:        util.RandNumber(6),
		Body:             body,
		Out_trade_no:     out_trade_no,
		Total_fee:        total_fee,
		Spbill_create_ip: spbill_create_ip,
		Notify_url:       config.Notify_url,
		Trade_type:       trade_type,
	}
	stringSignTemp := "appid=" + params.Appid + "&attach=" + params.Attach + "&body=" + params.Body + "&mch_id=" + params.Mch_id + "&nonce_str=" + params.Nonce_str + "&notify_url=" + params.Notify_url + "&out_trade_no=" + params.Out_trade_no + "&spbill_create_ip=" + params.Spbill_create_ip + "&total_fee=" + strconv.FormatInt(params.Total_fee, 10) + "&trade_type=" + params.Trade_type + "&key=" + config.Apikey
	fmt.Println(stringSignTemp)
	sign := util.MD5(stringSignTemp)
	params.Sign = strings.ToUpper(sign)
	xml, _ := xml.Marshal(params)
	fmt.Println("xml", string(xml))
	return string(xml)
}
