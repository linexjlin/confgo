package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"time"
)

type Client struct {
	conf *Config
}

func (c *Client) ClientDo() {
	for {
		fmt.Println("wait 3 second to get config")
		time.Sleep(time.Second * 3)
		iserr(c.Down())
	}
}

//download configure file
func (c *Client) Down() error {
	rurl := c.conf.url
	u, _ := url.Parse(rurl)
	saveName := c.conf.out
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

	//down the befor script to run
	if c.conf.befor != "" {
		fmt.Println("down load befor script")
		if e := c.DownFile("http://"+u.Host+"/"+c.conf.befor, c.conf.befor); e != nil {
			fmt.Println("error donwload befor script")
		}
		fmt.Println("run befor script")
		if iserr(c.runScript(c.conf.befor)) {
			fmt.Println("run befor script error")
		}
	}

	//Write the confiure file
	ioutil.WriteFile(saveName, body, 0644)

	//down the after script to run
	if c.conf.after != "" {
		fmt.Println("donwload after script")
		c.DownFile("http://"+u.Host+"/"+c.conf.after, c.conf.after)
		fmt.Println("run after script")
		if iserr(c.runScript(c.conf.after)) {
			fmt.Println("run after script error")
		}
	}

	return nil
}

//common file donwload function
func (c *Client) DownFile(url, saveName string) error {
	rsp, e := http.Get(url)
	if iserr(e) {
		return e
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return errors.New("No new file to found")
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	ioutil.WriteFile(saveName, body, 0644)
	return nil
}

//common run script function
func (c *Client) runScript(scriptName string) error {
	cmd := exec.Command(scriptName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
