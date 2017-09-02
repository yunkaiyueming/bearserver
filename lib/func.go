package lib

import (
	"fmt"
	"runtime"
	"time"

	"github.com/astaxie/beego/logs"
)

func GetMsgByCode(code int) string {
	fmt.Println(code)
	return ""
}

func DelSlice(sourSliceInterface interface{}, start int, lenth int) []interface{} {
	delSl, dstSl := make([]interface{}, 0), make([]interface{}, 0)

	sourSlice := sourSliceInterface.([]interface{})
	if start+lenth > len(sourSlice) {
		delSl = sourSlice[start:len(sourSlice)]
	} else {
		delSl = sourSlice[start : start+lenth]
	}

	copy(dstSl, sourSlice)
	sourSlice = nil
	sourSlice = dstSl[0:start]

	return delSl
}

func CreateCronJob(d time.Duration, f func() string, fNameInfo ...string) {
	Ticker := time.NewTicker(d)
	runLock := false

	var fName string
	if len(fNameInfo) > 0 {
		fName = fNameInfo[0]
	} else {
		pc, _, _, _ := runtime.Caller(1)
		fName = runtime.FuncForPC(pc).Name()
	}

	for {
		<-Ticker.C
		if !runLock {
			runLock = true

			start := time.Now()
			logs.Alert(fmt.Sprintf("fname:%v,start_time:%s\n", fName, start.Format("2006-01-02 15:04:05")))
			ret := f()
			end := time.Now()
			logs.Alert(fmt.Sprintf("fname:%v,end_time:%s,spend:%d s,ret:%s\n", fName, time.Now().Format("2006-01-02 15:04:05"), end.Unix()-start.Unix(), ret))

			runLock = false
		}
	}
}
