package znet

import (
	"errors"
	"fmt"
	"go/src/zinx/utils"
	"go/src/zinx/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer ziface.IServer
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前链接状态
	isClosed bool

	//当前链接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc

	//告知当前链接已经推出的/停止 channel
	ExitChan chan bool

	//无缓冲的管道和对应的处理业务API关系
	msgChan chan []byte

	//消息的管理MsgID 和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	fmt.Println("123123")
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "[Reader is exit], remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//buf := make([]byte, 512)
		//_, err := c.Conn.Read(buf)
		//
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		////调用当前链接所绑定的HandleAPI
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID ", c.ConnID, "handle is error", err)
		//	break
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端Msg Head二级制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据dataLen再次读取Data，放在msg。data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		////执行注册路由方法
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)
		//根据绑定好的MsgID 找到对应处理的API业务 执行
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

		//从路由中，找到注册绑定的Conn对应的router调用

	}

}

//提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，在发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDataLen MsgID data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))

	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Write msg id", msgId, "error:", err)
	//	return errors.New("conn Write error")
	//}

	c.msgChan <- binaryMsg
	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

//写消息的goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error", err)
				return
			}
		case <-c.ExitChan:
			//Reader 已经推出，此时Writer也要退出
			return
		}
	}
}

//启动链接 当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	//TODO 启动当前从链接写数据的业务
	go c.StartReader()

	go c.StartWriter()

	//按照开发者传递进来的，创建链接之后需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)

}

//停止链接 结束当前链接工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID=", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//调用开发者注册的销毁链接之间 需要执行的Hook函数
	c.TcpServer.CallOnConnStop(c)
	//关闭socket链接
	c.Conn.Close()
	c.ExitChan <- true

	//将当前链接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接IP
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
