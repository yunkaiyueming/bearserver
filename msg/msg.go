package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&UserLoginInfo{})
	Processor.Register(&LoginError{})

	Processor.Register(&RegisterUserInfo{})

	Processor.Register(&CodeState{})

	//房间会话注册
	Processor.Register(&Room{})        //基本信息
	Processor.Register(&RoomOperate{}) //用户输入密码 点击进入
}

type CodeState struct {
	MSG_STATE int    // const
	Message   string //警告信息
}

type Hello struct {
	Name string
}

type UserLoginInfo struct { //登录
	Name string
	Pwd  string
}

type LoginError struct {
	State   int
	Message string
}

type RegisterUserInfo struct { //注册
	Name string
	Pwd  string
}

type Room struct {
	RoomID   int
	RoomName string
	State    int    //房间状态
	UserNum  int    //玩家数目
	UserIds  []int  //玩家IDs
	RoomPwd  string //房间密码
}

type RoomOperate struct {
	RoomNumber string
	RoomPwd    string
	Type       int
}

type RoomPWDJoinCondition struct {
	Pwd string
}
