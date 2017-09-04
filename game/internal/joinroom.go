package internal

import (
	//"fmt"

	//"bearserver/msg"
	"fmt"
)

//type RoomInfo struct {
//	RoomID   int
//	RoomName string
//	State    int   //房间状态
//	UserNum  int   //玩家数目
//	UserIds  []int //玩家IDs
//}


func joinRoom(uid int) (bool, string) {
	//取最后一个房间，判断这个房间是否已经开始，如果没有开始，则加入，如果开始，就新建一个房间并且加入
	fmt.Println("room===========")
	fmt.Println(OnlineRooms)
	lastRoom := OnlineRooms[len(OnlineRooms):]
	fmt.Println(lastRoom)


	//for _, item := range OnlineRooms {
	//	if item.RoomName == roomName {
	//		return false, "房间名已经存在"
	//	}
	//}
	//RooCounter++
	//room := msg.Room{
	//	RoomID:   RooCounter,
	//	RoomName: roomName,
	//	State:    0,
	//	UserNum:  1,
	//	UserIds:  []int{uid},
	//	RoomPwd:  roomPwd,
	//}
	//OnlineRooms = append(OnlineRooms, room)
	return false,"11"

}


