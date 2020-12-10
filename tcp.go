package LY

import (
	"fmt"
	"net"
	"sync"
)


type tcp struct {
	sync.Mutex
	conn net.Conn
	protocol protocol
	defaultMaxBuffer int
	defaultMaxPageSize int
	name string
	property map[string]interface{}
	onConnStart func(connection Connection)
	onConnStop func(connection Connection)
	connID int
	currentBuf []byte
	onMessage func(con Connection, msg interface{})
}

type ConnInfo struct {
	name string
	value interface{}
}

func (tc *tcp)Read()(msg []interface{} , err error){
	tc.currentBuf = make([]byte,tc.defaultMaxBuffer )
	n, err :=tc.conn.Read(tc.currentBuf)
	tc.currentBuf = tc.currentBuf[:n]
	oneRequestBuf := make([]byte,len(tc.currentBuf) )
	if err != nil{
		fmt.Println("read buf err",err)
		return
	}
	if tc.protocol != nil{
		for{
			if string(tc.currentBuf) == ""{
				break
			}
			currentBufLen := tc.protocol.Input(tc.currentBuf[:n])
			if currentBufLen == 0{
				return
			}else if currentBufLen > 0 && currentBufLen <= tc.defaultMaxPageSize{
				if currentBufLen > len(tc.currentBuf){
					break
				}
			}
			if len(tc.currentBuf) == currentBufLen{
				oneRequestBuf = tc.currentBuf
				tc.currentBuf = []byte("")

			}else{
				oneRequestBuf = tc.currentBuf[:currentBufLen]
				tc.currentBuf  = tc.currentBuf[currentBufLen:]
			}
			currentBufLen = 0
			msg = append(msg,string(oneRequestBuf))

		}

	}
	return
}

func (tc *tcp)Write(data string){
	tc.Lock()
	if _,ok:= tc.property["send"];ok{
		tc.property["send"] = tc.property["send"].(int)+1
	}else{
		tc.property["send"] = 1
	}
	tc.Unlock()
	tc.conn.Write([]byte(data))
}

func (tc *tcp)Close(){
	tc.conn.Close()
}

func (tc *tcp)GetName() string{
	return tc.name
}

func (tc *tcp)GetProperty() map[string]interface{}{
	return tc.property
}

func (tc *tcp)OnConnStart(){
	if tc.onConnStart != nil {
		tc.onConnStart(tc)
	}
}

func (tc *tcp)OnConnStop(){
	if tc.onConnStop != nil {
		tc.onConnStop(tc)
	}
}

func (tc *tcp)SetOnConnStart(function func(connection Connection)){
	tc.onConnStart = function
}

func (tc *tcp)SetOnConnStop(function func(connection Connection)){
	tc.onConnStop = function
}

func(tc *tcp)GetConnID() int{
	return tc.connID
}




