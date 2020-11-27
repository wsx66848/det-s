package main

import (
	"io/ioutil"
	"net/http"
	"path"
)

type staticHandler struct {
}

// ContentType enum Content-Type
var ContentType = map[string]string{
	"ico":   "image/vnd.microsoft.icon",
	".css":  "text/css",
	".js":   "application/javascript",
	".woff": "font/woff2",
	".ttf":  "font/ttf",
}

func (h *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loger <- r.Method + ": " + r.RequestURI
	content, err := ioutil.ReadFile(args.publicDir + r.RequestURI)
	if err != nil {
		content, err = ioutil.ReadFile(args.publicDir + "/static" + r.RequestURI)
	}
	if err == nil {
		ext := path.Ext(r.RequestURI)
		if ct, ok := ContentType[ext]; ok {
			w.Header().Set("Content-Type", ct)
		}
		w.Write(content)
		return
	}
	content, err = ioutil.ReadFile(args.publicDir + "/index.html")
	if err == nil {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write(content)
		return
	}
	w.WriteHeader(404)
}
