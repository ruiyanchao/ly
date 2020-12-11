package ly

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

const  (
	DefaultMaxConn = 12000
	DefaultPoolWorkers = 10
	DefaultPoolWorkerDuration = 10 *time.Second
	Version = "1.0.0"
	DefaultMaxPacketSize = 1048576
	DefaultMaxBufferSize = 1048576
	Name = "LY"
	RUNNING = 1
	CLOSE = 2
)

type Worker struct {
	sync.Mutex
	// socket名称
	SocketName string
	// 最大连接数
	MaxConn int
	// 当前连接数，不能超过最大连接数
	currentConn int
	// 自定义协议
	Protocol protocol
	// 协程池 处理消息
	workerPool *workerPool
	// 协程池 进程数量
	PoolWorkerCounts int
	// 协程池 定时清理闲置协程时间
	PoolWorkerDuration time.Duration
	// 闭包 接收消息时处理方法
	OnMessage func(con Connection, msg interface{})
	// 主进程开始时 执行方法
	OnWorkerStart func()
	// 进程名
	Name string
	// 当前主进程关闭通道
	stopChan chan os.Signal
	// 主进程常态
	status int
	// 当前连接属性 如 tcp udp unix
	connection Connection
	// 每个链接 连接时处理方法
	OnConnStart func(connection Connection)
	// 每个链接 断开连接处理方法
	OnConnStop func(connection Connection)
	// 累加 作为连接的唯一ID
	cid int
	// 记录所有的链接
	connections map[int]Connection
	// 包的最大长度
	MaxPacketSize int
	// 包的最大大小
	MaxBufferSize int
	// 信息
	Statistics map[string]interface{}
}

/**
 * 运行
 */
func(w *Worker)RunAll(){
	w.displayUI()
	w.onWorkerStart()
	w.initWorker()
	w.run()
}

/**
 * 展示logo
 */
func(w *Worker)displayUI(){
	logo := `
 _      __   __
| |     \ \ / /
| |      \ V /
| |___    | |
|_____|   |_|`
	fmt.Println(logo)
}

/**
 * 进程开始前运行
 */
func(w *Worker)onWorkerStart(){
	if w.OnWorkerStart != nil {
		w.OnWorkerStart()
	}
}

/**
 * 获取 连接的cid
 */
func(w *Worker)getCID() int{
	w.Lock()
	if w.cid < 1{
		w.cid = 1
	}else{
		w.cid++
	}
	w.Unlock()
	return w.cid

}

/**
 * 新增连接到 connections 属性中
 */
func(w *Worker)addConn(cid int,connection Connection){
	w.Lock()
	w.connections[cid] = connection
	w.Unlock()
}

/**
 * 删除对应连接
 */
func(w *Worker)delConn(cid int){
	w.Lock()
	pro := w.connections[cid].GetProperty()
	w.Statistics["send"] = pro["send"]
	if _,ok:=w.Statistics["connections"];ok{
		w.Statistics["connections"] = w.Statistics["connections"].(int) + 1
	}else{
		w.Statistics["connections"] = 1
	}
	delete(w.connections, cid)
	w.Unlock()
}

/**
 * 初始化进程
 */
func(w *Worker)initWorker(){
	// 初始化进程名
	if w.Name == ""{
		w.Name = Name
	}
	// 初始化buff大小
	if w.MaxBufferSize == 0{
		w.MaxBufferSize = DefaultMaxBufferSize
	}
	// 初始化buff长度
	if w.MaxPacketSize == 0{
		w.MaxPacketSize = DefaultMaxPacketSize
	}
	// 初始化 协程池 进程数量
	if w.PoolWorkerCounts < 1{
		w.PoolWorkerCounts = DefaultPoolWorkers
	}
	// 初始化 协程池协程超时时间
	if w.PoolWorkerDuration < 1 * time.Second{
		w.PoolWorkerDuration = DefaultPoolWorkerDuration
	}
	// 实例化协程池
	w.workerPool = &workerPool{
		WorkerFunc: w.OnMessage,
		MaxWorkersCount:       w.PoolWorkerCounts,
		Logger:                Logger(log.New(os.Stderr, "", log.LstdFlags)),
		MaxIdleWorkerDuration: w.PoolWorkerDuration,
	}
	// 初始化当前连接数
	w.currentConn = 0
	// 开启协程池
	w.workerPool.Start()
	// 初始化最大连接
	if w.MaxConn < 1{
		w.MaxConn = DefaultMaxConn
	}
	// 初始化关闭通道
	w.stopChan = make(chan os.Signal)
	// 初始化运行状态
	w.status = RUNNING
	// 初始化监控
	w.Statistics = make(map[string]interface{})
	// 初始化cid
	w.cid = 0
	// 初始化 连接管理
	w.connections = make(map[int]Connection)
	//获取协议名
	protocolName := "nil"
	if w.Protocol != nil{
		protocolName = w.Protocol.GetName()
	}
	fmt.Printf("[LY] Version: %s, MaxConn: %d, MaxPacketSize: %d, Protocol: %s \n",Version,w.MaxConn,w.MaxPacketSize,protocolName)
}

/**
 * 运行监听
 */
func(w *Worker)run(){
	//获取对应 协议 名
	if !strings.Contains(w.SocketName,"://"){
		fmt.Println("err socket name")
	}
	snSplit := strings.Split(w.SocketName, "://")
	pro := snSplit[0]
	fmt.Printf("[START] Server name: %s,socket name : %s is starting \n",w.Name,w.SocketName)
	//TODO 添加 UDP UNIX
	switch pro{
	case "tcp":
		// 监听tcp
		w.listenTcp(snSplit[1])
	default:
		fmt.Println("err type")
	}
}

/**
 * 监听tcp
 */
func(w *Worker)listenTcp(address string){

	// 开启服务
	server, err := net.Listen("tcp",address)
	// 监听终端关闭信息
	w.monitorWorker(server)
	if err != nil{
		fmt.Println("server err:",err)
		return
	}
	fmt.Printf("start LY server  %s  succ, now listenning... \n" ,w.Name)
	for {
		// 接收链接
		conn, err := server.Accept()
		if err != nil {
			// 如果进程状态为关闭 跳出循环
			if w.status == CLOSE{
				break
			}
			fmt.Println("accept err")
			continue
		}
		// 链接超出 将关闭
		if w.currentConn >= w.MaxConn{
			conn.Close()
			fmt.Println("connection too much")
			continue
		}
		// 接收信息
		go w.connHandler(conn)
	}
}

/**
 * 处理链接
 */
func(w *Worker)connHandler(conn net.Conn){
	//初始化tcp链接
	ac := w.acceptTcpConnection(conn)
	//添加链接
	w.addConn(ac.GetConnID(),ac)
	//连接触发方法
	ac.OnConnStart()
	defer ac.Close()
	for {
		//读取消息
		msg,err := ac.Read()
		if err != nil{
			break
		}
		//将消息投入协程池处理
		for _,v :=range msg{
			w.workerPool.Serve(ac,v)
		}

	}
	w.delConn(ac.GetConnID())
	ac.OnConnStop()
}
/**
 * 初始化top链接
 */
func(w *Worker)acceptTcpConnection(conn net.Conn)(connection Connection){

	connection = &tcp{
		conn: conn,
		MaxBuffer: w.MaxBufferSize,
		MaxPageSize: w.MaxPacketSize,
		name:"tcp-server",
		onConnStart: w.OnConnStart,
		onConnStop: w.OnConnStop,
		protocol: w.Protocol,
		connID:w.getCID(),
		property:make(map[string]interface{}),
	}
	w.connection = connection
	return
}

/**
 * 监听信息
 */
func (w *Worker)monitorWorker(server net.Listener){
	signal.Notify(w.stopChan, os.Interrupt)
	go func() {
		<-w.stopChan
		w.status = CLOSE
		fmt.Print("\n")
		fmt.Println("get stop command. now stopping...")
		fmt.Printf("[STOP] Server name: %s,socket name : %s is stopped \n",w.Name,w.SocketName)
		send:=0
		connections := len(w.connections)
		for _,v:=range w.connections{
			p:=v.GetProperty()
			if _,ok := p["send"];ok{
				send =send + p["send"].(int)
			}
		}
		if _,ok:= w.Statistics["send"];ok{
			send  = send + w.Statistics["send"].(int)
		}
		if _,ok:= w.Statistics["connections"];ok{
			connections  = connections + w.Statistics["connections"].(int)
		}
		fmt.Printf("[LY] Totol connection nums : %d,Totol send : %d  \n",connections,send)
		if err := server.Close(); err != nil {
			fmt.Println("close listen err",err)
		}

	}()
}











