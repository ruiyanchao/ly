## LY
LY是一套轻量级的Go语言TCP框架（参考了 WorkerMan zinx 等优秀作品）。

## 快速开始
### server
1.自定义协议
```go
type MyProtocol struct {
	name string //定义协议名
}

// 获取一条信息长度
func (mp *MyProtocol)Input(buf []byte)( l int){
	request := string(buf)
	pos:= strings.Index(request,"\n")
	if pos == -1{
		return 0
	}
	return pos+1
}

//加密 信息 用于发送
func (mp *MyProtocol)Encode(data string){
	return
}

//解密 信息 用于读取
func (mp *MyProtocol)Decode(buf []byte){
	return
}

func (mp *MyProtocol)GetName() string{
	return mp.name
}
```
2.配置对应闭包，启动服务
```go
package main 

import (
	"fmt"
	"strings"
	"github.com/ruiyanchao/ly"
)

func main(){
    // 初始化worker
	w := ly.Worker{
		SocketName:"tcp://127.0.0.1:8080",
		Protocol: &MyProtocol{
			name:"rick",
		},
	}
   // 定义消息获取处理方法
	w.OnMessage = func(con ly.Connection,msg interface{}) {
		fmt.Println(msg)
		con.Write("hello client")
	}
   // 定义连接开始
	w.OnConnStart = func(con ly.Connection) {
		cid := con.GetConnID()
		fmt.Printf("get client [%d] \n",cid)
	}
   // 定义连接结束
	w.OnConnStop = func(con ly.Connection) {
		cid := con.GetConnID()
		fmt.Printf("client [%d] leave \n",cid)
	}
   // 运行
	w.RunAll()
}
```


## Contributors
- Sam 
- Rick

## TODO
- 添加UDP,UNIX协议
- 自定义协议完成 http websocket
- 压测




