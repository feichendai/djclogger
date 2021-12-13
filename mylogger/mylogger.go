package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

//定义logger类型级别常量
const (
	UNKNOWN LogLevel=iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func getLogLevelString(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "Debug"
	case TRACE:
		return "Trace"
	case INFO:
		return "Info"
	case WARNING:
		return "Warning"
	case ERROR:
		return "Error"
	case FATAL:
		return "Fatal"
	default:
		return "Debug"
	}
}

func parseLogLevel(s string) (LogLevel,error) {
	s=strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG,nil
	case "trace":
		return TRACE,nil
	case "info":
		return INFO,nil
	case "warning":
		return WARNING,nil
	case "error":
		return ERROR,nil
	case "fatal":
		return FATAL,nil
	default:
		return UNKNOWN,errors.New("无效的日志级别")
	}
}
//根据caller(skip int )得到函数调用的行//此处为0级,每包裹一层函数增加一级
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc,file,lineNo,ok:=runtime.Caller(skip)//lineNo为根据skip对应的调用函数所在的行
	if !ok{
		fmt.Println("runtime.Caller() failed")
		return "", "", 0
	}
	funcName=strings.Split(runtime.FuncForPC(pc).Name(),".")[1]//获取函数名字
	fileName=path.Base(file)//获取路径文件名
	return
}

type Logger interface {
	Debug(format string,a ...interface{})
	Trace(format string,a ...interface{})
	Info(format string,a ...interface{})
	Warning(format string,a ...interface{})
	Error(format string,a ...interface{})
	Fatal(format string,a ...interface{})
}


func NewLog() Logger {
	conf :=new(Config)
	conf.InitConfig("./mylogger/conf.ini")//当前目录下的文件
	logtype:=conf.Read("config","logtype")
	logpath:=conf.Read("config","logpath")
	logname:=conf.Read("config","logname")
	logsize:=conf.Read("config","maxsize")
	loglevel:=conf.Read("config","loglevel")
	msize,err:=strconv.Atoi(logsize)
	if err !=nil{
		panic(err)
	}
	t := strings.ToLower(logtype)
	if t == "file" {
		return NewFileLogger(loglevel,logpath,logname,int64(msize)*1024)
	}else{
		return NewConsoleLog(loglevel)
	}
}
