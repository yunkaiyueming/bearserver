package internal

import (
	"github.com/name5566/leaf/gate"
	"fmt"
	"bearserver/msg"
	"github.com/name5566/leaf/log"
)

func init() {//与gate 进行"交流"
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("LoginAgent", rpcLoginAgent)
	skeleton.RegisterChanRPC("RegisterAgent",rpcRigesterAgent)
}

func rpcNewAgent(args []interface{}) {
	fmt.Println("--rpcNew--",args)
	a := args[0].(gate.Agent)
	fmt.Println("args[0]:",a)
	fmt.Println("len():",len(args))
	for i := 0; i < len(args); i++{
		//fmt.Fprintln("i=%d,arg[%d]=%v",i,i,args[i])
		fmt.Printf("i=%d,arg[%d]=%v \n",i,i,args[i] )
	}

	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcLoginAgent(args []interface{})  {
	fmt.Println("-rpclon-:",args)
	a := args[0].(gate.Agent)
	fmt.Println("get m--:",a)
	fmt.Println("len--:",len(args))
	m := args[1].(*msg.UserLoginInfo)
	err := login(m)
	if err != nil{
		a.WriteMsg(&msg.CodeState{MSG_STATE:msg.MSG_DB_Error})
		return
	}
}

func rpcRigesterAgent(args []interface{})  {
	fmt.Println("resiter---")
	a := args[0].(gate.Agent)
	m := args[1].(*msg.RegisterUserInfo)
	err := checkExitedUser(m.Name)
	log.Debug("hello %v",m.Name)

	if err == nil {
		a.WriteMsg(&msg.CodeState{MSG_STATE:msg.MSG_Register_Existed})
		return
	}
	err = register(m)
	if err != nil{
		a.WriteMsg(&msg.CodeState{MSG_STATE:msg.MSG_DB_Error})
		return
	}
}

func rpcJoinRoomAgent(args []interface{})  {

}