package mlog

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	// log   *log.Logger
	out   *os.File
	color map[int]int
}

func NewConsoleLogger() LoggerInterface {
	return &ConsoleLogger{
		// log: log.New(os.Stdout, "", log.LstdFlags),
		out: os.Stdout,
		color: map[int]int{
			DEBUG: GREEN,
			ERROR: RED,
			FATAL: RED,
		},
	}
}

func (*ConsoleLogger) Init(conf string) {
	return
}

func (l *ConsoleLogger) WriteMsg(msg LogMsg) {
	if color, ok := l.color[msg.level]; ok {
		setConsleColor(color)
		defer unsetConsleColor()
	}
	// l.log.Println(msg.msg)
	l.out.WriteString(msg.msg)
	l.out.WriteString("\n")
	return
}

func (l *ConsoleLogger) Flush() {
	return
}

func (l *ConsoleLogger) Close() {
	return
}

func (l *ConsoleLogger) SetColor(logLevel int, color int) {
	l.color[logLevel] = color
	return
}

func setConsleColor(color int) {
	fmt.Printf("\033[%dm", color)
	return
}

func unsetConsleColor() {
	setConsleColor(0)
	return
}
