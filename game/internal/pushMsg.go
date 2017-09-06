/*
	推送类，关于推送的方法都会定义在这个里面


	joke：
	有一次我们公司组织我们出去玩，在山脚的水塘我遇到一只青蛙，它对我说：“切尔西2：1曼联”，然后就跳走了。当时我很奇怪，
	不知道为什么青蛙会说话，而且说的好像是足球比分。那天晚上正好是英超，我看完之后傻了，因为比分就是切尔西2：1曼联！
	第二天我又去找那只青蛙，它又对我说“3D，147。”我立即买了50注，结果真的中了，一下子中了5万。第三天我又去找青蛙，
	这次青蛙和我说了7个数字，我知道那是双色球，我买了5注，握着税后的5000万元支票，我深深感到我有今天都是拜青蛙所赐，于是我决定报答它。
	我问青蛙有什么愿望，它说要去见见人类的浴池。于是我包了最豪华的洗浴中心，在休息的时候我带它来到按摩床，问它还有什么愿望。青蛙看着我，
	说：“别说话，吻我！”我想青蛙都给了我几千万，难道吻一下不行吗？于是我亲了青蛙一下，这个时候青蛙忽然变成了一个16岁的漂亮小女孩，
	说她被施了魔法，现在我 救了她，她愿意报答我……警察同志，我说的是真的，床上那个女孩就是这么来的，你们不要抓我啊！
 */
package internal

import (
	"bearserver/msg"
	"fmt"
)

var cmd  = "pushMsg"
var rnum  = 1
type PushMsgModuel struct{}

//根据玩家的uid,给玩家推送消息
func (p *PushMsgModuel) pushMsgByUid(uid int,msgInfo interface{})(err error) {
	response := &msg.Response{Cmd:cmd,Rnum:rnum,Uid:uid}
	fmt.Println()
	for a, userid := range ConnMap {
		if uid == userid {
			response.Data = msgInfo
			a.WriteMsg(response)
		}
	}

	return
}
