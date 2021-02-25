package websocket

type Result struct {
	Code int         //指令代码, 0:服务端发送ping 1:发送消息， 2：接收消息, -1:服务器请求关闭连接 -2:其他用户登录导致掉线
	Data interface{} //数据
}

func CloseConnection() Result {
	return Result{
		Code: -1,
		Data: "服务器请求关闭连接",
	}
}

func SendPing() Result {
	return Result{
		Code: 0,
		Data: "ping",
	}
}
