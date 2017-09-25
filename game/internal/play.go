package internal

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"bearserver/conf"
	"bearserver/msg"

	"github.com/name5566/leaf/log"
)

type PlayerModuel struct{}

var MaxBaseCardNum = len(conf.BaseCards)

//处理客户端请求逻辑
func (p *PlayerModuel) HandlePlayCard(args []interface{}) *msg.Response {
	m := args[0].(*msg.Dispatch)
	response := &msg.Response{Cmd: m.Cmd, Rnum: m.Rnum, Uid: m.Uid, Ret: 0}

	if err := p.checkRoomMsg(m); err != nil {
		response.Ret = -1
		return response
	}

	uid := m.Uid
	params := m.Params
	mtype, ok := params["mtype"]
	if !ok {
		response.Ret = -3
		return response
	}

	//mtype:1发牌 2摸牌
	if mtype == "1" || mtype == "2" {
		if mtype == "1" {
			card, ok := params["card"]
			log.Debug("cards...", card)
			if !ok {
				response.Ret = -3
				return response
			}
		}
	} else {
		response.Ret = -3
		return response
	}
	log.Debug("step22")

	var card int
	log.Debug("card...", int(card))

	switch v := params["card"].(type) {
	case int:
		fmt.Println("整型", v)

	case string:
		fmt.Println("字符串", v)
		cardInt, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		card = cardInt
	}

	playerModuel := PlayerModuel{}
	roomModuel := RoomModule{}
	room, ok := roomModuel.getRoomByUid(uid)
	if !ok {
		response.Ret = -5
		return response
	}

	log.Debug("step3...")
	log.Debug("card...", card)

	if mtype == "1" {
		//发牌
		playRet := playerModuel.discard(uid, card, &room)
		if !playRet {
			response.Ret = -6
			return response
		}
	} else if mtype == "2" {
		//摸牌
		log.Debug("mtype2...", mtype)

		playRet := playerModuel.draw(uid, &room)
		if !playRet {
			response.Ret = -5
			return response
		}
	}
	//这个房间收到消息，不会自动发牌
	roomModuel.RecvRoomMsg(&room, params)
	perroomInfo, _ := roomModuel.getRoomInfo(uid, room.RoomID)
	roomModuel.pushRoomMsgToOthers(uid, &room)

	response.Data = perroomInfo

	return response
}

//洗牌
func (p *PlayerModuel) initPlayerCards() []int {
	roomCards := conf.BaseCards

	sortCards := make([]int, 0)
	for cardId, _ := range roomCards {
		sortCards = append(sortCards, int(cardId))
	}
	return sortCards
}

//第一次给玩家发牌
func (p *PlayerModuel) sendCard(room *Room, sortCards []int) {
	room.Cards = sortCards

	for _, uid := range room.UserIds {
		onePlayerCards := room.Cards[0:5]
		tmpCards := room.Cards[5:]
		room.Cards = tmpCards
		room.UserState[uid] = PlayerState{Uid: uid, Cards: onePlayerCards, Name: room.UserState[uid].Name, Status: 0}
	}
}

//摸牌
func (p *PlayerModuel) draw(uid int, room *Room) bool {
	if len(room.Cards) > 0 {
		card := room.Cards[0]
		tmpCard := room.Cards[1:]
		room.Cards = tmpCard
		newCard := append(room.UserState[uid].Cards, card)
		room.UserState[uid] = PlayerState{Uid: uid, Cards: newCard, Name: room.UserState[uid].Name, Status: 0}
		room.TurnTime = time.Now().Unix()

		room.Turn += 1
		if room.Turn > len(room.UserIds) {
			room.Turn = 1
		}
	} else {
		return false
	}

	return true
}

//玩家打牌
func (p *PlayerModuel) discard(uid int, card int, room *Room) bool {
	//验证颜色和形状
	roomCards := conf.BaseCards
	//cardInfo := roomCards[card]
	if roomCards[card]["color"] == roomCards[room.Center]["color"] || roomCards[card]["pic"] == roomCards[room.Center]["pic"] {
		var cardindex int
		for k, v := range room.UserState[uid].Cards {
			if card == v {
				cardindex = k
				break
			}
		}
		//删除你打的那张牌
		newCard := append(room.UserState[uid].Cards[:cardindex], room.UserState[uid].Cards[cardindex+1:]...)
		room.UserState[uid] = PlayerState{Uid: uid, Cards: newCard, Name: room.UserState[uid].Name, Status: 0}
		room.Center = card
		room.TurnTime = time.Now().Unix()
		room.Turn += 1
		if room.Turn > len(room.UserIds) {
			room.Turn = 1
		}
	} else {
		return false
	}

	return true
}

//游戏开始
func (p *PlayerModuel) start(room *Room) {
	(*room).State = PLAYING
	//初始化牌
	sortCards := p.initPlayerCards()
	//给每个人发牌
	p.sendCard(room, sortCards)

	//翻出第一张牌
	room.Center = room.Cards[0]
	tmpCards := room.Cards[1:]
	room.Cards = tmpCards

	room.Turn = 1
	room.TurnTime = time.Now().Unix()
	m := make(chan map[string]interface{})
	room.RecvCh = m

	log.Debug("resroom...", room)
	go p.StartPlay(room)
	//tmp := make(map[string]interface{})
	//room.RecvCh <- tmp
}

//如果玩家30s内没有发牌，则自动执行  room没有收到消息
func (p *PlayerModuel) StartPlay(room *Room) {
	stopFlag := false
	for {
		if stopFlag {
			break
		}
		//底牌出完，清算结局
		if len(room.Cards) == 0 {
			winerUid := p.CalculateEnd(room)
			roomModuel := RoomModule{}
			roomModuel.pushSuccessToOthers(winerUid, room)
			stopFlag = true
			break
		}

		//查看有没有人的牌出完了
		for k, _ := range room.UserState {
			if len(room.UserState[k].Cards) == 0 {
				winerUid := k
				roomModuel := RoomModule{}
				roomModuel.pushSuccessToOthers(winerUid, room)
				stopFlag = true
				room.State = 2

				break
			}
		}

		//监听玩家有没有出牌
		if !stopFlag {
			overCh := make(chan interface{})
			defer close(overCh)
			go p.PlayerSelTime(room, overCh)
			overRet := <-overCh
			if overRet == true {
				//自动出牌
				log.Debug("autoSendCard...")
				p.AutoSelCard(room)
			}
		} else {
			close(room.RecvCh)
			break
		}
	}
}

//检测超时
func (p *PlayerModuel) PlayerSelTime(room *Room, retCh chan interface{}) {
	for {
		select {
		case recvMsg := <-room.RecvCh:
			retCh <- recvMsg
			break
		case <-time.After(time.Second * 10):
			fmt.Println("10s time voer")
			retCh <- true
			break
		}
	}
}

//清算结局
func (p *PlayerModuel) CalculateEnd(room *Room) int {
	min := MaxBaseCardNum
	var winerUid int
	for uid, playerState := range room.UserState {
		cardNum := len(playerState.Cards)
		if cardNum < min {
			min = cardNum
			winerUid = uid
		}
	}

	return winerUid
}

func (p *PlayerModuel) AutoSelCard(room *Room) {
	//先检测目前是该谁出牌
	TurnUserPos := room.Turn
	turnUserId := room.UserIds[TurnUserPos-1]
	//判断有没有可以出的牌
	flag := false
	for _, v := range room.UserState[turnUserId].Cards {
		flag = p.discard(turnUserId, v, room)
		if flag {
			break
		}
	}

	if !flag {
		p.draw(turnUserId, room)
	}

	roomModuel := RoomModule{}
	roomModuel.pushRoomMsgToOthers(0, room)
}

//检查房间消息正确性
func (p *PlayerModuel) checkRoomMsg(msg *msg.Dispatch) error {
	//判断参数
	uid := msg.Uid
	params := msg.Params
	if _, ok := params["mtype"]; !ok {
		return errors.New("param mtype is not exist")
	}

	//判断该用户是否有出牌权
	roomModuel := RoomModule{}
	room, ok := roomModuel.getRoomByUid(uid)
	if !ok {

	}

	selectTurn := 0
	for _, roomUserId := range room.UserIds {
		if roomUserId == uid {
			selectTurn++
			break
		}
	}

	if selectTurn != room.Turn {
		return errors.New("the user have no such privace")
	}

	return nil
}
