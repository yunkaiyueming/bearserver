package internal

import (
	"fmt"
	"time"

	"bearserver/conf"
	"github.com/name5566/leaf/log"
)

type PlayerModuel struct{}

var MaxBaseCardNum = len(conf.BaseCards)

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
		}else{
			close(room.RecvCh)
		}
	}
}

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
