package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

// server接口
func NewServer(ip string, port int) *Server {
	server := Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return &server
}
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.RLock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

// 处理客户端连接
func (this *Server) Handler(conn net.Conn) {
	// fmt.Println("new client connected")
	//用户上线
	user := NewUser(conn, this)
	user.Online()
	//监听是否活跃的channel
	isLive := make(chan bool)
	//接受客户端消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err", err)
				return
			}
			msg := string(buf[:n-1])
			user.Domessage(msg)
			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 300):
			//超时，关闭连接
			user.SendMessage("你被踢了")
			close(user.C)
			conn.Close()
			return
		}
	}
}

// 启动服务器
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	defer listener.Close()
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err", err)
			continue
		}

		go this.Handler(conn)
		//do handler
	}

	//close listen socket

}
