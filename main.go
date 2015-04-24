package main

import "fmt"

func main() {
	var conf Config
	e := getConfig(&conf)
	if iserr(e) {
		return
	}

	client := Client{}
	server := ServeHTTP{}
	client.conf = &conf
	server.conf = &conf

	if conf.itype == "server" {
		fmt.Println("Listen on:", conf.port)
		server.httpListen()
	}
	if conf.itype == "client" {
		fmt.Println("Try to get:", conf.url)
		client.ClientDo()
	}
}
