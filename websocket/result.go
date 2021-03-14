package websocket

import "photoshare/models"

type Result struct {
	Code int         //指令代码, 0:服务端发送ping 1:发送消息， 2：接收消息, 3：清除未读消息数, -1:服务器请求关闭连接, -2:消息发送异常 -2:其他用户登录导致掉线
	Data interface{} //数据
}

//发送消息
func SendMessage(msg *models.Message) Result {
	return Result{
		Code: 1,
		Data: msg,
	}
}

//消息发送异常
func SendMessageErr(err string) Result {
	return Result{
		Code: -2,
		Data: err,
	}
}

//关闭连接请求
func CloseConnection() Result {
	return Result{
		Code: -1,
		Data: "服务器请求关闭连接",
	}
}

//ping
func SendPing() Result {
	return Result{
		Code: 0,
		Data: "ping",
	}
}
