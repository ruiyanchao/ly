package main

import (
	"fmt"
	"strings"
	"LY"
)

func main(){
	w := LY.Worker{
		SocketName:"tcp://127.0.0.1:8080",
		Protocol: &MyProtocol{
			name:"rick",
		},
	}
	w.OnMessage = func(con LY.Connection,msg interface{}) {
		fmt.Println(msg)
		con.Write("hello client")
	}
	w.OnConnStart = func(con LY.Connection) {
		cid := con.GetConnID()
		fmt.Printf("get client [%d] \n",cid)
	}
	w.OnConnStop = func(con LY.Connection) {
		cid := con.GetConnID()
		fmt.Printf("client [%d] leave \n",cid)
	}

	w.RunAll()
}

type MyProtocol struct {
	name string
}
func (mp *MyProtocol)Input(buf []byte)( l int){
	request := string(buf)
	pos:= strings.Index(request,"\n")
	if pos == -1{
		return 0
	}
	return pos+1
}

func (mp *MyProtocol)Encode(buf []byte){
	return
}

func (mp *MyProtocol)Decode(data string){
	return
}

func (mp *MyProtocol)GetName() string{
	return mp.name
}


