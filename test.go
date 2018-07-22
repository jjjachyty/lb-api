package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lb-api/util"
	"net/http"
	"strconv"
	"strings"
)

var UnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

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
}

type UnifiedOrderParams struct {
	Appid            string `xml:"appid"`  //公众账号ID
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

//获取骑牛云上传key
func main() {
	resp, err := http.Post(UnifiedOrderURL, "application/xml", strings.NewReader(getxml("9527", "支付测试", 1, "127.0.0.1", "")))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	returnMsg := new(ReturnMsg)
	err = xml.Unmarshal(body, &returnMsg)
	fmt.Println("%v", returnMsg)
}

//
func getxml(out_trade_no, body string, total_fee int64, spbill_create_ip string, trade_type string) string {
	if "" == trade_type {
		trade_type = "NATIVE"
	}

	var params = UnifiedOrderParams{
		Appid:            config.Appid,
		Mch_id:           config.Mch_id,
		Nonce_str:        util.RandNumber(6),
		Body:             body,
		Out_trade_no:     out_trade_no,
		Total_fee:        total_fee,
		Spbill_create_ip: spbill_create_ip,
		Notify_url:       config.Notify_url,
		Trade_type:       trade_type,
	}
	stringSignTemp := "appid=" + params.Appid + "&body=" + params.Body + "&mch_id=" + params.Mch_id + "&nonce_str=" + params.Nonce_str + "&notify_url=" + params.Notify_url + "&out_trade_no=" + params.Out_trade_no + "&spbill_create_ip=" + params.Spbill_create_ip + "&total_fee=" + strconv.FormatInt(params.Total_fee, 10) + "&trade_type=" + params.Trade_type + "&key=" + config.Apikey
	fmt.Println(stringSignTemp)
	sign := util.MD5(stringSignTemp)
	params.Sign = strings.ToUpper(sign)
	xml, _ := xml.Marshal(params)
	fmt.Println("xml", string(xml))
	return string(xml)
}
