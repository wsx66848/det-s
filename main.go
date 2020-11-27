package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Args is command line parameters
type Args struct {
	help      bool
	addr      string
	file      string
	publicDir string
	debug     bool
}

var args Args
var loger chan string

func parseArg() {
	args = Args{}
	flag.BoolVar(&args.help, "h", false, "help")
	flag.StringVar(&args.addr, "a", ":8003", "[ip][:port]")
	flag.StringVar(&args.file, "f", "fsc/detect.py", "python file")
	flag.StringVar(&args.publicDir, "public", "vue/dist/", "public dir")
	flag.BoolVar(&args.debug, "d", false, "debug")
	flag.Parse()
	loger <- fmt.Sprint(args)
}

func main() {
	http.HandleFunc("/api/upload", uploadHandler)       // POST 上传文件,返回识别后的位置信息,json格式
	http.HandleFunc("/api/option", styleHandler)        // GET/POST
	http.HandleFunc("/api/reloadsvg", reloadSvgHandler) // POST
	http.Handle("/", &staticHandler{})                  // GET 代理静态文件
	http.ListenAndServe(args.addr, nil)
}

func record() {
	if args.debug {
		for str := range loger {
			fmt.Println(str)
		}
	} else {
		for {
			<-loger
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		// w.Write(call("fsc/test.jpg"))
		w.WriteHeader(404)
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
	stdout, err := exec.Command("python", args.file, file).CombinedOutput()
	if err != nil {
		loger <- "执行脚本失败: file[" + file + "] " + err.Error()
		loger <- string(stdout)
		return []byte("error")
	}
	panel := &Panel{}
	panel.InitPanel()
	for _, line := range strings.Split(string(stdout), "\n") {
		panel.Add(line)
	}
	panel.Format()
	return panel.ToJSON()
	// ret, err := json.Marshal(panel)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return []byte(err.Error())
	// }
	// return ret
}

func init() {
	parseArg()
	go record()
	if args.help {
		flag.CommandLine.Usage()
		return
	}
}
