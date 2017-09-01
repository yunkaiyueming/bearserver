package internal

import (
	"github.com/name5566/leaf/gate"
	"fmt"
)

var ConnMap = make(map[gate.Agent]int)

func RegNewConn(a gate.Agent, uid int) {
	ConnMap[a] = uid
}

func LeaveConn(a gate.Agent) {
	fmt.Println("laozi zou le ")
	delete(ConnMap, a)
}
