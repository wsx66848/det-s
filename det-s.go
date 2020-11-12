package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var args map[string]interface{}
var loger chan string

// Ulenth 1U的长度
// Uwidth 1U的宽度
// Uheight 1U的高度
const (
	Ulenth  float32 = 370
	Uwidth  float32 = 435
	Uheight float32 = 44.45
)

// Panel 面板上所有目标
type Panel struct {
	x1, x2, y1, y2 float32
	U              int      `json:"u"`
	Network        []Target `json:"network"`
	Optical        []Target `json:"optical"`
	USB            []Target `json:"usb"`
	Manufacturer   []Target `json:"manufacturer"`
	Indicatorlight []Target `json:"indicatorlight"`
	Backplane      Target   `json:"backplane"`
}

// Target 面板上目标位置信息
type Target struct {
	Type        string  `json:"type"`
	Px          float32 `json:"px"`
	Py          float32 `json:"py"`
	Width       float32 `json:"width"`
	Height      float32 `json:"height"`
	Probability float32 `json:"probability"`
}

func parseArg() {
	args = make(map[string]interface{}, 10)
	help := flag.Bool("h", false, "help")
	addr := flag.String("a", ":8003", "[ip][:port]")
	file := flag.String("f", "fsc/detect.py", "python file")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()
	args["help"] = *help
	args["addr"] = *addr
	args["file"] = *file
	args["debug"] = *debug
	loger <- fmt.Sprint(args)
}

func main() {
	loger = make(chan string, 10)
	parseArg()
	if args["help"].(bool) {
		flag.CommandLine.Usage()
		return
	}
	go record()
	http.HandleFunc("/", index)
	http.ListenAndServe(args["addr"].(string), nil)
}

func record() {
	if args["debug"].(bool) {
		for str := range loger {
			fmt.Println(str)
		}
	} else {
		for {
			<-loger
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(`{"u":1,"network":[{"type":"netport","px":15.429741,"py":23.979607,"width":14.836289,"height":11.697369,"probability":1},{"type":"netport","px":93.76534,"py":9.942763,"width":14.242838,"height":11.697369,"probability":0.88},{"type":"netport","px":108.00819,"py":9.942763,"width":14.242838,"height":11.697369,"probability":0.88},{"type":"netport","px":93.76534,"py":24.564474,"width":14.242838,"height":11.1125,"probability":0.91},{"type":"netport","px":108.00819,"py":24.564474,"width":14.242838,"height":11.1125,"probability":0.91},{"type":"netport","px":186.93724,"py":24.564474,"width":13.946112,"height":11.1125,"probability":0.83},{"type":"netport","px":200.88335,"py":24.564474,"width":13.946112,"height":11.1125,"probability":0.83},{"type":"netport","px":188.12415,"py":10.527632,"width":13.946112,"height":11.1125,"probability":0.88},{"type":"netport","px":202.07025,"py":10.527632,"width":13.946112,"height":11.1125,"probability":0.88},{"type":"netport","px":36.200546,"py":9.942763,"width":14.539563,"height":11.697369,"probability":0.94},{"type":"netport","px":50.74011,"py":9.942763,"width":14.539563,"height":11.697369,"probability":0.94},{"type":"netport","px":65.27967,"py":9.942763,"width":14.539563,"height":11.697369,"probability":0.94},{"type":"netport","px":79.81924,"py":9.942763,"width":14.539563,"height":11.697369,"probability":0.94},{"type":"netport","px":37.387447,"py":23.979607,"width":14.242838,"height":11.1125,"probability":0.87},{"type":"netport","px":51.630287,"py":23.979607,"width":14.242838,"height":11.1125,"probability":0.87},{"type":"netport","px":65.87312,"py":23.979607,"width":14.242838,"height":11.1125,"probability":0.87},{"type":"netport","px":80.11596,"py":23.979607,"width":14.242838,"height":11.1125,"probability":0.87},{"type":"netport","px":131.1528,"py":10.527632,"width":14.242838,"height":11.1125,"probability":0.86},{"type":"netport","px":145.39563,"py":10.527632,"width":14.242838,"height":11.1125,"probability":0.86},{"type":"netport","px":159.63847,"py":10.527632,"width":14.242838,"height":11.1125,"probability":0.86},{"type":"netport","px":173.8813,"py":10.527632,"width":14.242838,"height":11.1125,"probability":0.86},{"type":"netport","px":131.1528,"py":23.979607,"width":14.687926,"height":11.1125,"probability":0.9},{"type":"netport","px":145.84071,"py":23.979607,"width":14.687926,"height":11.1125,"probability":0.9},{"type":"netport","px":160.52864,"py":23.979607,"width":14.687926,"height":11.1125,"probability":0.9},{"type":"netport","px":175.21657,"py":23.979607,"width":14.687926,"height":11.1125,"probability":0.9}],"optical":[{"type":"optical_netport","px":280.70258,"py":26.31908,"width":16.023191,"height":11.1125,"probability":0.99},{"type":"optical_netport","px":242.12823,"py":25.734211,"width":16.023191,"height":11.697369,"probability":0.99},{"type":"optical_netport","px":221.35744,"py":25.734211,"width":16.023191,"height":11.697369,"probability":0.99},{"type":"optical_netport","px":261.71213,"py":25.734211,"width":16.616644,"height":11.697369,"probability":0.98}],"usb":null,"manufacturer":[{"type":"manufacturer","px":0,"py":2.3394737,"width":22.551159,"height":8.773026,"probability":1}],"indicatorlight":[{"type":"indicatorlight","px":5.341064,"py":28.073685,"width":5.9345155,"height":6.4335527,"probability":0.45}],"backplane":{"type":"backplane","px":0,"py":0,"width":435,"height":44.45,"probability":1}}`))
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		loger <- "上传图片失败: " + err.Error()
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	fileName := fmt.Sprintf(os.TempDir()+"/%d.jpg", time.Now().UnixNano())
	dist, err := os.Create(fileName)
	if err != nil {
		loger <- "创建临时文件失败: " + err.Error()
		w.Write([]byte(err.Error()))
		return
	}
	_, err = io.Copy(dist, file)
	dist.Close()
	if err != nil {
		loger <- "临时文件写入失败: " + err.Error()
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(call(fileName))
	err = os.Remove(fileName)
	if err != nil {
		loger <- "删除临时文件失败: " + err.Error()
		return
	}
}

func call(file string) []byte {
	stdout, err := exec.Command("python", args["file"].(string), file).CombinedOutput()
	if err != nil {
		loger <- "执行脚本失败: file[" + file + "] " + err.Error()
		loger <- string(stdout)
		return []byte("{u:0}")
	}
	panel := &Panel{}
	panel.InitPanel()
	for _, line := range strings.Split(string(stdout), "\n") {
		panel.Add(line)
	}
	panel.Format()
	ret, err := json.Marshal(panel)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}
	return ret
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
	for index := range p.Optical {
		p.Optical[index].vector(mx, my, p.Backplane)
	}
	for index := range p.USB {
		p.USB[index].vector(mx, my, p.Backplane)
	}
	for index := range p.Indicatorlight {
		p.Indicatorlight[index].vector(mx, my, p.Backplane)
	}
	for index := range p.Manufacturer {
		p.Manufacturer[index].vector(mx, my, p.Backplane)
	}
	p.Backplane.vector(mx, my, p.Backplane)
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
