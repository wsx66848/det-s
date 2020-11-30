package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
)

type zipSaver struct {
	buffer *bytes.Buffer
	zip    *zip.Writer
}

func newZipSaver() zipSaver {
	buffer := new(bytes.Buffer)
	return zipSaver{buffer, zip.NewWriter(buffer)}
}

func (z *zipSaver) add(name, content string) {
	file, err := z.zip.Create(name)
	if err != nil {
		loger <- err.Error()
		return
	}
	file.Write([]byte(content))
}

func (z *zipSaver) toByte() []byte {
	z.zip.Close() // ignnore error
	return z.buffer.Bytes()
}

func downloadHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-zip-compressed; charset=UTF-8")
	panels := make(map[string]Panel, 0)
	json.Unmarshal(getBody(r), &panels)
	zip := newZipSaver()
	for k, p := range panels {
		zip.add(k+".svg", p.ToSvg())
	}
	w.Write(zip.toByte())
}
