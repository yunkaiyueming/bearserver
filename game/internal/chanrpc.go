package internal

import (
	"bearserver/msg"
	"fmt"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() { //与gate 进行"交流"
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("LoginAgent", rpcLoginAgent)
	skeleton.RegisterChanRPC("RegisterAgent", rpcRigesterAgent)
}

func rpcNewAgent(args []interface{}) {
	fmt.Println("--rpcNew--", args)
	a := args[0].(gate.Agent)
	fmt.Println("args[0]:", a)
	fmt.Println("len():", len(args))
	for i := 0; i < len(args); i++ {
		fmt.Printf("i=%d,arg[%d]=%v \n", i, i, args[i])
	}

	RegNewConn(a, 0)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	LeaveConn(a)
}

func rpcLoginAgent(args []interface{}) {
	fmt.Println("-rpclon-:", args)

	a := args[0].(gate.Agent)
	m := args[1].(*msg.UserLoginInfo)
	uid, err := login(m)
	if err != nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_Login_Error})
		return
	}
	RegNewConn(a, uid)
	a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_Login_OK})
}

func rpcRigesterAgent(args []interface{}) {
	fmt.Println("resiter---")
	a := args[0].(gate.Agent)
	m := args[1].(*msg.RegisterUserInfo)
	ok := checkExitedUser(m.Name)
	log.Debug("hello %v", m.Name)

	if !ok {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_Register_Existed})
		return
	}

	_, err := register(m)
	if err != nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_DB_Error})
		return
	}
}

func rpcJoinRoomAgent(args []interface{}) {

}
