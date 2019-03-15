package logs

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var file = new(lumberjack.Logger)
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger
var Debug *log.Logger

func SetLog(lp string, ms int, mb int, ma int) {
	file = &lumberjack.Logger{
		Filename:   lp,
		MaxSize:    ms,
		MaxBackups: mb,
		MaxAge:     ma,
		LocalTime:  true,
		Compress:   true,
	}
	//初始化错误日志记录器
	Info = log.New(io.MultiWriter(file, os.Stdout), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	//Info = log.New(file, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(io.MultiWriter(file, os.Stdout), "Warn:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stdout), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	//Error = log.New(file, "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(io.MultiWriter(file, os.Stdout), "Debug:", log.Ldate|log.Ltime|log.Lshortfile)

}
