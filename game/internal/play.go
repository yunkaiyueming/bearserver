package internal

import (
	"fmt"
	"time"

	"bearserver/conf"
	//"bearserver/lib"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

type PlayerModuel struct{}

var MaxBaseCardNum = len(conf.BaseCards)

//洗牌
func (p *PlayerModuel) initPlayerCards() []int {
	roomCards := conf.BaseCards

	sortCards := make([]int, 0)
	for cardId,_ := range roomCards {
		sortCards = append(sortCards, int(cardId))
	}
	return sortCards
}

//发牌
func (p *PlayerModuel) sendCard(room *Room, sortCards []int) {
	//newSortCards := make([]Card, 0)
	//var end int

	log.Debug("sendCard...",room)
	log.Debug("sortCards...",sortCards)
	room.Cards = sortCards

	for _, uid := range room.UserIds {
		onePlayerCards := room.Cards[0:5]
		tmpCards := room.Cards[5:]
		room.Cards = tmpCards
		//room.UserState[uid].Cards = onePlayerCards
		room.UserState[uid] = PlayerState{Uid: uid, Cards: onePlayerCards, Name:room.UserState[uid].Name,Status: 0}
	}

	//retCh
	//copy(newSortCards, sortCards)
	//sortCards = nil
	//copy(sortCards, newSortCards[end:len(newSortCards)])
}

//摸牌
func (p *PlayerModuel) draw(uid int,room *Room) bool{
	if len(room.Cards) > 0{
		card := room.Cards[0]
		newCard := append(room.UserState[uid].Cards,card)
		room.UserState[uid] = PlayerState{Uid: uid, Cards: newCard, Name:room.UserState[uid].Name,Status: 0}
		room.TurnTime = time.Now().Unix()
		room.Turn += 1
		if room.Turn > len(room.UserIds){
			room.Turn = 1
		}
	}else {
		return false
	}

	return true
}

//发牌
func (p *PlayerModuel) discard(uid int,card int,room *Room) bool{
	//验证颜色和形状
	roomCards := conf.BaseCards
	//cardInfo := roomCards[card]
	if roomCards[card]["color"] == roomCards[room.Center]["color"] || roomCards[card]["pic"] == roomCards[room.Center]["pic"]{
		var cardindex int
		for k,v := range room.UserState[uid].Cards{
			if card == v{
				cardindex = k
				break
			}
		}
		//删除你打的那张牌
		newCard := append(room.UserState[uid].Cards[:cardindex], room.UserState[uid].Cards[cardindex+1:]...)
		room.UserState[uid] = PlayerState{Uid: uid, Cards: newCard, Name:room.UserState[uid].Name,Status: 0}
		log.Debug("cardtocenter...",card)
		room.Center = card
		room.TurnTime = time.Now().Unix()
		room.Turn += 1
		if room.Turn > len(room.UserIds){
			room.Turn = 1
		}
		log.Debug("room666...",room)
	}else {
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
	//preCardS := lib.DelSlice(sortCards, 0, 1)
	//preCard := preCardS[0].(Card)
	//tmpPre := room.Cards[0:1]

	//firstCard := tmpPre

	room.Center = room.Cards[0]
	room.Turn = 1
	room.TurnTime = time.Now().Unix()

	log.Debug("resroom...",room)
	time.After(time.Second * 30)

	//go p.StartPlay(room)


	//fmt.Println("first card", preCard)
	//
	//stopFlag := false
	//winerUids := make([]int, 0)

	//for {
	//	if stopFlag {
	//		break
	//	}
	//
	//	for _, uid := range room.UserIds {
	//		//牌出完，清算结局
	//		if len(sortCards) == 0 {
	//			winerUids = p.CalculateEnd(room)
	//			fmt.Println("winer uids==>", winerUids)
	//			stopFlag = true
	//			break
	//		}
	//
	//		if len(room.UserState[uid].Cards) == 0 {
	//			winerUids := []int{uid}
	//			fmt.Println("winer uids==>", winerUids)
	//			stopFlag = true
	//			break
	//		}
	//
	//		overCh := make(chan interface{})
	//		go p.PlayerSelTime(room, overCh)
	//		overRet := <-overCh
	//		if overRet == true {
	//			//自动出牌
	//			p.AutoSelCard()
	//		}
	//		//解析牌，检测规则，更新玩家状态
	//		//room.UserState[uid].Cards
	//	}
	//}

	//a.WriteMsg(winerUids)
}

//如果玩家30s内没有发牌，则自动执行  room没有收到消息
func(p *PlayerModuel)StartPlay(room *Room){
	stopFlag := false
	for {
		if stopFlag {
			break
		}
		//监听玩家有没有出牌
		overCh := make(chan interface{})
		go p.PlayerSelTime(room, overCh)
		overRet := <-overCh
		if overRet == true {
			//自动出牌
			log.Debug("autoSendCard...")
			p.AutoSelCard()
		}
	}
}

func (p *PlayerModuel) PlayerSelTime(room *Room, retCh chan interface{}) {
	for {
		select {
			case recvMsg := <-room.RecvCh:
				retCh <- recvMsg
				break
			case <-time.After(time.Second * 30):
				fmt.Println("3s time voer")
				retCh <- true
				break
		}
	}
}

//清算结局
func (p *PlayerModuel) CalculateEnd(room *Room) []int {
	min := MaxBaseCardNum
	winerUidS := make([]int, 0)
	for uid, playerState := range room.UserState {
		cardNum := len(playerState.Cards)
		if min > cardNum {
			min = cardNum
			winerUidS = []int{uid}
		} else if min == cardNum {
			winerUidS = append(winerUidS, uid)
		}
	}

	return winerUidS
}

func (p *PlayerModuel) AutoSelCard() {
}

func (p *PlayerModuel) RecvRoomMsg(a gate.Agent, roomId int, params map[string]interface{}) {
	room, err := (&RoomModule{}).GetRoom(roomId)
	if err != nil {
		a.WriteMsg(err.Error())
		return
	}
	room.RecvCh <- params
}
