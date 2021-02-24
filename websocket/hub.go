package websocket

type Hub struct {
	clients    map[int32]*Client //在线客户端列表
	register   chan *Client      //注册客户端
	unregister chan *Client      //注销客户端
}

//启动集线器
func (hub Hub) Start() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client.user.Id] = client
		case client := <-hub.unregister:
			if _, ok := hub.clients[client.user.Id]; ok {
				delete(hub.clients, client.user.Id)
				close(client.send)
			}
		}
	}
}
