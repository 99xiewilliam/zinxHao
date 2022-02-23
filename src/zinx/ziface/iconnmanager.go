package ziface

/*
链接管理模块
*/

type IConnManager interface {
	//添加
	Add(conn IConnection)
	//删除
	Remove(conn IConnection)
	//根据ID获取
	Get(connID uint32) (IConnection, error)
	//得到当前链接总数
	Len() int
	//清楚并终止链接
	ClearConn()
}
