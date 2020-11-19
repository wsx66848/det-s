package main

import (
	"fmt"
	"io/ioutil"
)

const (
	errorSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"><text fill="rgb(192, 32, 32)" x="4" y="12">?</text></svg>`
)

var svgCache map[string]map[string]string

func getSvgString(tp, name string) string {
	if t, ok1 := svgCache[tp]; ok1 {
		if s, ok2 := t[name]; ok2 {
			return s
		}
	}
	loger <- fmt.Sprintf("获取svg失败: [%s] [%s]\n", tp, name)
	return errorSvg
}

func getSvgType(tp string) map[string]string {
	if t, ok1 := svgCache[tp]; ok1 {
		return t
	}
	return make(map[string]string)
}

func init() {
	root := "./component"
	svgCache = make(map[string]map[string]string)
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		panic("读取svg文件失败")
	}
	for _, info := range infos {
		dir := root + "/" + info.Name()
		m := make(map[string]string)
		svgCache[info.Name()] = m
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			path := dir + "/" + file.Name()
			content, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}
			m[file.Name()] = string(content)
		}
	}
	// fmt.Println(svgCache)
}
