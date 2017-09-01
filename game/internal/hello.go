/*
	测试游戏框架而写的第一个接口
 */
package internal

import (
	"bearserver/msg"
	"fmt"
)

func handleHello(args []interface{})(*msg.Response) {
	m := args[0].(*msg.Dispatch)
	response := &msg.Response{Cmd:m.Cmd,Rnum:m.Rnum,Uid:m.Uid}
	fmt.Println()

	//测试
	type ColorGroup struct {
		ID     int
		Name   string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
	}

	response.Data = group
	return response
}
