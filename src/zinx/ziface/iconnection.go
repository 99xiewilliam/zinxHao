package ziface

import "net"

type IConnection interface {
	//启动链接 当前链接准备开始工作
	Start()
	//停止链接 结束当前链接工作
	Stop()
	//获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前链接模块的链接IP
	GetConnID() uint32
	//获取远程客户端的TCP状态IP port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程的客户端
	SendMsg(uint32, []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})

	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error

//定义一个处理链接业务的方法
