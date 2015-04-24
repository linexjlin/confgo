package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"
)

type Config struct {
	paths []string
	port  string
	url   string
	befor string
	after string
	out   string
	itype string
}

func getConfig(conf *Config) error {
	fpath := "."
	fport := "60000"
	furl := "http://locahost:60000/myfile.txt"
	fbefor := ""
	fafter := ""
	fout := "out.txt"
	ftype := "client"
	flag.StringVar(&ftype, "type", "", "-type=client")
	flag.StringVar(&fpath, "path", ".", `-path ="/home/path1, /home/path2"`)
	flag.StringVar(&fport, "port", "60000", `-port="60000"`)
	flag.StringVar(&furl, "url", "", `-url="http://localhost:60000/myfile.txt"`)
	flag.StringVar(&fbefor, "befor", "", `-befor="befor.bat"`)
	flag.StringVar(&fafter, "after", "", `-after="after.bat"`)
	flag.StringVar(&fout, "out", "out.txt", `-out="./out.txt"`)
	flag.Parse()
	conf.paths = strings.Split(fpath, ",")
	conf.url = furl
	conf.itype = ftype
	if conf.itype == "" {
		return errors.New("Please spacify type")
	}
	conf.out = fout
	conf.port = fport
	conf.after = fafter

	if conf.itype == "server" {
		fmt.Println("Here is: ", conf.itype, "program")
		fmt.Println("Will Listen on port:", conf.port)
		fmt.Println("Found files from ", conf.paths)

		if conf.befor != "" {
			fmt.Println("Will the befor script is ", conf.befor)
		}

		if conf.after != "" {
			fmt.Println("Will the befor script is ", conf.after)
		}
	}

	if conf.itype == "client" {
		fmt.Println("Here is:", conf.itype, "program")

		if conf.url == "" {
			return errors.New("Url no given")
		}

		fmt.Println("Will get configure from ", conf.url, ", and save to ", conf.out)
		u, _ := url.Parse(conf.url)

		if conf.befor != "" {
			fmt.Println("Will get befor script from ", "http://"+u.Host+"/"+conf.befor)
		}

		if conf.after != "" {
			fmt.Println("Will get after script from ", "http://"+u.Host+"/"+conf.after)
		}
	}

	return nil
}
