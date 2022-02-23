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
//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle")
//
//	//先读取客户端数据，然后会写
//	fmt.Println("recv from client:msgID = ", request.GetMsgID(), "data=",
//		string(request.GetData()))
//
//	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
//	if err != nil {
//		fmt.Println(err)
//	}
//}

//在处理conn业务的主方法hook
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router handle")
	fmt.Println("recv from client:msgID = ", request.GetMsgID(), "data=",
		string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

//在处理conn业务的主方法hook
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter handle")
	fmt.Println("recv from client:msgID = ", request.GetMsgID(), "data=",
		string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello Zinx !!!"))
	if err != nil {
		fmt.Println(err)
	}
}

//在处理conn业务的钩子方法hook
//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping... ping... ping..."))
//	if err != nil {
//		fmt.Println("call back after ping... ping... ping... error")
//	}
//}

func main() {
	//创建一个server具柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.8]")

	//给框架添加一个router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//启动Server
	s.Serve()
}
