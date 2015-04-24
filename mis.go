package main

import "fmt"

func iserr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
