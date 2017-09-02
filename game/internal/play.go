package internal

import (
	"fmt"
	"time"

	"bearserver/conf"
	"bearserver/lib"

	"github.com/name5566/leaf/gate"
)

type PlayerModuel struct{}

var MaxBaseCardNum = len(conf.BaseCards)

//洗牌
func (p *PlayerModuel) initPlayerCards() []Card {
	roomCards := conf.BaseCards

	sortCards := make([]Card, 0)
	for cardId, cardInfo := range roomCards {
		color := cardInfo["color"]
		pic := cardInfo["pic"]

		card := Card{cardId, color, pic}
		sortCards = append(sortCards, card)
	}
	return sortCards
}

//发牌
func (p *PlayerModuel) sendCard(room *Room, sortCards []Card, a gate.Agent) {
	newSortCards := make([]Card, 0)
	var end int

	for i, uid := range room.UserIds {
		startPos := i * OneRoomPlayerNum
		end = (i + 1) * OneRoomPlayerNum
		onePlayerCards := sortCards[startPos:end]

		room.UserState[uid] = PlayerState{Uid: uid, Cards: onePlayerCards, Status: 0}
	}
	//retCh
	copy(newSortCards, sortCards)
	sortCards = nil
	copy(sortCards, newSortCards[end:len(newSortCards)])
}

//游戏开始
func (p *PlayerModuel) start(room *Room, a gate.Agent) {
	(*room).State = PLAYING
	sortCards := p.initPlayerCards()
	p.sendCard(room, sortCards, a)

	//翻出第一张牌
	preCardS := lib.DelSlice(sortCards, 0, 1)
	preCard := preCardS[0].(Card)
	fmt.Println("first card", preCard)

	stopFlag := false
	winerUids := make([]int, 0)

	for {
		if stopFlag {
			break
		}

		for _, uid := range room.UserIds {
			//牌出完，清算结局
			if len(sortCards) == 0 {
				winerUids = p.CalculateEnd(room)
				fmt.Println("winer uids==>", winerUids)
				stopFlag = true
				break
			}

			if len(room.UserState[uid].Cards) == 0 {
				winerUids := []int{uid}
				fmt.Println("winer uids==>", winerUids)
				stopFlag = true
				break
			}

			overCh := make(chan interface{})
			go p.PlayerSelTime(room, overCh)
			overRet := <-overCh
			if overRet == true {
				//自动出牌
				p.AutoSelCard()
			}
			//解析牌，检测规则，更新玩家状态
			//room.UserState[uid].Cards
		}
	}

	a.WriteMsg(winerUids)
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
