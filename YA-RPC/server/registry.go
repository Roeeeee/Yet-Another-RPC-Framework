package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Roeeeee/Yet-Another-RPC-Framewor/YA-RPC/client"
)

var serverMap map[string]*rserver
var funcMap map[string]*rfunction
var serverMapMu sync.RWMutex
var funcMapMu sync.RWMutex

type rserver struct {
	ip     string
	port   int
	addr   string
	funcs  map[string]*rfunction
	isLive chan bool
}

/*
	rserver超时，清除serverMap和funcMap对其的引用，结束heartbeatTimer协程
	GC对其回收
*/
func (rs *rserver) heartbeatTimer() {
	for {
		select {
		case <-rs.isLive:
			// 什么也不做
		case <-time.After(time.Millisecond * 300):
			log.Printf("%s is timeout\n", rs.addr)

			serverMapMu.Lock()
			delete(serverMap, rs.addr)
			serverMapMu.Unlock()

			funcMapMu.Lock()
			for _, v := range rs.funcs {
				delete(v.servers, rs.addr)
			}
			funcMapMu.Unlock()
			return
		}
	}
}

type rfunction struct {
	funcName string
	servers  map[string]*rserver
	mu       *sync.RWMutex
}

func (rf *rfunction) randGetServer() string {
	// 随机返回一个服务器。
	cnt := rand.Intn(len(rf.servers))
	for k := range rf.servers {
		if cnt == 0 {
			return k
		}
		cnt--
	}
	panic("unreachable")
}

/*
	rserver一旦创建则为只读，只能由GC进行回收
	rserver创建后，只有serverMap，funcMap，heartbeatTimer对其引用
*/
func Register(req RegisterRequest) RegisterResponse {
	s := &rserver{
		ip:     req.IP,
		port:   req.Port,
		addr:   fmt.Sprintf("%s:%d", req.IP, req.Port),
		funcs:  make(map[string]*rfunction),
		isLive: make(chan bool, 1),
	}
	funcMapMu.Lock()
	for _, fname := range req.Funcs {
		f, ok := funcMap[fname]
		if !ok {
			f = &rfunction{
				funcName: fname,
				servers:  make(map[string]*rserver),
				mu:       new(sync.RWMutex),
			}
			funcMap[fname] = f
		}
		f.servers[s.addr] = s
		s.funcs[fname] = f
		log.Printf("%s register %s\n", s.addr, fname)
	}
	funcMapMu.Unlock()

	serverMapMu.Lock()
	serverMap[s.addr] = s
	serverMapMu.Unlock()

	go s.heartbeatTimer()

	return RegisterResponse{Ok: true}
}

type RegisterRequest struct {
	IP    string
	Port  int
	Funcs []string
}

type RegisterResponse struct {
	Ok bool
}

/*
	通过isLive管道来实现心跳保活。
*/
func Heartbeat(req HeartbeatRequest) HeartbeatResponse {
	serverMapMu.RLock()
	s, ok := serverMap[fmt.Sprintf("%s:%d", req.IP, req.Port)]
	if !ok {
		serverMapMu.RUnlock()
		return HeartbeatResponse{Ok: false}
	}
	s.isLive <- true
	serverMapMu.RUnlock()
	return HeartbeatResponse{Ok: true}
}

type HeartbeatRequest struct {
	IP   string
	Port int
}

type HeartbeatResponse struct {
	Ok bool
}

func GetServer(req client.GetServerRequest) client.GetServerResponse {
	funcMapMu.RLock()
	f, ok := funcMap[req.FuncName]
	funcMapMu.RUnlock()
	if !ok {
		return client.GetServerResponse{Ok: false}
	}

	f.mu.RLock()
	if len(f.servers) == 0 {
		f.mu.RUnlock()
		return client.GetServerResponse{Ok: false}
	}
	addr := f.randGetServer()
	f.mu.RUnlock()

	return client.GetServerResponse{ServerAddr: addr, Ok: true}
}

type Registry struct {
	ip    string
	sport int
	cport int
}

func NewRegistry(ip string, cport int, sport int) *Registry {
	r := &Registry{
		ip:    ip,
		sport: sport,
		cport: cport,
	}
	return r
}

func (r *Registry) Start() {
	sServer := NewServer(r.ip, r.sport, r.ip, r.sport)
	sServer.Register("Register", Register)
	sServer.Register("Heartbeat", Heartbeat)
	sServer.Register("GetServer", GetServer)
	sServer.Start()
}

func init() {
	gob.Register(RegisterRequest{})
	gob.Register(RegisterResponse{})
	gob.Register(HeartbeatRequest{})
	gob.Register(HeartbeatResponse{})

	serverMap = make(map[string]*rserver)
	funcMap = make(map[string]*rfunction)
}
