/*
	推送类，关于推送的方法都会定义在这个里面
*/
package internal

import (
	"bearserver/msg"
)

type PushMsgModuel struct{}

func (p *PushMsgModuel) GetPushMsgCmd() string {
	return "pushMsg"
}

func (p *PushMsgModuel) GetPushMsgRNum() int {
	return 1
}

//根据玩家的uid,给玩家推送消息
func (p *PushMsgModuel) pushMsgByUid(uid int, msgInfo interface{}) (err error) {
	response := &msg.Response{Cmd: p.GetPushMsgCmd(), Rnum: p.GetPushMsgRNum(), Uid: uid}
	for a, userid := range ConnMap {
		if uid == userid {
			response.Data = msgInfo
			a.WriteMsg(response)
			break
		}
	}
	return nil
}
