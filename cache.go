package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	svgroot   = "./component"
	svgconfig = ".config.json"
	unkowICON = `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAAAQCAYAAACm53kpAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAAlUlEQVR42u2XwQrAMAhD+/8/7dhgMIbWRLQXPfQiEu2rabclIqvzWgNgACSL3pKP7ABoAEDbaAWAQqB+49/Ym6MtS0fTRTWsnJ0meRhxAMiGo42idRDdkglgpgS1CAtOy/c0BsBJAKwvo9bR6hMXZ8ybDKQMDeQy3cXgZ9Aaq+zmI3WOAvg3WH16jE28GPkdMv8C3QFc6+oge+3V3R8AAAAASUVORK5CYII=`
)

func errorSvg(err string) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="16"><text fill="rgb(192, 32, 32)" x="2" y="14">%s</text></svg>`, 8*len(err), err)
}

var (
	svgCache     map[string]interface{}
	defaultStyle map[string]string
)

// SelectOption label value
type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// getSvgString("manufacturer")
// getSvgString("network", "2m4.svg")
func getSvgString(tp, size string) string {
	df, ok := defaultStyle[tp]
	if !ok { // 没有配置默认样式(例如厂商)
		return errorSvg(tp)
	}
	svg, ok := svgCache[tp]
	if !ok { // 目标对应的文件夹不存在
		return errorSvg(tp + ":" + tp)
	}
	if strSvg, ok := svg.(string); ok {
		return strSvg // 当前不会出现这种结构
	}
	styleMap, _ := svg.(map[string]interface{})
	styles, ok := styleMap[df]
	if !ok { // 目标默认样式文件不存在
		return errorSvg(tp + ":" + df)
	}
	if strSvg, ok := styles.(string); ok {
		return strSvg // 默认样式
	}
	svgMap, ok := styles.(map[string]interface{})
	svg, ok = svgMap[size]
	if !ok { // 这个尺寸对应的文件不存在
		return errorSvg(tp + ":" + size)
	}
	if _, ok = svg.(string); !ok {
		return errorSvg(tp + ":" + size)
	}
	return svg.(string)
}

// getSvgTypeOption("manufacturer") => ["aruba.svg", "cisco.svg", ...]
// getSvgTypeOption("network") => ["style1", "style2", ...]
func getSvgTypeOption(tp string) []SelectOption {
	t, ok1 := svgCache[tp]
	if !ok1 {
		return []SelectOption{}
	}
	ret := []SelectOption{}
	switch t.(type) {
	case map[string]interface{}:
		for style1, if1 := range t.(map[string]interface{}) {
			if svg, ok := if1.(string); ok {
				ret = append(ret, SelectOption{svg, style1})
			}
			if style2, ok := if1.(map[string]interface{}); ok {
				if icon, ok := style2["base64"]; ok {
					ret = append(ret, SelectOption{icon.(string), style1})
				} else {
					ret = append(ret, SelectOption{unkowICON, style1})
				}
			}
		}
	default:
	}
	return ret
}

func setStyle(key, value string) {
	if _, ok := defaultStyle[key]; ok {
		defaultStyle[key] = value
	}
	content, _ := json.Marshal(defaultStyle)
	err := ioutil.WriteFile(svgroot+"/"+svgconfig, content, 0600)
	if err != nil {
		loger <- err.Error()
	}
}

func setStyleMul(opts []SelectOption) {
	for index := range opts {
		if _, ok := defaultStyle[opts[index].Label]; ok {
			defaultStyle[opts[index].Label] = opts[index].Value
		}
	}
	content, _ := json.Marshal(defaultStyle)
	err := ioutil.WriteFile(svgroot+"/"+svgconfig, content, 0600)
	if err != nil {
		loger <- err.Error()
	}
}

func getStyle(key string) string {
	if style, ok := defaultStyle[key]; ok {
		return style
	}
	return ""
}

func travelPath(root string) map[string]interface{} {
	ret := make(map[string]interface{})
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		path := root + "/" + info.Name()
		stat, _ := os.Stat(path)
		if stat.IsDir() {
			ret[info.Name()] = travelPath(path)
		} else {
			content, _ := ioutil.ReadFile(path)
			ret[info.Name()] = string(content)
		}
	}
	return ret
}

func init() {
	svgCache = travelPath(svgroot)
	defaultStyle = make(map[string]string)
	if config, ok := svgCache[svgconfig]; ok {
		json.Unmarshal([]byte(config.(string)), &defaultStyle)
		delete(svgCache, svgconfig)
	}
	for key, value := range svgCache {
		x, ok := value.(map[string]interface{})
		if !ok {
			loger <- "invalid key: " + key
			continue
		}
		if style, ok := defaultStyle[key]; ok {
			if _, ok = x[style]; ok {
				continue
			}
		}
		for k := range x {
			defaultStyle[key] = k
			break
		}
	}
	// setStyle("", "")
	// fmt.Println(svgCache)
}
