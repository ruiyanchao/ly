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
	i := 0
	for{
		time.Sleep(500*time.Microsecond)
		_,err:=conn.Write([]byte("Hello server\n"))
		if err != nil{
			fmt.Println("服务器已断开")
			break
		}
		buffer := make([]byte, 512)
		n, _ := conn.Read(buffer)
		if n > 0{
			i++
			fmt.Printf("第 %d 次发送 \n" ,i)
		}

		//fmt.Println(string(buffer[:n]))
	}



}
