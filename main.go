package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ServeHTTP struct {
	paths []string
}

func (s *ServeHTTP) DealReq(rw http.ResponseWriter, r *http.Request) {
	var file *os.File
	var fileName string
	fmt.Println(time.Unix(time.Now().Unix(), 0).String())
	fmt.Println(time.Now().Format("2006-01-02--15-04-05"))
	reqPath := r.URL.Path
	for _, basePath := range s.paths {
		fileName = basePath + reqPath
		fi, e := os.Stat(fileName)
		if e == nil && !fi.IsDir() {
			file, _ = os.Open(fileName)
			break
		}
	}
	if file == nil {
		http.NotFound(rw, r)
	}
	os.Rename(fileName, fileName+"_"+time.Now().Format("2006-01-02--15-04-05"))
	io.Copy(rw, file)
}

func main() {
	serve := new(ServeHTTP)
	serve.paths = []string{"."}
	h := http.HandlerFunc(serve.DealReq)
	http.ListenAndServe(":60000", h)
}
