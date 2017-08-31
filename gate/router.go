package gate

import (
	"bearserver/game"
	"bearserver/login"
	"bearserver/msg"
)

//消息在此进行交割
func init() {
	//处理游戏逻辑
	msg.Processor.SetRouter(&msg.Dispatch{}, game.ChanRPC)
	//用注册
	msg.Processor.SetRouter(&msg.RegisterUserInfo{}, login.ChanRPC)
	//登录
	msg.Processor.SetRouter(&msg.UserLoginInfo{}, login.ChanRPC)
}
