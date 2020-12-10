package LY

import (
	"fmt"
	"net"
	"sync"
)


type tcp struct {
	sync.Mutex
	// 当前连接
	conn net.Conn
	// 连接协议
	protocol protocol
	//包大小
	MaxBuffer int
	//包场地
	MaxPageSize int
	//协议名
	name string
	//链接属性
	property map[string]interface{}
	//连接开始方法
	onConnStart func(connection Connection)
	//链接断开方法
	onConnStop func(connection Connection)
	//连接ID
	connID int
	//当前 buff
	currentBuf []byte
}

/**
 * 读取信息
 */
func (tc *tcp)Read()(msg []interface{} , err error){
	//TODO SSL
	//读取conn 传来的信息
	tc.currentBuf = make([]byte,tc.MaxBuffer )
	n, err :=tc.conn.Read(tc.currentBuf)
	tc.currentBuf = tc.currentBuf[:n]
	//初始化一条信息buff
	oneRequestBuf := make([]byte,len(tc.currentBuf) )
	if err != nil{
		fmt.Println("read buf err",err)
		return
	}
	if tc.protocol != nil{
		for{
			//如果当前信息为空 跳出
			if string(tc.currentBuf) == ""{
				break
			}
			// 通过协议解析 获取 当前一条buf长度
			currentBufLen := tc.protocol.Input(tc.currentBuf)
			if currentBufLen == 0{
				return
			}else if currentBufLen > 0 && currentBufLen <= tc.MaxPageSize{
				if currentBufLen > len(tc.currentBuf){
					break
				}
			}
			//计算获取 一条buf 重置当前buf
			if len(tc.currentBuf) == currentBufLen{
				oneRequestBuf = tc.currentBuf
				tc.currentBuf = []byte("")
			}else{
				oneRequestBuf = tc.currentBuf[:currentBufLen]
				tc.currentBuf  = tc.currentBuf[currentBufLen:]
			}
			currentBufLen = 0
			//防止粘包
			msg = append(msg,string(oneRequestBuf))

		}

	}
	return
}

/**
 *发送
 */
func (tc *tcp)Write(data string){
	tc.conn.Write([]byte(data))
}

/**
 *关闭链接
 */
func (tc *tcp)Close(){
	tc.conn.Close()
}

/**
 *获取链接名
 */
func (tc *tcp)GetName() string{
	return tc.name
}

/**
 * 获取连接属性
 */
func (tc *tcp)GetProperty() map[string]interface{}{
	return tc.property
}

/**
 * 执行开始方法
 */
func (tc *tcp)OnConnStart(){
	if tc.onConnStart != nil {
		tc.onConnStart(tc)
	}
}

/**
 * 执行关闭方法
 */
func (tc *tcp)OnConnStop(){
	if tc.onConnStop != nil {
		tc.onConnStop(tc)
	}
}

/**
 * 设置开始方法
 */
func (tc *tcp)SetOnConnStart(function func(connection Connection)){
	tc.onConnStart = function
}

/**
 * 设置关闭方法
 */
func (tc *tcp)SetOnConnStop(function func(connection Connection)){
	tc.onConnStop = function
}

/**
 * 获取连接ID
 */
func(tc *tcp)GetConnID() int{
	return tc.connID
}




