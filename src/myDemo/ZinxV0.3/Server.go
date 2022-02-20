package main

import (
	"fmt"
	"go/src/zinx/ziface"
	"go/src/zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//在处理conn业务之前的钩子方法Hook
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router Prehanlde")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

//在处理conn业务的主方法hook
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping... ping..."))
	if err != nil {
		fmt.Println("call back ping... ping... ping... error")
	}
}

//在处理conn业务的钩子方法hook
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping... ping... ping..."))
	if err != nil {
		fmt.Println("call back after ping... ping... ping... error")
	}
}

func main() {
	//创建一个server具柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.3]")

	//给框架添加一个router
	s.AddRouter(&PingRouter{})

	//启动Server
	s.Serve()
}
