package api

import (
	"encoding/json"
	"mlog"
)

type jsonmap struct {
	data map[string]interface{}
}

func NewJsonMap(data []byte) (jsonmap, error) {
	js := jsonmap{}
	return js, js.Unmarshal(data)
}

func (js *jsonmap) Unmarshal(data []byte) error {
	body := make(map[string]interface{})
	js.data = make(map[string]interface{})
	err := json.Unmarshal(data, &body)
	if err != nil {
		return err
	}

	mlog.Debug(body)
	for k, v := range body {
		js.Set(k, v)
	}

	return nil
}

func (js *jsonmap) Get(path string) interface{} {
	return js.data[path]
}

func (js *jsonmap) Set(k string, v interface{}) {
	switch v.(type) {
	case []interface{}:
		set(js.data, k, v)
		for _, vv := range v.([]interface{}) {
			js.Set(k, vv)
		}
	case map[string]interface{}:
		set(js.data, k, v)
		for kk, vv := range v.(map[string]interface{}) {
			js.Set(k+"/"+kk, vv)
		}
	default:
		set(js.data, k, v)
	}

	return
}

func set(m map[string]interface{}, k string, v interface{}) {
	if _, ok := m[k]; !ok {
		m[k] = v
	}
}
