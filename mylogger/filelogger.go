package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	Level LogLevel
	filePath string //日志文件保存的路径
	fileName string //日志文件保存的文件名字
	maxFileSize int64 //日志文件大小，超出后新建一个文件
	fileObj *os.File //所有日志 存放文件
	errFileObj *os.File //错误日志存放文件
}

func NewFileLogger(levelStr, fp, fn string, maxsize int64) *FileLogger {
	logLevel,err:=parseLogLevel(levelStr)
	if err!=nil{
		panic(err)
	}
	fl:= &FileLogger{
		Level:logLevel,
		filePath:fp,
		fileName:fn,
		maxFileSize:maxsize,
	}
	err=fl.initFile()//按照文件路径和文件名将文件打开
	if err !=nil{
		panic(err)
	}
	return fl
}
func (f *FileLogger)initFile()(error){
	fullFileName:=path.Join(f.filePath,f.fileName)
	fileObj,err:=os.OpenFile(fullFileName,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
	if err !=nil{
		fmt.Println("open log file failed, err:%v\n",err)
		return err
	}
	errFileObj,err:=os.OpenFile(fullFileName+".err",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644)
	if err !=nil{
		fmt.Println("open log file failed, err:%v\n",err)
		return err
	}
	//日志文件都已经打开
	f.fileObj=fileObj
	f.errFileObj=errFileObj
	return nil
}

func (f *FileLogger) enable(loglevel LogLevel) bool {
	return f.Level<=loglevel
}

func (f *FileLogger) checkSize(file *os.File) bool {

	fileInfo,err:=file.Stat()
	if err !=nil{
		//fileObj,err:=os.OpenFile(file.Name(),os.O_APPEND|os.O_CREATE|os.O_WRONLY,0777)
		//if err !=nil {
		fmt.Println("check size file stat err",err)
		return false
		//}
		//fileInfo,_ =fileObj.Stat()
	}
	fileSize:=fileInfo.Size()
	//如果当前文件大小大于等于文件的最大值就应该返回true
	if f.maxFileSize<fileSize {
		return true
	}

	return false
}
func (f *FileLogger) splitLogFile(file *os.File) (*os.File,error){
	//获取文件的状态信息
	//1.关闭当前的日志文件
	file.Close()
	//拿到当前日志文件完整路径
	logName:=file.Name()
	fmt.Println("file name:",logName)
	//新日志文件名字加入日期时间
	nowStr:=time.Now().Format("20060102150405000")
	newLogName:=fmt.Sprintf("%s.bak%s",logName,nowStr)

	//2.备份一下rename
	err:=os.Rename(logName,newLogName)
	if err !=nil {
		fmt.Println("os.rename err:",err)
		return nil,err
	}
	//3.打开一个新的日志文件

	fileObj,err:=os.OpenFile(logName,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0777)
	if err !=nil {
		fmt.Println("open log file failed, err:",err)
		return nil,err
	}
	//4.将打开的新日志文件返回
	return fileObj,nil
}

func (f *FileLogger) log(lv LogLevel, format string,a ...interface{}) {
	if f.enable(lv) {
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNo := getInfo(3)
		levelstr := getLogLevelString(lv)
		if f.checkSize(f.fileObj) {
			//需要切割日志文件
			newFile,err:=f.splitLogFile(f.fileObj)
			if err !=nil {
				fmt.Println("切割日志文件错误，err:",err)
				return
			}
			//fmt.Println("file name:",logName)

			f.fileObj=newFile
		}
		msg := fmt.Sprintf(format, a...)
		fmt.Println("output:",msg)
		_,err:=fmt.Fprintf(f.fileObj,"[%s] [%s] [%s:%s:%d] %s\n", now, levelstr, fileName, funcName, lineNo, msg)
		if err !=nil {
			e := fmt.Sprintf("f.fileObj err :%v",err)
			fmt.Println(e)
			panic(e)
		}
		//logInfo:=fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", now, levelstr, fileName, funcName, lineNo, msg)
		//_,err:=f.fileObj.WriteString(logInfo)
		//if err !=nil {
		//	e := fmt.Sprintf("f.fileObj err :%v",err)
		//	fmt.Println(e)
		//	panic(e)
		//}
		if lv>=ERROR {
			if f.checkSize(f.errFileObj) {
				//需要切割日志文件
				newFile, err := f.splitLogFile(f.errFileObj)
				if err != nil {
					fmt.Println("切割日志文件错误，err:", err)
					return
				}
				f.errFileObj = newFile
			}
			_, err := fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n", now, levelstr, fileName, funcName, lineNo, msg)
			//logInfo:=fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", now, levelstr, fileName, funcName, lineNo, msg)
			//_,err:=f.errFileObj.WriteString(logInfo)
			if err !=nil {
				e := fmt.Sprintf("f.errFIleObj err :%v", err)
				fmt.Println(e)
				panic(e)
			}
		}
	}
}
func (f *FileLogger) Debug(format string,a ...interface{})  {
	f.log(DEBUG,format,a...)
}
func (f *FileLogger) Trace(format string,a ...interface{})  {
	f.log(TRACE,format,a...)
}
func (f *FileLogger) Info(format string,a ...interface{})  {
	f.log(INFO,format,a...)
}
func (f *FileLogger) Warning(format string,a ...interface{}) {
	f.log(WARNING,format,a...)
}
func (f *FileLogger) Error(format string,a ...interface{})  {
	f.log(ERROR,format,a...)
}
func (f *FileLogger) Fatal(format string,a ...interface{})  {
	f.log(FATAL,format,a...)
}
func (f *FileLogger) Close()  {
	f.errFileObj.Close()
	f.fileObj.Close()

}
