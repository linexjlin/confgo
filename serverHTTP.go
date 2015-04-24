package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ServeHTTP struct {
	conf *Config
}

//Ceate http Listen
func (s *ServeHTTP) httpListen() {
	h := http.HandlerFunc(s.DealReq)
	http.ListenAndServe(":"+s.conf.port, h)
}

//Coming request handle
func (s *ServeHTTP) DealReq(rw http.ResponseWriter, r *http.Request) {
	var file *os.File
	var fileName string
	reqPath := r.URL.Path
	fmt.Println(reqPath)
	fmt.Println("the path length is:", len(s.conf.paths))
	for _, basePath := range s.conf.paths {
		fileName = basePath + reqPath
		fmt.Println("try to find file", fileName)
		fi, e := os.Stat(fileName)
		if e == nil && !fi.IsDir() {
			file, _ = os.Open(fileName)
			defer file.Close()
			break
		}
	}

	if file == nil {
		http.NotFound(rw, r)
		return
	}

	io.Copy(rw, file)
	file.Close()

	//do not rename befor script or after script
	if r.URL.Path == "/"+s.conf.befor || r.URL.Path == "/"+s.conf.after {
		return
	}

	newName := fileName + "_" + time.Now().Format("2006-01-02--15-04-05")

	fmt.Println("rename", fileName, "to", newName)
	os.Rename(fileName, newName)
}
