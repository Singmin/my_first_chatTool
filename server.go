package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Port int
	//在线用户的列表
	OnlineMap map[string]*User
	mapLock		sync.RWMutex
	Message		chan string
}

//创建一个server的端口
func NewServer(ip string, port int) *Server {
	server :=&Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

func (this *Server) ListenMessage(){
	for  {
		msg := <-this.Message

		//将msg 发送给全部的在线User
		this.mapLock.Lock()
		for _,cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//广播消息的方法
func(this *Server) BroadCast(user *User,msg string){
	sendMsg := "["+user.Addr+"]"+user.Name + ":"  +msg
	this.Message <- sendMsg
}


func (this *Server) Handler(conn net.Conn){


	user := NewUser(conn,this)
	user.Online()
	//监听用户是否活跃
	isLive := make(chan bool)
	//接受传递消息
	go func() {
		buf := make([]byte,4096)
		for {
			n,err :=conn.Read(buf)
			if n==0 {
				user.Offline()
				return
			}

			if err!= nil && err !=io.EOF{
				fmt.Println("Conn Read err:",err)
				return
			}
			//提取用户消息
			msg :=string(buf[:n-1])


			//广播
			user.DoMessage(msg)
			//用户任意消息代表用户正在活跃
			isLive <- true
		}
	}()
	for{
		select {
			case <-isLive:
				//重置定时器，更新下面的定时器

			case <-time.After(time.Second *100):
				//已经超时
				//将当前的User强制关闭
				user.SendMsg("你已经被踢了")
				close(user.C)
				conn.Close()
				return
		}
	}

}

func(this *Server) Start(){
	//
	listener,err :=net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err != nil{
		fmt.Println("net.Listen err:",err)
		return
	}
	defer  listener.Close()
//启动监听Message
	go this.ListenMessage()


	for {
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err", err)
			continue
		}
		go this.Handler(conn)
	}

}