package mlog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
)

var _Pid = os.Getpid()

//shell color
const (
	RED    = 31
	GREEN  = 32
	YELLOW = 33
	WHITE  = 37
)

//log levels
const (
	NONE = iota
	TRACE
	DEBUG
	INFO
	ERROR
	FATAL
)

var levelPrefix = map[int]string{
	NONE:  "",
	TRACE: "[T]",
	DEBUG: "[D]",
	INFO:  "[I]",
	ERROR: "[E]",
	FATAL: "[F]",
}

//LoggerInterface define
type LoggerInterface interface {
	Init(conf string)
	WriteMsg(msg LogMsg)
	Flush()
	Close()
}

//Mlogger define
type Mlogger struct {
	mu        sync.Mutex
	level     int
	calldepth int
	out       LoggerInterface
}

type LogMsg struct {
	level int
	msg   string
}

func New(out LoggerInterface) *Mlogger {
	logger := new(Mlogger)
	logger.level = INFO
	logger.calldepth = 2
	logger.out = out
	return logger
}

func newConsoleLogger() *Mlogger {
	consoleLogger := NewConsoleLogger()
	logger := New(consoleLogger)
	return logger
}

func (this *Mlogger) Trace(format string, v ...interface{}) {
	this.writeMsg(TRACE, format, v...)
	return
}

func (this *Mlogger) Debug(format string, v ...interface{}) {
	this.writeMsg(DEBUG, format, v...)
	return
}

func (this *Mlogger) Info(format string, v ...interface{}) {
	this.writeMsg(INFO, format, v...)
	return
}

func (this *Mlogger) Error(format string, v ...interface{}) {
	this.writeMsg(ERROR, format, v...)
	return
}

func (this *Mlogger) Fatal(format string, v ...interface{}) {
	this.writeMsg(FATAL, format, v...)
	return
}

func (this *Mlogger) writeMsg(logLevel int, format string, v ...interface{}) error {
	if this.level > logLevel {
		return nil
	}

	var perfix string
	if this.calldepth > 0 && logLevel < INFO {
		_, file, line, ok := runtime.Caller(this.calldepth)
		if ok {
			_, filename := path.Split(file)
			perfix = fmt.Sprintf("[%d]%s %s:%d", _Pid, levelPrefix[logLevel], filename, line)
		} else {
			perfix = fmt.Sprintf("[%d]%s", _Pid, levelPrefix[logLevel])
		}
	} else {
		perfix = fmt.Sprintf("[%d]", _Pid)
	}

	msg := perfix + fmt.Sprintf(format, v...)
	logmsg := LogMsg{
		level: logLevel,
		msg:   msg,
	}

	if logLevel == FATAL {
		panic(logmsg.msg)
	}

	this.out.WriteMsg(logmsg)
	return nil
}

var dfl = newConsoleLogger()

func Trace(v ...interface{}) {
	dfl.Trace("%s", fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	dfl.Debug("%s", fmt.Sprint(v...))
}

func Info(v ...interface{}) {
	dfl.Info("%s", fmt.Sprint(v...))
}

func Error(v ...interface{}) {
	dfl.Error("%s", fmt.Sprint(v...))
}

func Fatal(v ...interface{}) {
	dfl.Fatal("%s", fmt.Sprint(v...))
}
