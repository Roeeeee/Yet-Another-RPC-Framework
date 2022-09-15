package server

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/Roeeeee/YA-RPC/client"
	"github.com/Roeeeee/YA-RPC/comm"
)

type Server struct {
	IP    string
	Port  int
	RIP   string
	RPort int
	Funcs map[string]reflect.Value
}

func NewServer(ip string, port int, rip string, rport int) *Server {
	server := &Server{
		IP:    ip,
		Port:  port,
		RIP:   rip,
		RPort: rport,
		Funcs: make(map[string]reflect.Value),
	}
	return server
}

func (s *Server) Register(funcName string, funcBody interface{}) {
	if _, ok := s.Funcs[funcName]; ok {
		log.Printf("Warning: Func %s had been registered!\n", funcName)
	}
	method := reflect.ValueOf(funcBody)
	s.Funcs[funcName] = method
}

func (s *Server) RemoteRegister() {
	c := client.NewClient(fmt.Sprintf("%s:%d", s.RIP, s.RPort))
	rgReq := &RegisterRequest{
		IP:    s.IP,
		Port:  s.Port,
		Funcs: make([]string, 0),
	}
	for k := range s.Funcs {
		rgReq.Funcs = append(rgReq.Funcs, k)
	}

	rgResI, err := c.Call("Register", rgReq)
	if err != nil {
		log.Fatalln(err)
	}
	rgRes := rgResI.(RegisterResponse)
	if !rgRes.Ok {
		log.Fatalln("Remote Register Fail")
	} else {
		log.Println("Register Success!")
	}

	// 开一个协程定时发送心跳信息。
	go func() {
		hbReq := &HeartbeatRequest{
			IP:   s.IP,
			Port: s.Port,
		}
		for {
			// fmt.Println("Send Heartbeat")
			hbResI, err := c.Call("Heartbeat", hbReq)
			if err != nil {
				log.Fatalln(err)
			}
			hbRes := hbResI.(HeartbeatResponse)
			if !hbRes.Ok {
				log.Fatalln("Hearbeat Fail!")
			}
			<-time.After(time.Millisecond * 100)
		}
	}()
}

func (s *Server) Execute(funcName string, arg interface{}) (interface{}, error) {
	f, ok := s.Funcs[funcName]
	if !ok {
		err := fmt.Errorf("can't find func: %s", funcName)
		return nil, err
	}
	inArg := []reflect.Value{reflect.ValueOf(arg)}
	ret := f.Call(inArg)
	res := ret[0].Interface()
	return res, nil
}

func (s *Server) Handler(conn net.Conn) {
	reqPkg, err := comm.Read(conn)
	if err != nil {
		log.Println(conn.RemoteAddr(), " read error: ", err)
		return
	}
	// fmt.Printf("Get RPC-Call %s from %s\n", reqPkg.FuncName, conn.RemoteAddr())
	res, err := s.Execute(reqPkg.FuncName, reqPkg.Arg)
	resPkg := comm.NewRPCPackage(reqPkg.FuncName, res, err)
	comm.Write(conn, resPkg)
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	log.Println("Server Start!")
	// 每到达一个请求就单独开一个协程处理，主协程不会阻塞。
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.Handler(conn)
	}
}
