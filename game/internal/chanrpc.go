package internal

import (
	"bearserver/msg"
	"fmt"

	"github.com/name5566/leaf/gate"
	//"github.com/name5566/leaf/log"
)

func init() { //与gate 进行"交流"
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("LoginAgent", rpcLoginAgent)
	skeleton.RegisterChanRPC("RegisterAgent", rpcRigesterAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
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
	fmt.Println("=======")
	fmt.Println(m)
	uid, err := login(m)
	fmt.Println(uid)
	if err != nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE: msg.MSG_Login_Error, Message: err.Error()})
		return
	}

	RegNewConn(a, uid)

	//登录成功之后就开始加入房间
	roomModuel := &RoomModule{}
	room, _ := roomModuel.JoinRoom(uid)
	a.WriteMsg(&msg.Response{Uid: uid, Cmd: "login", Ret: 0, Data: room, Rnum: 1})
}

func rpcRigesterAgent(args []interface{}) {
	fmt.Println("resiter---")
	a := args[0].(gate.Agent)
	m := args[1].(*msg.RegisterUserInfo)
	ok := checkExitedUser(m.Name)
	if ok {
		response := &msg.Response{Cmd: "rigester", Rnum: 1, Ret: -1}
		a.WriteMsg(response)
		return
	}

	_, err := register(m)
	if err == nil {
		response := &msg.Response{Cmd: "rigester", Rnum: 1, Ret: 0}
		a.WriteMsg(response)
		return
	}
}
