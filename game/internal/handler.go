package internal

import (
	"reflect"
	"bearserver/msg"

	"github.com/name5566/leaf/gate"
)

func init() {
	handler(&msg.Dispatch{}, handleDispatch)      //处理dispatch
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleDispatch(args []interface{}){
	m := args[0].(*msg.Dispatch)
	a := args[1].(gate.Agent)
	method := m.Cmd
	//这里以后会处理相应的参数要求逻辑


	var response *msg.Response
	switch method {
		case "hello":
			response = handleHello(args)
		default:
			response.Cmd = method
	}

	a.WriteMsg(response)
}
