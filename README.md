# Go-JSON-RPC使用

JSON-RPC，是一个无状态且轻量级的远程过程调用传送协议，其传递内容主要以JSON数据为主，相较于一般的 RESTFul 通过 URL 地址，如 GET /student 调用远程服务器，JSON-RPC直接在内容中定义了想要调用的方法名称如，{"method": "Service": "getStudent"}，这也令开发者不会陷于该使用 PUT 还是 POST 的问题中。在RPC服务定义中主要定义一些数据结构及其相关的处理规则。在Golang中所有注册的RPC服务方法需要满足三个条件，第一方法有一个输入参数，第二方法有一个指针类型的输出参数，第三方法返回一个error类型的返回值，满足这三个条件即可注册为RPC服务方法。在Golang中如何使用请看以下内容...

#### 示例项目目录结构

```shell
$ tree
.
├── README.md
├── client
│   └── demo
│       └── demo.go
├── server
│   └── demo
│       └── demo.go
└── services
    └── demo
        └── demo.go
```

#### 定义RPC服务方法

```go
// services/demo/demo.go
package demo

import "errors"

// 定义RPC服务参数类型
type Args struct {
	A, B int
}

// 定义RPC服务
type Service struct {}

// Division 计算args参数之除法运算
func (*Service) Division(args Args, result *float64) error {
	if args.B == 0 {
		return errors.New("division by zero")
	}
	*result = float64(args.A) / float64(args.B)
	return nil
}
```

上面定义了一个名为 `Service` 的服务，注意这是在 `services/demo` 包中进行的定义的，这个服务有一个方法 `Division`，这个方法就是进行简单的业务处理也就是进行除法运算，并把运算结果写入到 `result` 变量内。

#### 注册RPC服务

编写RPC服务端程序如下：

```go
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
```

在这段程序中主要做了以下三件事情，第一就是注册RPC服务，第二就是通过 tcp 方式监听服务端口，以供外部调用之网络端口(说白了就是使用一个tcp网络协议8859端口提供RPC服务)，第三就是通过Goroutine的方式等待客户端连接到JSON-RPC服务。

那么这时一个简单的JSON-RPC服务也构建完毕了，服务构建完成之后那就需要使用，一个不使用的服务存在是没有意义的，我们可以通过telnet工具对其RPC服务端程序进行验证，也可以通过在其它项目中对RPC服务进行调用，为了简单起见，将客户端程序定义到 client/demo/demo.go 中，请直接向下看。

#### 调用RPC服务

调用RPC服务首先要将RPC服务程序运行，可以通过以下命令将其运行：

```go
$ go run server/server.go
```

##### 使用telnet工具进行RPC服务调用(也就是手动调用模式)

打开另一终端，输入telnet命令对RPC服务进行调用

```shell
$ telnet 127.0.0.1 8859
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
```

输入以JSON格式的参数数据，进行JSON-RPC服务调用：

```json
{"id": 123456, "method": "Service.Division", "params": [{"A": 88, "B": 6}]}
```

关于JSON格式的参数数据说明如下：

- `id`: 字段为调用编号，调用结束后服务端会原样返回(可供调用者根据此ID进行其它的业务处理)。
- `method`: 字段为要调用的RPC服务方法。
- `params`: 字段为调用RPC服务方法需要传递的参数。

-------------------------------------------------

如果调用成功返回JSON格式的结果如下所示：

```go
{"id":123456,"result":14.666666666666666,"error":null}
```

##### 使用RPC客户端进行RPC服务调用(也就是自动调用模式)

对于RPC服务也不可能只是通过telnet进行简单的调用，更多的是面向其它的服务，举个例了，比如说在一个大型电商项目中，有很多服务如 `订单`，`派送`，`短信`等等吧，这些服务的业务都是相对独立的，如用户下单是不是要生成订单，发送短信等等事件，那再一步说这些事件它肯定也不是同步的一个一个去完成，一般都会把下单之后的事件操作放入到消息中间件去异步处理，消息中间件再处理各个场景事件时再去调用其它的服务，这个时候就需要用到RPC调用，因为各个服务相对独立(如：不同地域服务，不同主机或端口号，不同服务器实例或Docker容器)，在这个示例中以 `client/demo/demo.go` 模拟其它项目中的调用，程序定义如下：

```go
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
```

运行客户端程序模拟其它项目对JSON-RPC服务调用：

```go
$ go run client/demo/demo.go
```

--------------------------------------------------------------
调用结果如下所示：

```shell
2020/04/06 18:50:47 JSON-RPC call result: 14.666667
```

这便是Golang中JSON-RPC如何使用的简单示例，当然你也可以使用gRPC框架进行RPC服务开发~

#### 示例代码

[Go-JSON-RPC](https://github.com/wumoxi/Go-JSON-RPC)