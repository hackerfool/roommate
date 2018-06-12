package main

import (
	"api"
	"flag"
)

var (
	listenAddr = flag.String("a", "", "listen addr")
	listenPort = flag.String("p", "8080", "listen port")
	mode       = flag.String("m", "release", "run mode:test,debug,release")
	showHelp   = flag.Bool("h", false, "show this help list")
)

func main() {
	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}

	// blockchain.Init()
	// ch := blockchain.NewBlockChain()
	// ch.AddBlock([]byte("Hello world!"))
	// for _, i := range ch.Node {
	// 	fmt.Printf("Data: %s\r\nHash: %x\r\nPrevHash: %x\r\n\r\n\r\n", i.Data, i.Hash, i.Prevhash)
	// }

	// api.RegisterHandle()
	api.Run(*listenAddr, *listenPort, *mode)
	// mlog.SetLevel(mlog.DEBUG)
	// mlog.Info("server start ok,listen on ", *listenPort, ".")
	// err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil)
	// mlog.Fatal(err)

	return
}
