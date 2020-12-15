package ly

import (
	"net"
	"sync"
)

type udp struct {
	sync.Mutex
	//连接开始方法
	onConnStart func(connection Connection)
	//链接断开方法
	onConnStop func(connection Connection)
	conn *net.UDPConn
	protocol protocol
	property map[string]interface{}
	name string
	connID int
	
}

/**
 * 读取信息
 */
func (ud *udp)Read()(msg []interface{} , err error){
	return 
}

/**
 *发送
 */
func (ud *udp)Write(data string){
	ud.setProperty("send","incr")
	if ud.protocol != nil{
		ud.conn.Write([]byte(ud.protocol.Encode(data)))
	}else{
		ud.conn.Write([]byte(data))
	}
}

/**
 *关闭链接
 */
func (ud *udp)Close(){
	ud.conn.Close()
}

/**
 *获取链接名
 */
func (ud *udp)GetName() string{
	return ud.name
}

/**
 * 获取连接属性
 */
func (ud *udp)GetProperty() map[string]interface{}{
	return ud.property
}

/**
 * 执行开始方法
 */
func (ud *udp)OnConnStart(){
	if ud.onConnStart != nil {
		ud.onConnStart(ud)
	}
}

/**
 * 执行关闭方法
 */
func (ud *udp)OnConnStop(){
	if ud.onConnStop != nil {
		ud.onConnStop(ud)
	}
}

/**
 * 设置开始方法
 */
func (ud *udp)SetOnConnStart(function func(connection Connection)){
	ud.onConnStart = function
}

/**
 * 设置关闭方法
 */
func (ud *udp)SetOnConnStop(function func(onnection Connection)){
	ud.onConnStop = function
}

/**
 * 获取连接ID
 */
func(ud *udp)GetConnID() int{
	return ud.connID
}

/**
 * 设置发送
 */
func(ud *udp)setProperty(key string,value string){
	ud.Lock()
	if value == "incr"{
		if _,ok := ud.property[key];ok{
			ud.property[key] = ud.property[key].(int) + 1
		}else{
			ud.property[key] = 1
		}
	}else{
		ud.property[key] = value
	}

	ud.Unlock()
}
