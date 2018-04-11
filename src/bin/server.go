package main

import (
	"api"
	"flag"
	"fmt"
	"mlog"
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

	api.RegisterHandle()
	mlog.SetLevel(mlog.DEBUG)
	mlog.Info("server start ok,listen on ", *listenPort, ".")
	err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil)
	mlog.Fatal(err)

	return
}
