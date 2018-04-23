package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mlog"
	"net/http"
	"time"
)

var ()

var (
	sessionMap = make(map[string][]byte)
)

func RegisterHandle() {
	http.HandleFunc("/geo", geoHandle)
	http.HandleFunc("/regeo", regeoHandle)
	http.HandleFunc("/wx_login", wxLoginHandle)
	http.HandleFunc("/wx_decode", wxDecodeHandle)
	http.HandleFunc("/wx_run", wxRunHandle)
	return
}

func geoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	city := req.Form["city"][0]
	for _, address := range req.Form["addres"] {
		url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/geo?city=%s&address=%s&output=json&key=%s", city, address, gdkey)
		body, err := httpHandler.Get(url)
		if err != nil {
			mlog.Error(err)
			continue
		}
		mlog.Info(body)
		js, err := NewJsonMap(body)
		if err != nil {
			mlog.Error(err)
			continue
		}
		mlog.Debug(js)

		w.Write([]byte(js.Get("geocodes/formatted_address").(string) + "\n"))
	}

	return
}

func regeoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	location := req.Form["location"][0]
	url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/regeo?output=json&location=%s&key=%s", location, gdkey)
	body, err := httpHandler.Get(url)
	if err != nil {
		mlog.Error(err)
		return
	}

	js, err := NewJsonMap(body)
	if err != nil {
		mlog.Error(err)
		return
	}

	mlog.Debug(js)
	w.Write([]byte(js.Get("regeocode/formatted_address").(string)))
}

func wxLoginHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	mlog.Info(req.URL)
	code := req.Form["code"][0]
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
	body, err := httpHandler.Get(url)
	if err != nil {
		mlog.Error(err)
		return
	}

	js, err := NewJsonMap(body)
	if err != nil {
		mlog.Error(err)
		return
	}
	id := js.Get("openid").(string)
	keystr := js.Get("session_key").(string)
	key, _ := base64.StdEncoding.DecodeString(keystr)
	sessionMap[id] = key

	w.Write([]byte(id))
}

func wxDecodeHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	mlog.Info(req.URL)
	encryptedString := req.Form["encrypt"][0]
	ivString := req.Form["iv"][0]
	id := req.Form["openId"][0]

	encryptedData, _ := base64.StdEncoding.DecodeString(encryptedString)
	iv, _ := base64.StdEncoding.DecodeString(ivString)

	decodeData, err := aes128cbc(sessionMap[id], encryptedData, iv)
	if err != nil {
		mlog.Error(err)
		return
	}

	w.Write([]byte(decodeData))
}

func wxRunHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	mlog.Info(req.URL)
	encryptedString := req.Form["encrypt"][0]
	ivString := req.Form["iv"][0]
	id := req.Form["openId"][0]

	encryptedData, _ := base64.StdEncoding.DecodeString(encryptedString)
	iv, _ := base64.StdEncoding.DecodeString(ivString)

	decodeData, err := aes128cbc(sessionMap[id], encryptedData, iv)
	if err != nil {
		mlog.Error(err)
		return
	}

	wxRunDataReq := &wxRunDataRequest{}
	err = json.Unmarshal(decodeData, wxRunDataReq)
	if err != nil {
		mlog.Error(err)
		return
	}

	wxRunDataRes := &wxRunDataResponse{}
	for _, info := range wxRunDataReq.StepInfoList {
		date := time.Unix(info.Timestamp, 0)
		timestamp := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
		wxRunDataRes.StepInfoList = append(wxRunDataRes.StepInfoList, stepInfoList2db{Timestamp: timestamp, Step: info.Step})
	}
	res, err := json.Marshal(wxRunDataRes)
	if err != nil {
		mlog.Error(err)
		return
	}
	w.Write(res)

	return
}
