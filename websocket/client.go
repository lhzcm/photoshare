package websocket

import (
	"errors"
	"photoshare/models"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	//允许向对等体写消息的时间
	writeWait = 10 * time.Second
	//允许从对等端读取下一条pong消息的时间
	pongWait = 60 * time.Second
	//在这段时间内发送ping到对等体，一定比pongWait小
	pingPeriod = (pongWait * 9) / 10
	//消息最大长度
	maxMsgSize = 2048
)

type Client struct {
	hub  *Hub         //集线器
	conn *ws.Conn     //当前客户端ws连接
	user *models.User //连接用户信息
	send chan []byte  //向客户端发送消息channel
}

//向用户发送消息
func (c Client) Send(userid int32, msg []byte) error {
	if client, ok := c.hub.clients[userid]; ok {

		select {
		case client.send <- msg:
		default:
			delete(c.hub.clients, client.user.Id)
			close(client.send)
			return errors.New("客户端没有响应")
		}
	}
	return errors.New("当前用户不在线")
}
