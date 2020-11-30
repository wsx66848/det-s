package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Ulenth 1U的长度
// Uwidth 1U的宽度
// Uheight 1U的高度
const (
	Ulenth  float32 = 1063
	Uwidth  float32 = 1250 // 435
	Uheight float32 = 127  // 44.45
)

// Nwidth 网口的宽度
// Nheight 网口高度
const (
	Nwidth  float32 = 41
	Nheight float32 = 32
)

// Lwidth 网口的宽度
// Lheight 网口高度
const (
	Lwidth  float32 = 46
	Lheight float32 = 32
)

// NPrefix 网口ID前缀
// OPrefix 光口ID前缀
const (
	NPrefix = "g-"
	OPrefix = "e-"
)

// Panel 面板上所有目标
type Panel struct {
	x1, x2, y1, y2 float32
	U              int         `json:"u"`
	Network        TargetSlice `json:"-"`
	Optical        TargetSlice `json:"-"`
	NetworkGroup   []Group     `json:"network"`
	OpticalGroup   []Group     `json:"optical"`
	USB            TargetSlice `json:"usb"`
	Manufacturer   TargetSlice `json:"manufacturer"`
	Indicatorlight TargetSlice `json:"indicatorlight"`
	Backplane      Target      `json:"backplane"`
}

// Target 面板上目标位置信息
type Target struct {
	Type        string  `json:"type"`
	Px          float32 `json:"px"`
	Py          float32 `json:"py"`
	Width       float32 `json:"width"`
	Height      float32 `json:"height"`
	Probability float32 `json:"-"`
}

// InitPanel 初始化面板边缘
func (p *Panel) InitPanel() {
	p.x1, p.y1, p.x2, p.y2 = 1e+4, 1e+4, 0, 0
}

// 如果里面的目标超出背板边缘说明背板检查不准,需要更新边缘
func (p *Panel) updateRectangle(edge [4]float32) {
	if p.x1 > edge[0] {
		p.x1 = edge[0]
	}
	if p.x2 < edge[2] {
		p.x2 = edge[2]
	}
	if p.y1 > edge[1] {
		p.y1 = edge[1]
	}
	if p.y2 < edge[3] {
		p.y2 = edge[3]
	}
}

// Add parse line and add target to panel
func (p *Panel) Add(line string) {
	if !strings.HasPrefix(line, "bboxes") {
		return
	}
	boxlabel := strings.Split(line, ",")
	if len(boxlabel) < 2 {
		return
	}
	loger <- line
	// bboxes
	tmp := strings.Split(boxlabel[0], ":")[1]
	tmp = strings.TrimLeft(tmp, " [")
	tmp = strings.TrimRight(tmp, "]")
	var boxes [4]float32
	_, err := fmt.Sscanf(tmp, "%f %f %f %f", &boxes[0], &boxes[1], &boxes[2], &boxes[3])
	if err != nil {
		loger <- "字符串解析失败: " + err.Error()
		return
	}
	// label_text
	tmp = strings.Split(boxlabel[1], ":")[1]
	tmp = strings.TrimLeft(tmp, " ")
	tmp = strings.TrimRight(tmp, "\r")
	label := strings.Split(tmp, "|")
	probability, err := strconv.ParseFloat(label[1], 32)
	if err != nil {
		loger <- "数字解析失败: " + err.Error()
		return
	}
	target := Target{
		label[0],
		float32(boxes[0]),
		float32(boxes[1]),
		float32(boxes[2] - boxes[0]),
		float32(boxes[3] - boxes[1]),
		float32(probability),
	}
	p.updateRectangle(boxes)
	// fmt.Printf("target := %v\n", target)
	switch target.Type {
	case "netport":
		// p.Console = append(p.Console, target) // 这个应该是console口
		p.Network = append(p.Network, target)
		break
	case "two_netport":
		p.Network = append(p.Network, target.vsplite(2, "netport")...)
		break
	case "four_netport":
		p.Network = append(p.Network, target.vsplite(4, "netport")...)
		break
	case "optical_netport":
		p.Optical = append(p.Optical, target)
		break
	case "two_optical_netport":
		p.Optical = append(p.Optical, target.vsplite(2, "optical")...)
		break
	case "four_optical_netport":
		p.Optical = append(p.Optical, target.vsplite(4, "optical")...)
		break
	case "backplane":
		p.Backplane = target
		break
	case "manufacturer":
		p.Manufacturer = append(p.Manufacturer, target)
		break
	case "indicatorlight":
		p.Indicatorlight = append(p.Indicatorlight, target)
		break
	case "usb":
		p.USB = append(p.USB, target)
		break
	default:
		break
	}
}

// Format 按照实际尺寸(1U/2U/4U)格式化数值
func (p *Panel) Format() {
	p.Backplane.Px, p.Backplane.Py = p.x1, p.y1
	p.Backplane.Width, p.Backplane.Height = p.x2-p.x1, p.y2-p.y1
	p.U = int(.5 + p.Backplane.Height*Uwidth/p.Backplane.Width/Uheight)
	mx, my := Uwidth/p.Backplane.Width, float32(p.U)*Uheight/p.Backplane.Height
	loger <- fmt.Sprintf("base:[%dU %f %f], proportion:[%f %f]\n",
		p.U, p.Backplane.Px, p.Backplane.Py, mx, my)
	for index := range p.Network {
		p.Network[index].vector(mx, my, p.Backplane)
	}
	p.NetworkGroup = groupTargets(p.Network)
	for index := range p.Optical {
		p.Optical[index].vector(mx, my, p.Backplane)
	}
	p.OpticalGroup = groupTargets(p.Optical)
	for index := range p.USB {
		p.USB[index].vector(mx, my, p.Backplane)
		p.USB[index].round()
	}
	for index := range p.Indicatorlight {
		p.Indicatorlight[index].vector(mx, my, p.Backplane)
		p.Indicatorlight[index].round()
	}
	for index := range p.Manufacturer {
		p.Manufacturer[index].vector(mx, my, p.Backplane)
		p.Manufacturer[index].round()
	}
	p.Backplane.vector(mx, my, p.Backplane)
	p.Backplane.round()
}

// ToSvg return panel's svg file
func (p *Panel) ToSvg() string {
	var str string
	str += fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.f" height="%.f">`, p.Backplane.Width, p.Backplane.Height) + "\n"
	// 背板
	str += "<g>"
	str += getSvgString("backplane", fmt.Sprintf("%d.svg", p.U))
	str += "</g>\n"
	// 网口
	start := 1
	serialType, _ := strconv.Atoi(getStyle("serial"))
	for _, g := range p.NetworkGroup {
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, g.MinX, g.MinY)
		tmp := getSvgString("network", fmt.Sprintf("%dm%d.svg", g.Row, g.Col)) + "</g>\n"
		if g.Row == 1 {
			str += tmp
			continue
		}
		for i := 0; i < g.Col*g.Row; i += 2 {
			tmp = svgIDReplace(tmp, "${ID}", NPrefix+fmt.Sprint(start+i+serialType))
			tmp = svgIDReplace(tmp, "${ID}", NPrefix+fmt.Sprint(start+i+1-serialType))
		}
		str += tmp
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, g.MinX, g.MinY-17)
		str += serial(g.Col, start, Nwidth)
		str += "</g>\n"
		start += g.Col * g.Row
	}
	// 光口
	start = 1
	for _, g := range p.OpticalGroup {
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, g.MinX, g.MinY)
		tmp := getSvgString("optical", fmt.Sprintf("%dm%d.svg", g.Row, g.Col)) + "</g>\n"
		if g.Row == 1 {
			str += tmp
			continue
		}
		for i := 0; i < g.Col*g.Row; i += 2 {
			tmp = svgIDReplace(tmp, "${ID}", OPrefix+fmt.Sprint(start+i+serialType))
			tmp = svgIDReplace(tmp, "${ID}", OPrefix+fmt.Sprint(start+i+1-serialType))
		}
		str += tmp
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, g.MinX, g.MinY-17)
		str += serial(g.Col, start, Nwidth)
		str += "</g>\n"
		start += g.Col * g.Row
	}
	// 指示灯
	for _, target := range p.Indicatorlight {
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, target.Px, target.Py)
		str += getSvgString("usb", fmt.Sprintf("1m1.svg"))
		str += "</g>\n"
	}
	// USB
	for _, target := range p.USB {
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, target.Px, target.Py)
		str += getSvgString("usb", fmt.Sprintf("1m1.svg"))
		str += "</g>\n"
	}
	// 厂商
	for _, target := range p.Manufacturer {
		str += fmt.Sprintf(`<g transform="translate(%.1f,%.1f)">`, target.Px, target.Py)
		str += getSvgString("manufacturer", fmt.Sprintf("sugon.svg"))
		str += "</g>\n"
	}
	str += `</svg>`
	return str
}

// ToJSON Panel to json byte array
func (p *Panel) ToJSON() []byte {
	ret, _ := json.Marshal(*p)
	return ret
}

// FromJSON init Panel from json byte array
func (p *Panel) FromJSON(b []byte) {
	json.Unmarshal(b, p)
}

func (t Target) vsplite(n int, tp string) []Target {
	width := t.Width / float32(n)
	ret := []Target{}
	for i := 0; i < n; i++ {
		target := Target{
			tp,
			t.Px + float32(i)*width,
			t.Py,
			width,
			t.Height,
			t.Probability,
		}
		ret = append(ret, target)
	}
	return ret
}

func (t *Target) vector(mx, my float32, base Target) {
	t.Px = mx * (t.Px - base.Px)
	t.Py = my * (t.Py - base.Py)
	t.Width *= mx
	t.Height *= my
}

func (t *Target) round() {
	t.Px = round(t.Px, 1)
	t.Py = round(t.Py, 1)
	t.Width = round(t.Width, 1)
	t.Height = round(t.Height, 1)
}
