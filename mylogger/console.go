package mylogger

import (
	"fmt"
	"time"
)

//往中断写日志相关内容

//logger 日志结构体
type LogLevel uint16


type ConsoleLogger struct {
	Level LogLevel
}
func NewConsoleLog(levelStr string) ConsoleLogger  {
	level,err:=parseLogLevel(levelStr)
	if err !=nil {
		panic(err)
	}
	return ConsoleLogger{
		Level:level,
	}
}
func (c ConsoleLogger) enable(loglevel LogLevel) bool {
	return c.Level<=loglevel
}
func (c ConsoleLogger)log(lv LogLevel, format string,a ...interface{}) {
	if c.enable(lv) {
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNo := getInfo(3)
		levelstr := getLogLevelString(lv)
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now, levelstr, fileName, funcName, lineNo, msg)
	}
}
func (c ConsoleLogger) Debug(format string,a ...interface{})  {
		c.log(DEBUG,format,a...)
}
func (c ConsoleLogger) Trace(format string,a ...interface{})  {
		c.log(TRACE,format,a...)
}
func (c ConsoleLogger) Info(format string,a ...interface{})  {
		c.log(INFO,format,a...)
}
func (c ConsoleLogger) Warning(format string,a ...interface{}) {
		c.log(WARNING,format,a...)
}
func (c ConsoleLogger) Error(format string,a ...interface{})  {
		c.log(ERROR,format,a...)
}
func (c ConsoleLogger) Fatal(format string,a ...interface{})  {
		c.log(FATAL,format,a...)
}
