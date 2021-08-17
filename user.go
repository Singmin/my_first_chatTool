package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C  chan string
	conn net.Conn
	server *Server
}

func  NewUser(conn net.Conn,server *Server)  *User{
	userAddr := conn.RemoteAddr().String()

	user :=&User{
		Name:  userAddr,
		Addr:  userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}
//启动监听的GO层
	go user.ListenMessage()

	return user
}

func (this *User)  Online(){
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this,"已上线")
}

func (this *User)  Offline(){
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this,"Quit")
}

func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

func (this *User)  DoMessage(msg string){
	if msg == "who"{
		//查询当前在线用户

		this.server.mapLock.Lock()
		for _, user :=range this.server.OnlineMap{
			onlineMsg :="[" +user.Addr +"]"+user.Name +":"+"Online\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	}else if len(msg)>7 && msg[:7]=="rename|"{
		newName := strings.Split(msg,"|")[1]

		//判断Name是否已经存在
		_,ok :=this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名被使用\n")
		}else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName]=this
			this.server.mapLock.Unlock()

			this.Name=newName
			this.SendMsg("您已经更新用户名"+newName+"\n")
		}
	}else if len(msg)>4&&msg[0:3]=="to|"{
		//to|张三|消息内容
		remoteName := strings.Split(msg,"|")[1]
		if remoteName ==""{
			this.SendMsg("消息格式不对")
			return
		}

		remoteUser, ok :=this.server.OnlineMap[remoteName]
		if !ok{
			this.SendMsg("该用户名不存在\n")
			return
		}

		//获取消息内容，通过对方的User对象发送内容过去
		content := strings.Split(msg,"|")[2]
		if content == ""{
			this.SendMsg("无效内容，请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name+"对您说:"+content+"\n")
	}else{
		this.server.BroadCast(this,msg)
	}
}

//


func (this *User) ListenMessage() {
	for{
		msg := <- this.C

		this.conn.Write([]byte(msg +"\n"))
	}
}



