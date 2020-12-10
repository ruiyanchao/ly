package LY

type Connection interface {
	Read() (msg []interface{},err error)
	Write(data string)
	Close()
	GetName() string
	GetConnID() int
	GetProperty() map[string]interface{}
	OnConnStart()
	OnConnStop()
	SetOnConnStart(func(connection Connection))
	SetOnConnStop(func(connection Connection))
}