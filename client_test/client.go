package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Roeeeee/YA-RPC/client"
)

var registryIp string
var registryPort int

func init() {
	flag.StringVar(&registryIp, "rip", "localhost", "设置registry IP地址，默认：127.0.0.1")
	flag.IntVar(&registryPort, "rport", 8000, "设置registry Port，默认：8000")
}

type SumRequest struct {
	A float64
	B float64
}

type SumResponse struct {
	C float64
}

type UpperRequest struct {
	Str string
}

type UpperResponse struct {
	Str string
}

func main() {
	flag.Parse()

	c := client.NewClient(fmt.Sprintf("%s:%d", registryIp, registryPort))
	gob.Register(SumRequest{})
	gob.Register(SumResponse{})
	gob.Register(UpperRequest{})
	gob.Register(UpperResponse{})

	for i := 0; i < 10; i++ {
		sResI, err := c.RPCCall("Sum", SumRequest{A: 7.8, B: 3.5})
		if err != nil {
			log.Fatalln(err)
		}
		sRes := sResI.(SumResponse)
		log.Println(sRes)

		uResI, err := c.RPCCall("UpperCase", UpperRequest{Str: "distributed system"})
		if err != nil {
			log.Fatalln(err)
		}
		uRes := uResI.(UpperResponse)
		log.Println(uRes)

		<-time.After(time.Second * 10)
	}
}
