package main

import (
	"encoding/json"
	"io"
	"net/http"
)

var bodyBuf []byte

func getBody(r *http.Request) []byte {
	n, err := r.Body.Read(bodyBuf)
	if err != nil && err != io.EOF {
		loger <- err.Error()
	}
	return bodyBuf[0:n]
}

func toJSON(v interface{}) []byte {
	j, err := json.Marshal(v)
	if err == nil {
		return j
	}
	loger <- err.Error()
	return []byte{}
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method == "POST" {
		options := make([]SelectOption, 0)
		json.Unmarshal(getBody(r), &options)
		setStyleMul(options)
		w.Write([]byte(`{"status":true}`))
		return
	}
	key := r.URL.Query().Get("key")
	w.Write(toJSON(
		map[string]interface{}{"status": true, "option": getSvgTypeOption(key), "default": getStyle(key)},
	))
}

func reloadSvgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	panels := make(map[string]Panel, 0)
	json.Unmarshal(getBody(r), &panels)
	ret := make(map[string]string, 0)
	for k, p := range panels {
		ret[k] = p.ToSvg()
	}
	w.Write(toJSON(ret))
}

func init() {
	bodyBuf = make([]byte, 2*1024*1024) //2M
}
