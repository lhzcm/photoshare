package websocket

import "log"

type Hub struct {
	clients    map[int32]*Client //在线客户端列表
	register   chan *Client      //注册客户端
	unregister chan *Client      //注销客户端
}

var MainHub *Hub

//启动集线器
func Start() {
	MainHub = &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int32]*Client),
	}
	for {
		select {
		case client := <-MainHub.register:
			MainHub.clients[client.user.Id] = client
		case client := <-MainHub.unregister:
			if _, ok := MainHub.clients[client.user.Id]; ok {
				log.Printf("客户端又关闭： %#v  %v \n", MainHub.clients, client.user.Id)
				delete(MainHub.clients, client.user.Id)
				close(client.send)
			}
		}
	}
}
