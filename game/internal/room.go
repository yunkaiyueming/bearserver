package internal

import (
	"errors"
	"fmt"

	"bearserver/msg"

	"github.com/name5566/leaf/gate"
)

type Room struct {
	RoomID    int                 //房间ID
	RoomName  string              //房间名
	State     int                 //房间状态: 0问开局， 1开局中， 2已结束
	UserNum   int                 //玩家数目
	UserIds   []int               //玩家IDs
	RoomPwd   string              //房间密码
	UserState map[int]PlayerState //玩家信息
	RecvCh    chan map[string]interface{}
}

type RoomModule struct{}

type PlayerState struct {
	Uid    int
	Cards  []Card
	Status int //玩家状态：0进入房间， 1开始， 2有出牌权，3没有出牌权，4赢了，5输了
}

const (
	OneRoomPlayerNum = 4

	READY    = 0
	PLAYING  = 1
	GAMEOVER = 2
)

var RoomCounter int
var OnlineRooms = make([]Room, 0)

func (r *RoomModule) createRoom(uid int) (Room, error) {
	RoomCounter++

	userState := map[int]PlayerState{uid: {Uid: uid, Status: 0}}
	room := Room{
		RoomID:    RoomCounter,
		RoomName:  fmt.Sprintf("房间%d", RoomCounter),
		State:     READY,
		UserNum:   1,
		UserIds:   []int{uid},
		UserState: userState,
	}
	OnlineRooms = append(OnlineRooms, room)
	return room, nil
}

func (r *RoomModule) GetRoom(roomId int) (Room, error) {
	for _, room := range OnlineRooms {
		if room.RoomID == roomId {
			return room, nil
		}
	}
	return Room{}, errors.New("无此房间信息")
}

func (r *RoomModule) IsNeedCreateRoom() bool {
	lastRoom := OnlineRooms[len(OnlineRooms)-1]
	if lastRoom.UserNum < OneRoomPlayerNum {
		return false
	}
	return true
}

func (r *RoomModule) JoinRoom(uid int) (Room, error) {
	if r.IsNeedCreateRoom() {
		return r.createRoom(uid)
	}
	lastRoom := OnlineRooms[len(OnlineRooms)-1]
	lastRoom.UserIds = append(lastRoom.UserIds, uid)
	lastRoom.UserNum++
	if lastRoom.UserState == nil {
		lastRoom.UserState = make(map[int]PlayerState)
	}
	lastRoom.UserState[uid] = PlayerState{Uid: uid, Status: 0}

	return lastRoom, nil
}

func (r *RoomModule) Start(args []interface{}) {
	m := args[0].(*msg.Dispatch)
	a := args[1].(gate.Agent)

	playerModuel := PlayerModuel{}
	params := m.Params.(map[string]interface{})
	roomId := params["room_id"].(int)
	if m.Cmd == "room.start" {
		//检查房间
		room, err := r.GetRoom(roomId)
		if err != nil {
			a.WriteMsg(err.Error())
			return
		}

		//检查人数
		if room.UserNum != OneRoomPlayerNum {
			a.WriteMsg(fmt.Sprintf("房间人数不足，不能开局"))
			return
		}

		//开局
		go playerModuel.start(&room, a)
	}

	go playerModuel.RecvRoomMsg(a, roomId, params)
}
