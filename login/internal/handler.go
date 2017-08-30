package internal

import (
	"bearserver/game"
	"bearserver/msg"
	"reflect"

	"github.com/name5566/leaf/gate"
)

func init() {
	handler(&msg.RegisterUserInfo{}, handlRegisterUserInfo)
	handler(&msg.UserLoginInfo{}, handlLoginUser)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handlRegisterUserInfo(args []interface{}) {
	//收到注册信息
	m := args[0].(*msg.RegisterUserInfo)
	//获取发送者
	a := args[1].(gate.Agent)

	//交给 game 模块处理
	game.ChanRPC.Go("RegisterAgent", a, m)
}

func handlLoginUser(args []interface{}) {
	m := args[0].(*msg.UserLoginInfo)
	a := args[1].(gate.Agent)

	//交给 game
	game.ChanRPC.Go("LoginAgent", a, m)
}
