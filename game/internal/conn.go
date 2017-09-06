package internal

import (
	"github.com/name5566/leaf/gate"
	"fmt"
)

var ConnMap = make(map[gate.Agent]int)

func init() {

}

func beatHeart() {

}

func RegNewConn(a gate.Agent, uid int) {
	ConnMap[a] = uid
	//fmt.Println("v%\n",ConnMap)
	for k, v := range ConnMap {
		fmt.Println(k, v)
	}
}

func LeaveConn(a gate.Agent) {
	delete(ConnMap, a)
}
