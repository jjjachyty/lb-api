package wx

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lb-api/util"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

//统一下单URL
var UnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

//退款URL
var ReFundURL = "https://api.mch.weixin.qq.com/secapi/pay/refund"

//系统收取手续费
const systemCharge = 0

//微信全局参数配置
var config = struct {
	AppID           string
	MchID           string
	ApiKey          string
	PayNotifyURL    string
	RefundNotifyURL string
}{
	AppID:           "wxb5f20ab8f1933772",
	MchID:           "1501291661",
	ApiKey:          "jjjachyty929133jjjachyty92913364",
	PayNotifyURL:    "https://www.fortravel.cn/api/v1/pay/wx/notifiy",
	RefundNotifyURL: "https://www.fortravel.cn/api/v1/pay/wx/refund/notifiy",
}

//回调返回信息
type ReturnMsg struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

//支付请求参数
type NativePayParams struct {
	Appid  string `xml:"appid"` //公众账号ID
	Attach string `xml:"attach"`
	Body   string `xml:"body"`   //商品描述
	MchID  string `xml:"mch_id"` //商户号

	NonceStr  string `xml:"nonce_str"`  //随机字符串
	NotifyURL string `xml:"notify_url"` //通知地址
	// Openid         string `xml:"openid"`           //用户标识	trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识
	OutTradeNO string `xml:"out_trade_no"` //商户订单号
	// ProductID      string `xml:"product_id"`       //非必填,trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	SpbillCreateIP string `xml:"spbill_create_ip"` //终端IP
	TotalFee       int64  `xml:"total_fee"`        //标价金额
	TradeType      string `xml:"trade_type"`       //交易类型
	Sign           string `xml:"sign"`             //签名

	// deviceInfo string //设备号,非必填
	// signType   string //签名类型，默认为MD5，支持HMAC-SHA256和MD5。非必填
	// detail     string //商品详情 ,非必填
	// feeType    string //标价币种,非必填 默认CNy
	// timeStart  string //交易起始时间,非必填
	// timeExpire string //交易结束时间,非必填
	// goodsTag   string //订单优惠标记,非必填
	// limitPay   string //非必填,上传此参数no_credit--可限制用户不能使用信用卡支付

}

//通知回调数据
type PayNotify struct {
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

//退款申请请求参数
type RefundParams struct {
	Appid       string `xml:"appid"`
	MchID       string `xml:"mch_id"`
	NonceStr    string `xml:"nonce_str"`
	NotifyURL   string `xml:"notify_url"`
	OutRefundNo string `xml:"out_refund_no"`

	OutTradeNo string `xml:"out_trade_no"`
	RefundDesc string `xml:"refund_desc"`

	RefundFee int64 `xml:"refund_fee"`

	Sign     string `xml:"sign"`
	TotalFee int64  `xml:"total_fee"`
}

//退款申请请求返回信息
type RefundReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

type RefundNotify struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	ReqInfo    string `xml:"req_info"`
}

type RefundNotifyDecrypt struct {
	OutRefundNO         string  `xml:"out_refund_no"`
	OutTradeNO          string  `xml:"out_trade_no"`
	RefundAccount       string  `xml:"refund_account"`
	RefundFee           string  `xml:"refund_fee"`
	RefundID            string  `xml:"refund_id"`
	RefundRecvAccout    string  `xml:"refund_recv_accout"`
	RefundRequestSource string  `xml:"refund_request_source"`
	RefundStatus        string  `xml:"refund_status"`
	SettlementRefundFee float64 `xml:"settlement_refund_fee"`
	SettlementTotalFee  float64 `xml:"settlement_total_fee"`
	SuccessTime         string  `xml:"success_time"`
	TotalFee            float64 `xml:"total_fee"`
	TransactionID       string  `xml:"transaction_id"`
}

type NotifyReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// 根据实体获取签名
func getSign(nt interface{}) string {
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
			case int64:
				fieldValueStr = strconv.FormatInt(rtFieldValue.(int64), 10)
			}
			if "sign" != xmltage { //签名不用

				if i > 0 {
					checkStr += "&" + xmltage + "=" + fieldValueStr
				} else {
					checkStr += xmltage + "=" + fieldValueStr
				}
			}
		}
	}
	checkStr += "&key=" + config.ApiKey
	util.Glog.Debugf("签名串", checkStr)
	return checkStr
}

//校验签名
func checkSign(checkStr string, sigin string) bool {
	fmt.Println("MD5加密串", checkStr)
	if strings.ToUpper(util.MD5(checkStr)) == sigin {
		return true
	}
	return false
}

//获取微信支付发送数据XML
func getxml(params interface{}) string {

	rtValue := reflect.ValueOf(params)
	fmt.Println("金额", rtValue.Elem().FieldByName("TotalFee").Interface())
	signStr := getSign(params)
	sign := util.MD5(signStr)
	md5Sign := strings.ToUpper(sign)
	rtValue.Elem().FieldByName("Sign").SetString(md5Sign)

	xml, _ := xml.Marshal(params)
	fmt.Println("xml", string(xml))
	return string(xml)
}

func requestWx(posturl string, params interface{}, retStruc interface{}) (err error) {
	sendData := getxml(params)
	resp, err := http.Post(posturl, "application/xml", strings.NewReader(sendData))
	if nil == err {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if nil == err {
			err = xml.Unmarshal(body, &retStruc)
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
		util.Glog.Errorf("微信支付-HTTPPOST%s-失败%v", posturl, err)
		err = &util.GError{Code: -1, Err: "请求微信支付网络连接失败"}
	}
	return err
}

func requestWxRefund(posturl string, params interface{}, retStruc interface{}) (err error) {

	tlsConfig, err := getTLSConfig()
	if err != nil {
		return err
	}

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	sendData := getxml(params)
	resp, err := client.Post(posturl, "application/xml", strings.NewReader(sendData))
	if nil == err {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if nil == err {
			err = xml.Unmarshal(body, &retStruc)
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
		util.Glog.Errorf("微信支付-HTTPPOST%s-失败%v", posturl, err)
		err = &util.GError{Code: -1, Err: "请求微信支付网络连接失败"}
	}
	return err
}

func getTLSConfig() (*tls.Config, error) {
	var _tlsConfig *tls.Config
	if _tlsConfig != nil {
		return _tlsConfig, nil
	}

	// load cert
	cert, err := tls.LoadX509KeyPair("./controlers/pay/wx/apiclient_cert.pem", "./controlers/pay/wx/apiclient_key.pem")
	if err != nil {
		util.Glog.Errorln("load wechat keys fail", err)
		return nil, err
	}

	// load root ca
	// caData, err := ioutil.ReadFile(wechatCAPath)
	// if err != nil {
	// 	glog.Errorln("read wechat ca fail", err)
	// 	return nil, err
	// }
	// pool := x509.NewCertPool()
	// pool.AppendCertsFromPEM(caData)

	_tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		// RootCAs:      pool,
	}
	return _tlsConfig, nil
}
