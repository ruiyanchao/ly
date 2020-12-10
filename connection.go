package ly

type Connection interface {
	Read() (msg []interface{},err error)  // 读取信息
	Write(data string)                     // 发送信息
	Close()                               // 关闭链接
	GetName() string                      // 获取连接名
	GetConnID() int                       // 获取当前连接ID
	GetProperty() map[string]interface{}   // 获取连接属性
	OnConnStart()                           // 执行连接前方法
	OnConnStop()                            //执行断开连接方法
	SetOnConnStart(func(connection Connection)) // 设置连接开始方法
	SetOnConnStop(func(connection Connection))   // 设置断开连接方法
}