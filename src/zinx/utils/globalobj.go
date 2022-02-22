package utils

import (
	"encoding/json"
	"go/src/zinx/ziface"
	"io/ioutil"
)

/*
存储一切有关Zinx的全局参数，供其他模块使用
一些蚕食是可以通过
Zinx。json用户进行配置
*/

type GlobalObj struct {
	TcpServer ziface.IServer //当亲啊Zinx全局Server对象
	Host      string         //服务器主机监听ip
	TcpPort   int            //服务器主机监听端口号
	Name      string         //服务器名称

	Version        string //当前Zinx版本号
	MaxConn        int    //服务器主机允许的最大链接数
	MaxPackageSize uint32 //Zinx框架数据包的最大值
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	//如果配置文件没有家在，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//GlobalObject.Reload()
}
