package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

//获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

//获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

//设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

//设置消息长度
func (m *Message) SetDataLLen(len uint32) {
	m.DataLen = len
}
