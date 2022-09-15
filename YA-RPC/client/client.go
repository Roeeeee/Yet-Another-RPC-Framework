package client

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Roeeeee/Yet-Another-RPC-Framewor/YA-RPC/comm"
)

type Client struct {
	Server string
}

// 实例化一个client
func NewClient(s string) *Client {
	c := &Client{
		Server: s,
	}
	return c
}

// 以直接调用方式进行RPC请求。
func (c *Client) Call(funcName string, req interface{}) (interface{}, error) {
	conn, err := net.Dial("tcp", c.Server)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()

	reqPkg := comm.NewRPCPackage(funcName, req, nil)
	err = comm.Write(conn, reqPkg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resPkg, err := comm.Read(conn)
	if err != nil || resPkg.Err != nil {
		log.Println(err)
		return nil, err
	}
	return resPkg.Arg, nil
}

var func2server map[string]string
var func2serverMu sync.RWMutex

func addServer(funcName string, serverAddr string) {
	func2serverMu.Lock()
	func2server[funcName] = serverAddr
	func2serverMu.Unlock()
}

func deleteServer(funcName string) {
	func2serverMu.Lock()
	delete(func2server, funcName)
	func2serverMu.Unlock()
}

func (c *Client) RPCCall(funcName string, req interface{}) (interface{}, error) {
	func2serverMu.RLock()
	serverAddr, ok := func2server[funcName]
	func2serverMu.RUnlock()
	if !ok {
		// 表中没有提供func的server，向registry查询
		addr, err := c.getServer(funcName)
		if err != nil {
			return nil, err
		}
		serverAddr = addr
	}

	cc := NewClient(serverAddr)
	res, err := cc.Call(funcName, req)
	if err != nil {
		/*
			RPC首次调用失败
			从表中删去无效的server
			再试一次
		*/
		deleteServer(funcName)
		newAddr, err := c.getServer(funcName)
		if err != nil {
			return nil, err
		}
		cc = NewClient(newAddr)
		return cc.Call(funcName, req)
	}
	return res, nil
}

func (c *Client) getServer(funcName string) (string, error) {
	gsReq := &GetServerRequest{
		FuncName: funcName,
	}
	gsResI, err := c.Call("GetServer", gsReq)
	if err != nil {
		return "", err
	}
	gsRes := gsResI.(GetServerResponse)
	if !gsRes.Ok {
		err = fmt.Errorf("can't find func: %s", funcName)
		return "", err
	}

	addServer(funcName, gsRes.ServerAddr)

	return gsRes.ServerAddr, nil
}

type GetServerRequest struct {
	FuncName string
}

type GetServerResponse struct {
	ServerAddr string
	Ok         bool
}

func init() {
	gob.Register(GetServerRequest{})
	gob.Register(GetServerResponse{})
	func2server = make(map[string]string)
}
