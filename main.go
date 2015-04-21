package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ServeHTTP struct {
	paths []string
}

func (s *ServeHTTP) DealReq(rw http.ResponseWriter, r *http.Request) {
	var file *os.File
	var fileName string
	reqPath := r.URL.Path
	for _, basePath := range s.paths {
		fileName = basePath + reqPath
		fi, e := os.Stat(fileName)
		if e == nil && !fi.IsDir() {
			file, _ = os.Open(fileName)
			defer file.Close()
			break
		}
	}
	if file == nil {
		http.NotFound(rw, r)
	}
	os.Rename(fileName, fileName+"_"+time.Now().Format("2006-01-02--15-04-05"))
	io.Copy(rw, file)
	file.Close()
}

type Config struct {
	paths []string
	port  string
	url   string
	befor string
	after string
	out   string
	itype string
}

func getConfig(conf *Config) {
	fpath := "."
	fport := "60000"
	furl := "http://locahost:60000/myfile.txt"
	fafter := ""
	fout := "out.txt"
	ftype := "client"
	flag.StringVar(&ftype, "type", "client", "-type=client")
	flag.StringVar(&fpath, "path", ".", `-path ="/home/path1, /home/path2"`)
	flag.StringVar(&fport, "port", "60000", `-port="60000"`)
	flag.StringVar(&furl, "url", "", `-url="http://localhost:60000/myfile.txt"`)
	flag.StringVar(&fafter, "after", "", `-after="after_script.bat"`)
	flag.StringVar(&fout, "out", "", `-out="./out.txt"`)
	flag.Parse()
	conf.paths = strings.Split(fpath, ",")
	conf.url = furl
	conf.out = fout
	conf.port = fport
	conf.after = fafter
	conf.itype = ftype
}

func httpListen(conf *Config) {
	serve := new(ServeHTTP)
	serve.paths = conf.paths
	h := http.HandlerFunc(serve.DealReq)
	http.ListenAndServe(":"+conf.port, h)
}

func clientWait(conf *Config) {
	rqurl := conf.url
	saveName := conf.out
	for {
		rsp, _ := http.Get(rqurl)
		defer rsp.Body.Close()

		if rsp.StatusCode != 200 {
			break
		}
		body, _ := ioutil.ReadAll(rsp.Body)
		ioutil.WriteFile(saveName, body, 0644)
		if conf.after != "" {
			runScript(conf.after)
		}
		time.Sleep(time.Second * 3)
	}
}

func runScript(scriptName string) {
	c := exec.Command("call", scriptName)
	if err := c.Run(); err != nil {
		fmt.Println("error", err)
	}
}

func main() {
	conf := new(Config)
	getConfig(conf)
	if conf.itype == "server" {
		httpListen(conf)
	}
	if conf.itype == "client" {
		clientWait(conf)
	}
}
