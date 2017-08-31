package internal

import (
	"fmt"

	"bearserver/msg"
)

//type RoomInfo struct {
//	RoomID   int
//	RoomName string
//	State    int   //房间状态
//	UserNum  int   //玩家数目
//	UserIds  []int //玩家IDs
//}

var RooCounter int
var OnlineRooms = make([]msg.Room, 0)

func createRoom(uid int, roomName string, roomPwd string) (bool, string) {
	if roomName == "" {
		return false, "房间名不能为空"
	}
	for _, item := range OnlineRooms {
		if item.RoomName == roomName {
			return false, "房间名已经存在"
		}
	}
	RooCounter++
	room := msg.Room{
		RoomID:   RooCounter,
		RoomName: roomName,
		State:    0,
		UserNum:  1,
		UserIds:  []int{uid},
		RoomPwd:  roomPwd,
	}
	OnlineRooms = append(OnlineRooms, room)
}

func getRooms() {

}
