package main

import (
	"djclogger/mylogger"
	"time"
)
var L mylogger.Logger
func main() {

	L = mylogger.NewLog()
	////log:=myConsoleLogger.NewLog("debug")
	//Log:=mylogger.NewFileLogger("Info","./","djc.log",1024)
	for {

		L.Debug("这是一条Debug日志,%d,%s",1000,"时间未来")
		L.Trace("这是一条Trace日志,%d,%s",1000,"时间未来")
		L.Info("这是一条Info日")
		L.Warning("这是一条Warning日志")
		L.Error("这是一条Error日志")
		L.Fatal("这是一条Fatal日志")
		time.Sleep(time.Second)
	}


}
