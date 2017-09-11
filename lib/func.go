package lib

import (
	"fmt"
	"runtime"
	"time"

	"github.com/astaxie/beego/logs"
	"reflect"
	"errors"
)

func GetMsgByCode(code int) string {
	fmt.Println(code)
	return ""
}

//根据索引删除切片的值
func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	} else if (length - 1) >= index {
		return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
	}
	return nil, errors.New("error")
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


