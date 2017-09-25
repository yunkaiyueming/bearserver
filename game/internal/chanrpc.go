package internal

import (
	"fmt"

	"bearserver/msg"

	"github.com/name5566/leaf/gate"
)

func init() {
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
	fmt.Println("login start---")
	a := args[0].(gate.Agent)
	m := args[1].(*msg.UserLoginInfo)

	uid, err := (&UserModule{}).Login(m)
	if err != nil {
		a.WriteMsg(&msg.Response{Cmd: "login", Ret: msg.MSG_Login_Error, Data: err.Error(), Rnum: 1})
		return
	}

	RegNewConn(a, uid)

	//登录成功之后就开始加入房间
	roomModuel := &RoomModule{}
	room, _ := roomModuel.JoinRoom(uid)
	a.WriteMsg(&msg.Response{Uid: uid, Cmd: "login", Ret: 0, Data: room, Rnum: 1})
}

func rpcRigesterAgent(args []interface{}) {
	fmt.Println("resiter start---")

	a := args[0].(gate.Agent)
	m := args[1].(*msg.RegisterUserInfo)
	ok := (&UserModule{}).CheckExitedUser(m.Name)
	response := &msg.Response{}
	if ok {
		response = &msg.Response{Cmd: "register", Rnum: 1, Ret: -1, Data: "uname is exist"}
		a.WriteMsg(response)
		return
	}

	_, err := (&UserModule{}).Register(m)
	if err != nil {
		response = &msg.Response{Cmd: "rigester", Rnum: 1, Ret: -1, Data: err.Error()}
	} else {
		response = &msg.Response{Cmd: "rigester", Rnum: 1, Ret: 0}
	}

	a.WriteMsg(response)
}
