package internal

import (
	"fmt"

	"github.com/name5566/leaf/gate"
)

var ConnMap = make(map[gate.Agent]int)

func init() {

}

func beatHeart() {

}

func RegNewConn(a gate.Agent, uid int) {
	ConnMap[a] = uid
	fmt.Printf("new conn==>%s \n", a.RemoteAddr().String())

	fmt.Println("all conn list:")
	for k, _ := range ConnMap {
		fmt.Printf("%s \n", k.RemoteAddr().String())
	}
}

func LeaveConn(a gate.Agent) {
	fmt.Printf("remote conn leave: %s \n", a.RemoteAddr().String())
	delete(ConnMap, a)
}
