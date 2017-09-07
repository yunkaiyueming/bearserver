/*
	测试游戏框架而写的第一个接口
 */
package internal

import (
	"bearserver/msg"
	//"fmt"
	"github.com/name5566/leaf/log"
	"fmt"
	"strconv"
)

func handlePlayCard(args []interface{}) (*msg.Response) {
	m := args[0].(*msg.Dispatch)
	response := &msg.Response{Cmd: m.Cmd, Rnum: m.Rnum, Uid: m.Uid, Ret: 0}

	//判断参数
	uid := m.Uid
	params := m.Params
	mtype, ok := params["mtype"]
	if !ok {
		response.Ret = -3
		return response
	}
	log.Debug("step11",mtype)
	//mtype = int(mtype)
	//if mtype.(type)

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
	//card, _ := params["card"].(int)
	var card int
	log.Debug("card...",int(card))

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
	room,ok := roomModuel.getRoomByUid(uid)
	if !ok{
		response.Ret = -5
		return response
	}

	log.Debug("step3...")
	log.Debug("card...",card)


	if mtype == "1"{
		//发牌
		playRet := playerModuel.discard(uid,card,&room)
		if !playRet{
			response.Ret = -6
			return response
		}
	}else if mtype == "2"{
		//摸牌
		log.Debug("mtype2...",mtype)

		playRet := playerModuel.draw(uid,&room)
		if !playRet{
			response.Ret = -5
			return response
		}
		perroomInfo, _ := roomModuel.getRoomInfo(uid, room.RoomID)
		roomModuel.pushRoomMsgToOthers(uid,&room)

		response.Data = perroomInfo
	}

	perroomInfo, _ := roomModuel.getRoomInfo(uid, room.RoomID)
	roomModuel.pushRoomMsgToOthers(uid,&room)

	response.Data = perroomInfo

	return response
}
