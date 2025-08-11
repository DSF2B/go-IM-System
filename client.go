package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 0:公聊 1:私聊
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}
func (client *Client) DealResponse() {
	//一旦client.conn有数据，就读取并打印出来，永久阻塞监听
	io.Copy(os.Stdout, client.conn)

}
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.私聊")
	fmt.Println("2.公聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("输入错误")
		return false
	}
}

var serverIp string
var serverPort int

func init() {
	//./client -ip 127.0.0.1 -port 7777
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip:127.0.0.1 default")
	flag.IntVar(&serverPort, "port", 7777, "set server port:7777 default")
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("私聊模式")
			client.PrivateChat()
			break
		case 2:
			fmt.Println("公聊模式")
			client.PublicChat()
			break
		case 3:
			fmt.Println("更新用户名")
			client.UpdataName()
			break
		case 4:
			fmt.Println("退出")
			break

		}
	}
}
func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("send who msg error:", err)
		return
	}
}
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUser()
	fmt.Println("请输入要私聊的用户名:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入要发送的消息:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("send private chat msg error:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入要发送的消息:")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUser()
		fmt.Println("请输入要私聊的用户名:")
		fmt.Scanln(&remoteName)
	}
}
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入要发送的消息:")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("send public chat msg error:", err)
				break
			}

		}

		chatMsg = ""
		fmt.Println("请输入要发送的消息:")
		fmt.Scanln(&chatMsg)
	}
}
func (client *Client) UpdataName() bool {
	fmt.Println("请输入新的用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("send rename msg error:", err)
		return false
	}
	return true
}
func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("connect to server failed")
		return
	}
	go client.DealResponse()
	fmt.Println("connect to server success")
	client.Run()
}
