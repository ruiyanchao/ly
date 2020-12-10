package LY

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
	SocketName string
	MaxConn int
	currentConn int
	Protocol protocol
	workerPool *workerPool
	PoolWorkerCounts int
	PoolWorkerDuration time.Duration
	OnMessage func(con Connection, msg interface{})
	OnWorkerStart func()
	Name string
	stopChan chan os.Signal
	status int
	connection Connection
	OnConnStart func(connection Connection)
	OnConnStop func(connection Connection)
	cid int
	connections map[int]Connection
}

func(w *Worker)RunAll(){
	w.displayUI()
	w.onWorkerStart()
	w.initWorker()
	w.run()
}


func(w *Worker)displayUI(){
	logo := `
 _      __   __
| |     \ \ / /
| |      \ V /
| |___    | |
|_____|   |_|`
	fmt.Println(logo)
}

func(w *Worker)onWorkerStart(){
	if w.OnWorkerStart != nil {
		w.OnWorkerStart()
	}
}

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

func(w *Worker)addConn(cid int,connection Connection){
	w.Lock()
	w.connections[cid] = connection
	w.Unlock()
}


func(w *Worker)initWorker(){
	if w.Name == ""{
		w.Name = Name
	}
	if w.PoolWorkerCounts < 1{
		w.PoolWorkerCounts = DefaultPoolWorkers
	}
	if w.PoolWorkerDuration < 1 * time.Second{
		w.PoolWorkerDuration = DefaultPoolWorkerDuration
	}
	w.workerPool = &workerPool{
		WorkerFunc: w.OnMessage,
		MaxWorkersCount:       w.PoolWorkerCounts,
		Logger:                Logger(log.New(os.Stderr, "", log.LstdFlags)),
		MaxIdleWorkerDuration: w.PoolWorkerDuration,
	}
	w.currentConn = 0
	w.workerPool.Start()
	if w.MaxConn < 1{
		w.MaxConn = DefaultMaxConn
	}
	w.stopChan = make(chan os.Signal)
	w.status = RUNNING
	w.cid = 0
	w.connections = make(map[int]Connection)
	protocolName := "nil"
	if w.Protocol != nil{
		protocolName = w.Protocol.GetName()
	}
	fmt.Printf("[LY] Version: %s, MaxConn: %d, MaxPacketSize: %d, Protocol: %s \n",Version,w.MaxConn,DefaultMaxPacketSize,protocolName)
}

func(w *Worker)run(){
	if !strings.Contains(w.SocketName,"://"){
		fmt.Println("err socket name")
	}
	snSplit := strings.Split(w.SocketName, "://")
	pro := snSplit[0]
	fmt.Printf("[START] Server name: %s,socket name : %s is starting \n",w.Name,w.SocketName)
	switch pro{
	case "tcp":
		w.listenTcp(snSplit[1])
	default:
		fmt.Println("err type")
	}
}

func(w *Worker)listenTcp(address string){

	server, err := net.Listen("tcp",address)
	w.monitorWorker(server)
	if err != nil{
		fmt.Println("server err:",err)
		return
	}
	fmt.Printf("start LY server  %s  succ, now listenning... \n" ,w.Name)
	for {
		conn, err := server.Accept()
		if err != nil {
			if w.status == CLOSE{
				break
			}
			fmt.Println("accept err")
			continue
		}
		if w.currentConn >= w.MaxConn{
			conn.Close()
			fmt.Println("connection too much")
			continue
		}
		go w.connHandler(conn)
	}
}


func(w *Worker)connHandler(conn net.Conn){
	ac := w.acceptTcpConnection(conn)
	w.addConn(ac.GetConnID(),ac)
	ac.OnConnStart()
	defer ac.Close()
	for {
		msg,err := ac.Read()
		if err != nil{
			break
		}
		for _,v :=range msg{
			w.workerPool.Serve(ac,v)
		}

	}
	ac.OnConnStop()
}

func(w *Worker)acceptTcpConnection(conn net.Conn)(connection Connection){

	connection = &tcp{
		conn: conn,
		defaultMaxBuffer: DefaultMaxBufferSize,
		defaultMaxPageSize: DefaultMaxPacketSize,
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

func (w *Worker)monitorWorker(server net.Listener){
	signal.Notify(w.stopChan, os.Interrupt)
	go func() {
		<-w.stopChan
		w.status = CLOSE
		fmt.Print("\n")
		fmt.Println("get stop command. now stopping...")
		fmt.Printf("[STOP] Server name: %s,socket name : %s is stopped \n",w.Name,w.SocketName)
		send:=0
		for _,v:=range w.connections{
			p:=v.GetProperty()
			if _,ok := p["send"];ok{
				send =send + p["send"].(int)
			}
		}
		fmt.Printf("[LY] Connection nums : %d,Totol send : %d  \n",len(w.connections),send)
		if err := server.Close(); err != nil {
			fmt.Println("close listen err",err)
		}

	}()
}











