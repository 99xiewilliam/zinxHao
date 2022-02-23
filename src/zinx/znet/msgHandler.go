package znet

import (
	"fmt"
	"go/src/zinx/utils"
	"go/src/zinx/ziface"
	"strconv"
)

type MsgHandle struct {
	//存放每个MsgID 对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker去任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsghandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is not found! need Register")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的API处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}

//启动一个worker工作池（开启工作池的动作只能发生一次，一个zinx框架只能又一个worker工作池）
func (mh *MsgHandle) StartWorkerPool() {

	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//1个worker被启动
		//给当前的worker对应的channel消息队列开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当时的worker，阻塞等待消息从channel传递出来
		go mh.StartOneWorker(i, mh.TaskQueue[i])

	}
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started...")

	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不通过的worker
	//根据客户端建立的ConnID进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgID(), "to WorkerID", workerID)

	//将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
