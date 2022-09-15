package main

import (
	"flag"

	"github.com/Roeeeee/YA-RPC/server"
)

var listenIp string
var listenServerPort int
var listenClientPort int

func init() {
	flag.StringVar(&listenIp, "ip", "localhost", "设置registry IP，默认：127.0.0.1")
	flag.IntVar(&listenServerPort, "sport", 8000, "设置registry对server Port，默认：8000")
	flag.IntVar(&listenClientPort, "cport", 8001, "设置registry对client Port，默认8001")
}

func main() {
	flag.Parse()

	r := server.NewRegistry(listenIp, listenClientPort, listenServerPort)
	r.Start()
}
