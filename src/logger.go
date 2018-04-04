package main

import (
	"fmt"
	"log"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	white  = 37
)

type nlog struct {
}

var (
	logger = &nlog{}
)

func (*nlog) debug(v ...interface{}) {
	setConsleColor(green)
	defer unsetConsleColor()

	log.Println("[D]", fmt.Sprint(v...))
	return
}

func (*nlog) info(v ...interface{}) {
	log.Println("[I]", fmt.Sprint(v...))
	return
}

func (*nlog) error(v ...interface{}) {
	setConsleColor(red)
	defer unsetConsleColor()

	log.Println("[E]", fmt.Sprint(v))
	return
}

func setConsleColor(color int) {
	fmt.Printf("\033[%dm", color)
	return
}

func unsetConsleColor() {
	setConsleColor(0)
	return
}
