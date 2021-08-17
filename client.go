package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

func NewClient(serverIp string,serverPort int) *Client{
	//创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag:999,
	}
	//连接server
	conn,err :=net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err !=nil{
		fmt.Println("net.Dial:",err)
		return nil
	}

	client.conn =conn
	return client
}

func(client *Client) DealResponse(){
	io.Copy(os.Stdout,client.conn)
}

func (Client *Client) menu() bool {
	var flag int
	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)

	if flag >=0 && flag<=3{
		Client.flag = flag
		return true
	}else {
		fmt.Println(">>>>>>>>>>请输入合法范围内的数字>>>>>>>>>>>")
		return false
	}
}

func (client *Client) SelectUsers() {
	sendMsg :="who\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err !=nil{
		fmt.Println("conn Write err :",err)
		return
	}
}

func (client *Client) PrivateChat(){
	var remoteName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println(">>>>>>>请输入聊天对象[用户名]>>>>输入exit退出>>>>>>")
	fmt.Scanln(&remoteName)
	for remoteName != "exit"{
		fmt.Println(">>>>>>>>>>请输入聊天内容>>>>>输入exit退出>>>>")
		fmt.Scanln(&chatMsg)
		for chatMsg !="exit"{
			if len(chatMsg) !=0{
				sendMsg := "to|"+remoteName+"|"+chatMsg +"\n"
				_,err :=client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:",err)
					break
				}
			}
			chatMsg=""
			fmt.Println(">>>>>>>>>>请输入聊天内容>>>>输入exit退出>>>>>")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>>>>请输入聊天对象[用户名]>>>>输入exit退出>>>>>>")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat(){
	var chatMsg string
	fmt.Println(">>>>>>>>请输入聊天内容，exit退出>>>>>>>>>>")
	fmt.Scanln(&chatMsg)
	for chatMsg !="exit"{
		if len(chatMsg) != 0 {
			sendMsg := chatMsg +"\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {}
				fmt.Println("conn Write err:",err)
			break
		}
	}
	chatMsg = " "
	fmt.Println(">>>>>>>>请输入聊天内容，exit退出>>>>>>>>>>")
	fmt.Scanln(&chatMsg)
}

func (client *Client)  UpdateName() bool{
	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|"+ client.Name +"\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err !=nil {
		fmt.Println("conn.Write err:",err)
		return false
	}
	return true
}


func (client *Client) Run(){
	for client.flag !=0{
		for client.menu() !=true{
		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式选择....")
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式选择....")
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init(){
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器IP地址是（默认是127.0.0.1）")
	flag.IntVar(&serverPort,"port",8888,"设置服务器端口(默认是8888)")

}

func main(){
	flag.Parse()

	client :=NewClient(serverIp,serverPort)
	if client ==nil {
		fmt.Println(">>>>>连接失败")
		return
	}

	go client.DealResponse()
	fmt.Println(">>>>>>连接服务器成功>>>>>>")

	client.Run()
}
