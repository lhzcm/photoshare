package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"photoshare/models"
	"photoshare/redis"
	"photoshare/service"
	"strings"
	"time"

	"github.com/goinggo/mapstructure"
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

var (
	newline  = []byte{'\n'}
	upgrader = ws.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Client struct {
	hub  *Hub         //集线器
	conn *ws.Conn     //当前客户端ws连接
	user *models.User //连接用户信息
	send chan []byte  //向客户端发送消息channel
}

//向用户发送消息
func (c *Client) Send(userid int32, msg Result) error {
	if client, ok := c.hub.clients[userid]; ok {
		msgbyte, _ := json.Marshal(msg)
		select {
		case client.send <- msgbyte:
		default:
			delete(c.hub.clients, client.user.Id)
			close(client.send)
			return errors.New("客户端没有响应")
		}
		return nil
	} else {
		return errors.New("当前用户不在线")
	}
}

//向客户端发送消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteJSON(CloseConnection())
				log.Println("send channel has close")
				return
			}

			if err := c.conn.WriteMessage(ws.TextMessage, message); err != nil {
				log.Println("WriteMessage error")
				log.Println(err.Error())
				return
			}
			for len(c.send) > 0 {
				if err := c.conn.WriteMessage(ws.TextMessage, <-c.send); err != nil {
					log.Println("WriteMessage error")
					log.Println(err.Error())
					return
				}
			}

			// w, err := c.conn.NextWriter(ws.TextMessage)
			// if err != nil {
			// 	log.Println("nextwriter create fail")
			// 	log.Panicln(err.Error())
			// 	return
			// }

			// w.Write(message)
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			// //关闭消息写入
			// if err := w.Close(); err != nil {
			// 	log.Printf("close error")
			// 	log.Println(err)
			// 	return
			// }
		case <-ticker.C:
			//c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

//向客户端读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	//c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(msg string) error {
	//	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// 	return nil
	// })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msgs := strings.Split(string(message), string(newline))
		for _, item := range msgs {
			c.MessageDeal(item)
		}
	}
}

//对消息进行处理
func (c *Client) MessageDeal(msg string) {
	var result Result
	if err := json.Unmarshal([]byte(msg), &result); err != nil {
		c.Send(c.user.Id, SendMessageErr("消息解析错误"))
		return
	}

	//接受到发送消息
	if result.Code == 2 {
		var msg models.Message
		data, ok := result.Data.(map[string]interface{})
		if ok {
			if err := mapstructure.Decode(data, &msg); err != nil {
				c.Send(msg.Senderid, SendMessageErr("消息解析错误"))
				return
			}
		}
		msg.Senderid = c.user.Id

		if err := service.SendMessage(&msg); err != nil {
			c.Send(msg.Senderid, SendMessageErr(err.Error()))
			return
		}

		//发送给在线用户
		msgResult := SendMessage(&msg)
		if err := c.Send(msg.Receiverid, msgResult); err != nil {
			//用户不在线缓存到redis
			bytemsg, _ := json.Marshal(msgResult)
			redis.RedisAddMsg(msg.Receiverid, bytemsg)
		}
	}
	//else TODO

}

//启用websocket客户端
func StartClient(w http.ResponseWriter, r *http.Request, user *models.User) {
	log.Println(user)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: MainHub, conn: conn, send: make(chan []byte, 256), user: user}
	MainHub.register <- client

	go client.writePump()
	go client.readPump()

	//读取redis缓存的消息
	for msg, err := redis.RedisGetMsg(client.user.Id); err == nil && msg != nil; {
		client.send <- msg
	}
}
