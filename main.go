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

var args map[string]interface{}
var loger chan string

func parseArg() {
	args = make(map[string]interface{}, 10)
	help := flag.Bool("h", false, "help")
	addr := flag.String("a", ":8003", "[ip][:port]")
	file := flag.String("f", "fsc/detect.py", "python file")
	publicDir := flag.String("public", "vue/dist/", "public dir")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()
	args["help"] = *help
	args["addr"] = *addr
	args["file"] = *file
	args["public"] = *publicDir
	args["debug"] = *debug
	loger <- fmt.Sprint(args)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", &staticHandler{})
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
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
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
		return []byte("error")
	}
	panel := &Panel{}
	panel.InitPanel()
	for _, line := range strings.Split(string(stdout), "\n") {
		panel.Add(line)
	}
	panel.Format()
	return []byte(panel.ToSvg())
	// ret, err := json.Marshal(panel)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return []byte(err.Error())
	// }
	// return ret
}

func init() {
	loger = make(chan string, 10)
	parseArg()
	if args["help"].(bool) {
		flag.CommandLine.Usage()
		return
	}
	go record()
}
