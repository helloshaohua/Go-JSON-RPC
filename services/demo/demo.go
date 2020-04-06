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