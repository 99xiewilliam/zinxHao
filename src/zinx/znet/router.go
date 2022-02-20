package znet

import "go/src/zinx/ziface"

//实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写就好了
type BaseRouter struct {
}

//这里之所以BashHandle方法都为空
//是因为有的router不希望有PreHandle、PostHandle这两个业务
//所以Router全部继承来自BashRouter的好处是，不需要实现PreHandle和PostHandle
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (br *BaseRouter) Handle(request ziface.IRequest) {

}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
