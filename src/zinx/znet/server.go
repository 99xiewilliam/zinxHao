package znet

import (
	"fmt"
	"go/src/zinx/utils"
	"go/src/zinx/ziface"
	"net"
)

//IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定ip版本
	IPVersion string
	//服务器监听ip
	IP string
	//服务器监听端口
	Port int

	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//该server的链接管理器
	ConnMgr ziface.IConnManager

	//该Server穿gain链接之后自动调用Hook函数--OnConnStart
	OnConnStart func(conn ziface.IConnection)
	//该Server销毁链接之前自动调用Hook函数--OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

//当前客户端链接的所绑定的handle api(目前这个handle是写死的，以后优化应该由用户自定义)
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	//回显业务
//	fmt.Println("[Conn Handle] CallbackToClient...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err", err)
//		return errors.New("CallBackToClient error")
//	}
//
//	return nil
//}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server name:%s, listenner at IP: %s, Port:%d is starting",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("{Start} Server listenner at IP :%s, Port %d is starting \n", s.IP, s.Port)
	go func() {
		s.MsgHandler.StartWorkerPool()
		//1。获取tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		//2。监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ, ", s.Name, "succ, Listening...")
		var cid uint32
		cid = 0
		//3。阻塞等待客户端链接，处理客户端链接业务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO给客户端相应一个超出最大的连接错误包
				fmt.Println("====> Too many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//已经与客户端建立连接，做一些业务，做一个最基本的最大512字节长度的回显业务
			//将处理心链接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cmt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err", err)
			//		}
			//		fmt.Printf("recv client buf %s, cnt %d\n", buf, cmt)
			//		//会显功能
			//		if _, err := conn.Write(buf[:cmt]); err != nil {
			//			fmt.Println("wirte book buf err", err)
			//			continue
			//		}
			//	}
			//}()
		}
	}()
}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些已经开始的链接信息 进行停止或者回收
	fmt.Println("[STOP] Zinx server name", s.Name)
	s.ConnMgr.ClearConn()

}

func (s *Server) Serve() {
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("add Router succ!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//初始化Server模块方法

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsghandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

//注册OnConnStart钩子函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop钩子函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()")
		s.OnConnStop(conn)
	}
}
