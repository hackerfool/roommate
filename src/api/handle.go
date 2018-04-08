package api

import (
	"encoding/json"
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

		js := &geoRespone{}
		err = json.Unmarshal(body, js)
		if err != nil {
			mlog.Error(err)
			continue
		}
		mlog.Debug(js)

		w.Write(body)
	}

	return
}
