package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil{
		fmt.Println("conn err",err)
		return
	}
	for{
		time.Sleep(1*time.Second)
		_,err:=conn.Write([]byte("Hello server\n"))
		if err != nil{
			fmt.Println("服务器已断开")
			break
		}
		buffer := make([]byte, 512)
		n, _ := conn.Read(buffer)
		fmt.Println(string(buffer[:n]))
	}



}
