package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Dispatch{})
	Processor.Register(&Response{})
	Processor.Register(&UserLoginInfo{})
	Processor.Register(&RegisterUserInfo{})
}

//客户端发送请求格式
type Dispatch struct {
	Uid    int                    `json:"uid"`
	Cmd    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
	Rnum   int                    `json:"rnum"`
	Ts     int                    `json:"ts"`
}

//服务端返回数据格式
type Response struct {
	Uid  int         `json:"uid"`
	Cmd  string      `json:"cmd"`
	Ret  int         `json:"ret"`
	Data interface{} `json:"data"`
	Rnum int         `json:"rnum"`
}

type UserLoginInfo struct { //登录
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

type RegisterUserInfo struct { //注册
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}
