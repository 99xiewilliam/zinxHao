package main

import "go/src/zinx/znet"

func main() {
	//创建一个server具柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	//启动Server
	s.Serve()
}
