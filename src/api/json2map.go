package api

import (
	"encoding/json"
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
		for _, vv := range v.([]interface{}) {
			js.Set(k, vv)
		}
	case map[string]interface{}:
		for kk, vv := range v.(map[string]interface{}) {
			js.Set(k+"/"+kk, vv)
		}
	default:
		if _, ok := js.data[k]; !ok {
			js.data[k] = v
		}
	}

	return
}
