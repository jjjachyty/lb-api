package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

const apiurl = "https://api-cn.faceplusplus.com/cardpp/v1/ocridcard"
const apiKey = "HwdG52YWaf8tgX-FNwyHY68jKcROhjx7"
const apiSecret = "o00516a7fjOAP3CoofuzSGy60jxPXJcf"

type results struct {
}

func GetIDCardText(imageBase64 string) (gjson.Result, error) {
	var err error
	var params = make(url.Values)
	var result gjson.Result
	params.Set("api_key", apiKey)
	params.Set("api_secret", apiSecret)

	params.Set("image_base64", imageBase64)

	resp, err := http.PostForm(apiurl, params)
	if err != nil {
		Glog.Errorf("身份识别-PostForm%v", err)
		return gjson.Result{}, err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	Glog.Debugf("身份证识别-%s", string(bt))
	cards := gjson.Get(string(bt), "cards")
	fmt.Println("len(cards.Array())", len(cards.Array()), cards)
	if 0 == len(cards.Array()) {
		err = &GError{Code: 3001, Err: "图片非身份证图片"}
	} else {
		result = cards.Array()[0]
	}

	return result, err
}
