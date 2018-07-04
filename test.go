package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type People interface {
	Speak(string) string
}

type Address struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	ParentCode string `json:"parentCode"`
	Level      int    `json:"level"`
}

var p = make([]Address, 0)

var ct = make([]Address, 0)
var cot = make([]Address, 0)
var all = make([]Address, 0)

func main() {

	p = getAddress("A000086000")
	for _, pp := range p {
		tmp := getAddress(pp.Code)
		ct = append(ct, tmp...)
		for _, tmpp := range tmp {
			tmppp := getAddress(tmpp.Code)
			cot = append(cot, tmppp...)
		}
	}
	all = append(all, p...)
	all = append(all, ct...)
	all = append(all, cot...)

	jsondata, _ := json.Marshal(all)
	f, _ := os.Create("./address.json")
	f.Write(jsondata)
	f.Close()
}
func getAddress(code string) []Address {
	var temp = make([]Address, 0)
	var url = "http://www.sf-express.com/sf-service-owf-web/service/region/" + code + "/subRegions?level=1&lang=sc&region=cn&translate="
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	err = json.Unmarshal(body, &temp)
	return temp
}
