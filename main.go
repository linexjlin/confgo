package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getConfig(conf *Config) error {
	fpath := "."
	fport := "60000"
	furl := "http://locahost:60000/myfile.txt"
	fbefor := ""
	fafter := ""
	fout := "out.txt"
	ftype := "client"
	flag.StringVar(&ftype, "type", "client", "-type=client")
	flag.StringVar(&fpath, "path", ".", `-path ="/home/path1, /home/path2"`)
	flag.StringVar(&fport, "port", "60000", `-port="60000"`)
	flag.StringVar(&furl, "url", "http://localhost:60000/aa.txt", `-url="http://localhost:60000/myfile.txt"`)
	flag.StringVar(&fbefor, "befor", "", `-befor="befor.bat"`)
	flag.StringVar(&fafter, "after", "", `-after="after.bat"`)
	flag.StringVar(&fout, "out", "out.txt", `-out="./out.txt"`)
	flag.Parse()
	conf.paths = strings.Split(fpath, ",")
	conf.url = furl
	if conf.url == "" {
		return errors.New("Url no given")
	}
	conf.out = fout
	conf.port = fport
	conf.after = fafter
	conf.itype = ftype
	return nil
}

type ServeHTTP struct {
	paths []string
}

func (s *ServeHTTP) DealReq(rw http.ResponseWriter, r *http.Request) {
	var file *os.File
	var fileName string
	reqPath := r.URL.Path
	for _, basePath := range s.paths {
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
	}
	io.Copy(rw, file)
	file.Close()
	fmt.Println("rename ", fileName, " to ", fileName+"_"+time.Now().Format("2006-01-02--15-04-05"))
	e := os.Rename(fileName, fileName+"_"+time.Now().Format("2006-01-02--15-04-05"))
	fmt.Println(e)
}

func (s *ServeHTTP) httpListen(conf *Config) {
	serve := new(ServeHTTP)
	serve.paths = conf.paths
	h := http.HandlerFunc(serve.DealReq)
	http.ListenAndServe(":"+conf.port, h)
}

type Client struct{}

func (c *Client) ClientDo(conf *Config) {
	for {
		fmt.Println("wait 3 second to get config")
		time.Sleep(time.Second * 3)
		c.Down(conf)
	}
}

func (c *Client) Down(conf *Config) error {
	rurl := conf.url
	u, _ := url.Parse(rurl)
	saveName := conf.out
	rsp, e := http.Get(rurl)
	if e != nil {
		fmt.Println(e)
		return e
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return errors.New("No new file to found")
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	if conf.befor != "" {
		fmt.Println("down load befor script")
		if e := c.DownFile("http://"+u.Host+"befor", conf.befor); e != nil {
			fmt.Println("error donwload befor script")
		}
		fmt.Println("run befor script")
		c.runScript(conf.befor)
	}
	ioutil.WriteFile(saveName, body, 0644)
	if conf.after != "" {
		fmt.Println("donwload after script")
		c.DownFile("http://"+u.Host+"after", conf.after)
		fmt.Println("run after script")
		c.runScript(conf.after)
	}
	return nil
}

func (c *Client) DownFile(url, saveName string) error {
	rsp, _ := http.Get(url)
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return errors.New("No new file to found")
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	ioutil.WriteFile(saveName, body, 0644)
	return nil
}

func (c *Client) runScript(scriptName string) error {
	cmd := exec.Command(scriptName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	conf := new(Config)
	e := getConfig(conf)
	if e != nil {
		fmt.Println(e)
		return
	}

	client := Client{}
	server := new(ServeHTTP)

	fmt.Println("here is:", conf.itype)
	server.paths = conf.paths
	if conf.itype == "server" {
		fmt.Println("Listen on:", conf.port)
		server.httpListen(conf)
	}
	if conf.itype == "client" {
		fmt.Println("Try to get:", conf.url)
		client.ClientDo(conf)
	}
}
