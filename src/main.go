package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listenPort = flag.Int("p", 8080, "listen port")
	showHelp   = flag.Bool("h", false, "show this help list")
)

func main() {
	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}

	registerHandle()
	logger.info("server start ok,listen on ", *listenPort, ".")
	err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil)
	log.Fatal(err)
	return
}
