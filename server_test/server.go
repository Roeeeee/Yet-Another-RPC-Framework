package main

import (
	"encoding/gob"
	"flag"
	"log"
	"strings"

	"github.com/Roeeeee/Yet-Another-RPC-Framewor/YA-RPC/server"
)

type SumRequest struct {
	A float64
	B float64
}

type SumResponse struct {
	C float64
}

func Sum(req SumRequest) SumResponse {
	log.Println("Sum Run: ", req)
	res := SumResponse{req.A + req.B}
	return res
}

type UpperRequest struct {
	Str string
}

type UpperResponse struct {
	Str string
}

func UpperCase(req UpperRequest) UpperResponse {
	log.Println("UpperCase Run: ", req)
	res := UpperResponse{strings.ToUpper(req.Str)}
	return res
}

var registryIp string
var registryPort int
var serverIp string
var serverPort int

func init() {
	flag.StringVar(&registryIp, "rip", "localhost", "设置registry IP地址，默认：127.0.0.1")
	flag.IntVar(&registryPort, "rport", 8000, "设置registry Port，默认：8000")
	flag.StringVar(&serverIp, "sip", "localhost", "设置server监听IP，默认：127.0.0.1")
	flag.IntVar(&serverPort, "sport", 8002, "设置server监听Port，默认：8002")
}

func main() {
	flag.Parse()

	s := server.NewServer(serverIp, serverPort, registryIp, registryPort)

	s.Register("Sum", Sum)
	gob.Register(SumRequest{})
	gob.Register(SumResponse{})

	s.Register("UpperCase", UpperCase)
	gob.Register(UpperRequest{})
	gob.Register(UpperResponse{})

	s.RemoteRegister()
	s.Start()
}
