package gate

import (
	"bearserver/game"
	"bearserver/login"
	"bearserver/msg"
)

//消息在此进行交割
func init() {
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC) //参数消息内容 通信桥chanRPC
	msg.Processor.SetRouter(&msg.JoinRoomInfo{}, game.ChanRPC)

	//用注册
	msg.Processor.SetRouter(&msg.RegisterUserInfo{}, login.ChanRPC)
	//登录
	msg.Processor.SetRouter(&msg.UserLoginInfo{}, login.ChanRPC)
}
