package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const APPKEY = "3222320339d628c7d4aa099c9150bfc4" //您申请的APPKEY
const JuheURL = "http://v.juhe.cn/sms/send"

//2.发送注册短信
func SendSMS(mobile, tplId, tplValue string) error {
	var err error
	//请求地址
	switch tplId {
	case "80014": //注册

		//初始化参数
		param := url.Values{}

		//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
		param.Set("mobile", mobile)                                          //接收短信的手机号码
		param.Set("tpl_id", tplId)                                           //短信模板ID，请参考个人中心短信模板设置
		param.Set("tpl_value", url.QueryEscape("#code#="+tplValue+"&#m#=1")) //变量名和变量值对。如果你的变量名或者变量值中带有#&amp;=中的任意一个特殊符号，请先分别进行urlencode编码后再传递，&lt;a href=&quot;http://www.juhe.cn/news/index/id/50&quot; target=&quot;_blank&quot;&gt;详细说明&gt;&lt;/a&gt;
		param.Set("key", APPKEY)                                             //应用APPKEY(应用详细页查询)
		param.Set("dtype", "json")
		fmt.Println("#code#="+tplValue+"&#m#=1", url.QueryEscape("#code#="+tplValue+"&#m#=1"))
		//返回数据的格式,xml或json，默认json
		Glog.Debugf("发送短信-请求参数%v", param)

		//发送请求
		data, err := Post(JuheURL, param)
		if err == nil {

			var netReturn map[string]interface{}
			json.Unmarshal(data, &netReturn)
			Glog.Debugf("发送短信-返回信息%s", string(data))
			if netReturn["error_code"].(float64) != 0 {
				err = &GError{Code: netReturn["error_code"].(float64), Err: "短信接口错误，请联系管理员"}
				return err
			}
		} else {
			Glog.Errorf("发送短信-请求失败%v", err)

		}

	}
	fmt.Println("err:---2--", err)
	return err

}

// get 网络请求
func Get(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		Glog.Errorf("发送短信-解析url错误%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		Glog.Errorf("发送短信-错误%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values) (rs []byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		Glog.Errorf("发送短信-PostForm%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
