package rlock

import (
	"log"
	"os"
)

type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type logx struct {
	infoL  logger
	errorL logger
	debugL logger
}

func newLogger() *logx {
	Info := log.New(os.Stdout, "[Info]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error := log.New(os.Stdout, "[Error]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug := log.New(os.Stdout, "[Debug]: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &logx{infoL: Info, errorL: Error, debugL: Debug}
}

func (l *logx) Info(v ...any) {
	l.infoL.Println(v...)
}

func (l *logx) Infof(format string, v ...any) {
	l.infoL.Printf(format, v...)
}

func (l *logx) Error(v ...any) {
	l.errorL.Println(v...)
}

func (l *logx) Errorf(format string, v ...any) {
	l.errorL.Printf(format, v...)
}

func (l *logx) Debug(v ...any) {
	l.debugL.Println(v...)
}

func (l *logx) Debugf(format string, v ...any) {
	l.debugL.Printf(format, v...)
}
