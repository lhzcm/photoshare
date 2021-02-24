package websocket

type Result struct {
	Code int         //指令代码, 0:发送消息， 1：接收消息, -1:其他用户登录导致掉线
	Data interface{} //数据
}
