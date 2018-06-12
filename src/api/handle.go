package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mlog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	gdkey  = ""
	webkey = ""
	appid  = ""
	secret = ""
)

var (
	server *gin.Engine
)

var (
	sessionMap = make(map[string][]byte)
)

func Run(addr, port, mode string) {
	switch mode {
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	server = gin.Default()
	registerHandle()
	listenAddr := fmt.Sprintf("%s:%s", addr, port)
	server.Run(listenAddr)
}

func registerHandle() {
	// http.HandleFunc("/geo", geoHandle)
	// http.HandleFunc("/regeo", regeoHandle)
	// http.HandleFunc("/wx_login", wxLoginHandle)
	// http.HandleFunc("/wx_decode", wxDecodeHandle)
	// http.HandleFunc("/wx_run", wxRunHandle)
	server.GET("/geo", geoHandle)
	server.GET("/regeo", regeoHandle)

	return
}

func geoHandle(c *gin.Context) {
	city := c.Query("city")
	for _, address := range c.QueryArray("addres") {
		url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/geo?city=%s&address=%s&output=json&key=%s", city, address, gdkey)
		body, err := httpHandler.Get(url)
		if err != nil {
			mlog.Error(err)
			continue
		}
		js, err := NewJsonMap(body)
		if err != nil {
			mlog.Error(err)
			continue
		}

		// c.Writer.Write( + "\n"))
		c.String(http.StatusOK, "%s\n", js.Get("geocodes/formatted_address"))

	}

	return
}

func regeoHandle(c *gin.Context) {
	location := c.Param("location")
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

	c.String(http.StatusOK, "%s", []byte(js.Get("regeocode/formatted_address").(string)))
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
	l := len(wxRunDataRes.StepInfoList)
	seventotal := 0
	for i := 1; i < 8; i++ {
		seventotal += wxRunDataRes.StepInfoList[l-i].Step
	}
	wxRunDataRes.SevenDayAvg = seventotal / 7
	wxRunDataRes.DiffLastDay = wxRunDataRes.StepInfoList[l-1].Step - wxRunDataRes.StepInfoList[l-2].Step
	res, err := json.Marshal(wxRunDataRes)
	if err != nil {
		mlog.Error(err)
		return
	}
	w.Write(res)

	return
}
