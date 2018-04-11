package api

import (
	"fmt"
	"mlog"
	"net/http"
)

var (
	key    = "7d18f2bef987ddcb9fedc0987ea406e3"
	webkey = "bb3109b05244845dfdc49743152e5838"
)

func RegisterHandle() {
	http.HandleFunc("/geo", geoHandle)
	http.HandleFunc("/regeo", regeoHandle)
	return
}

func geoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	city := req.Form["city"][0]
	for _, address := range req.Form["addres"] {
		url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/geo?city=%s&address=%s&output=json&key=%s", city, address, key)
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
		mlog.Debug(js)

		w.Write([]byte(js.Get("geocodes/formatted_address").(string) + "\n"))
	}

	return
}

func regeoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	location := req.Form["location"][0]
	url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/regeo?output=json&location=%s&key=%s", location, key)
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
	w.Write([]byte(js.Get("regeocodes/formatted_address").(string)))
}
