package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	leaflog "github.com/name5566/leaf/log"
)

var (
	// log conf
	LogFlag = log.LstdFlags

	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000

	//database
	hostname = "192.168.8.81:3307"
	username = "root"
	password = ""
	database = "bear"

	//base card
	BaseCards = make(map[int]map[string]string)
)

var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
	DBUrl       string
}

func init() {

}

func initServerCf() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		leaflog.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		leaflog.Fatal("%v", err)
	}
}

func initCardCf() {
	data, err := ioutil.ReadFile("conf/card.json")
	if err != nil {
		leaflog.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &BaseCards)
	if err != nil {
		leaflog.Fatal("%v", err)
	}
}

func initCodeCf() {

}
