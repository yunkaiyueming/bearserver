package internal

import (
	"reflect"

	"bearserver/game"
	"bearserver/msg"

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
	m := args[0].(*msg.RegisterUserInfo)
	a := args[1].(gate.Agent)
	game.ChanRPC.Go("RegisterAgent", a, m)
}

func handlLoginUser(args []interface{}) {
	m := args[0].(*msg.UserLoginInfo)
	a := args[1].(gate.Agent)
	game.ChanRPC.Go("LoginAgent", a, m)
}
