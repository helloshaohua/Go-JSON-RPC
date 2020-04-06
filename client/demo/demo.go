package main

import (
	"GO-JSON-RPC/services/demo"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

func main() {
	// Errors recover
	defer func() {
		if e := recover(); e != nil {
			log.Printf("demo JSON RPC client has errors: %s\n", e)
		}
	}()
	// 连接到远程RPC服务
	conn, err := net.Dial("tcp", ":8859")
	if err != nil {
		panic(err)
	}

	// 声明JSON-RPC调用回复结果变量
	var reply float64

	// 创建JSON-RPC客户端
	client := jsonrpc.NewClient(conn)

	// 调用RPC服务方法
	err = client.Call("Service.Division", demo.Args{A: 88, B: 6}, &reply)
	if err != nil {
		panic(err)
	}
	log.Printf("call result: %f\n", reply)
}