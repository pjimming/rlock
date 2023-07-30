package rlock

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type DLog struct {
	infoL  Logger
	errorL Logger
	debugL Logger
}

func init() {
	Info := log.New(os.Stdout, "[Info]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error := log.New(os.Stdout, "[Error]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug := log.New(os.Stdout, "[Debug]: ", log.Ldate|log.Ltime|log.Lshortfile)

	dlog = &DLog{infoL: Info, errorL: Error, debugL: Debug}
}

func (l *DLog) Info(v ...any) {
	l.infoL.Println(v...)
}

func (l *DLog) Infof(format string, v ...any) {
	l.infoL.Printf(format, v...)
}

func (l *DLog) Error(v ...any) {
	l.errorL.Println(v...)
}

func (l *DLog) Errorf(format string, v ...any) {
	l.errorL.Printf(format, v...)
}

func (l *DLog) Debug(v ...any) {
	l.debugL.Println(v...)
}

func (l *DLog) Debugf(format string, v ...any) {
	l.debugL.Printf(format, v...)
}
