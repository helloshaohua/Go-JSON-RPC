// server/demo/demo.go
package main

import (
	"GO-JSON-RPC/services/demo"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main()  {
	// Errors recover
	defer func() {
		if e := recover(); e != nil {
			log.Printf("demo JSON RPC server has errors: %s\n", e)
		}
	}()
	// 服务端注册RPC服务
	err := rpc.Register(new(demo.Service))
	if err != nil {
		panic(err)
	}

	// 监听服务端口
	listener, err := net.Listen("tcp", ":8859")
	if err != nil {
		panic(err)
	}

	// 接受连接并处理服务调用
	for {
		accept, err := listener.Accept()
		if err != nil {
			log.Printf("accpet has errors: %s\n", err)
			continue
		}
		go jsonrpc.ServeConn(accept)
	}
}

