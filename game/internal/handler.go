package internal

import (
	"bearserver/msg"
	"reflect"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	handler(&msg.Hello{}, handleHello)      //具体处理函数调用
	handler(&msg.RoomOperate{}, handleRoom) //具体处理函数调用
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleHello(args []interface{}) {
	m := args[0].(*msg.Hello)
	a := args[1].(gate.Agent)

	log.Debug("hello %v", m.Name)
	a.WriteMsg(&msg.Hello{Name: "ClientHaHa"})
}

func handleRoom(args []interface{}) {
	m := args[0].(*msg.RoomOperate)
	a := args[1].(gate.Agent)

	a.WriteMsg(&msg.Hello{Name: "ClientHaHa"})
}
