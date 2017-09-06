package internal

import (
	"errors"
	"fmt"

	//"bearserver/msg"

	//"github.com/name5566/leaf/gate"
	"bearserver/gamedata/db"
	"github.com/name5566/leaf/log"
)

type Room struct {
	RoomID    int                 //房间ID
	RoomName  string              //房间名
	State     int                 //房间状态: 0问开局， 1开局中， 2已结束
	UserNum   int                 //玩家数目
	UserIds   []int               //玩家IDs
	RoomPwd   string              //房间密码
	UserState map[int]PlayerState //玩家信息
	Center    Card 				//房间中间那张牌
	RecvCh    chan map[string]interface{}
}

type RoomModule struct{}

type PlayerState struct {
	Uid int
	Cards []Card
	Name string
	Status int //玩家状态：0进入房间， 1开始， 2有出牌权，3没有出牌权，4赢了，5输了
}

const (
	OneRoomPlayerNum = 3

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
	if len(OnlineRooms) == 0{
		return true
	}else {
		lastRoom := OnlineRooms[len(OnlineRooms)-1]
		if lastRoom.UserNum < OneRoomPlayerNum{
			return false
		}
	}
	return true
}

func (r *RoomModule) JoinRoom(uid int) (interface{}, error) {
	if r.IsNeedCreateRoom() {
		r.createRoom(uid)
	}
	lastRoom := OnlineRooms[len(OnlineRooms)-1]
	lastRoom.UserIds = append(lastRoom.UserIds, uid)
	lastRoom.UserNum++
	if lastRoom.UserState == nil {
		lastRoom.UserState = make(map[int]PlayerState)
	}

	userModel := db.ModelUser{}
	uInfo,_ := userModel.GetUserById(uid)

	lastRoom.UserState[uid] = PlayerState{Uid: uid, Status: 0,Name:uInfo.Name}

	//检查是否可以开始游戏

	r.Start(&lastRoom)
	roomInfo,_ := r.getRoomInfo(uid,lastRoom.RoomID)
	//给房间里面的其他人推送信息
	for k,_ := range lastRoom.UserState{
		if k != uid {
			perroomInfo,_ := r.getRoomInfo(k,lastRoom.RoomID)
			PushMsgModuel := PushMsgModuel{}
			log.Debug("push...",k,perroomInfo)
			PushMsgModuel.pushMsgByUid(k,perroomInfo)
		}
	}
	return roomInfo, nil
}

//给client返回房间信息
func (r *RoomModule) getRoomInfo(uid int,roomId int) (interface{}, error) {
	type Players struct {
		P1 struct{
			Uid int
			Name string
			CardNum int
		}
		P2 struct{
			Uid int
			Name string
			CardNum int
		}
		P3 struct{
			Uid int
			Name string
			CardNum int
		}
	}

	type Roominfo struct{
		Player Players
		Center Card
		MyCards []Card
		MathcFlag bool
		Turn int
		TurnTime int
		MyPos int
	}
	//
	room,_ := r.GetRoom(roomId)
	players := Players{}
	//
	index := 1
	myPos := 1
	resRoom := Roominfo{}
	for _,v := range room.UserState{
		if index == 1{
			players.P1.Uid = v.Uid
			players.P1.Name = v.Name
			players.P1.CardNum = len(v.Cards)
		}else if index == 2{
			players.P2.Uid = v.Uid
			players.P2.Name = v.Name
			players.P2.CardNum = len(v.Cards)
		}else if index == 3{
			players.P3.Uid = v.Uid
			players.P3.Name = v.Name
			players.P3.CardNum = len(v.Cards)
		}
		if v.Uid == uid{
			myPos = index
			resRoom.MyCards = v.Cards
		}
		index++
	}

	if room.State == PLAYING{
		resRoom.Center = room.Center
		resRoom.MathcFlag = true
	}else {
		resRoom.MathcFlag = false
	}
	resRoom.Player = players
	resRoom.MyPos = myPos

	return  resRoom,nil



}

func (r *RoomModule) Start(room *Room) {
	playerModuel := PlayerModuel{}
	log.Debug("start...111")
	//params := m.Params.(map[string]interface{})
	//检查房间
	//room, err := r.GetRoom(roomId)
	//if err != nil {
	//	return
	//}
	//检查人数
	log.Debug("UserNum...",room.UserNum )
	if room.UserNum != OneRoomPlayerNum {
		return
	}
	log.Debug("start...222")
	//开局
	go playerModuel.start(room)





	//go playerModuel.RecvRoomMsg(a, roomId, params)
}
