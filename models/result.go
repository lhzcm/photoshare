package models

const (
	noLogin = -11
	fail    = -1
	success = 0
)

type Result struct {
	Data interface{}
	Msg  string
	Code int
}

func Success(data interface{}, msg string) Result {
	return Result{
		Data: data,
		Msg:  msg,
		Code: success,
	}
}

func Fail(msg string) Result {
	return Result{
		Data: nil,
		Msg:  msg,
		Code: fail,
	}
}

func Nologin() Result {
	return Result{
		Code: noLogin,
	}
}
