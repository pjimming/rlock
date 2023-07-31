package logx

import (
	"log"
	"os"
)

type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type Logger struct {
	infoL  logger
	errorL logger
	debugL logger
}

func NewLogger() *Logger {
	Info := log.New(os.Stdout, "[Info]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error := log.New(os.Stdout, "[Error]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug := log.New(os.Stdout, "[Debug]: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{infoL: Info, errorL: Error, debugL: Debug}
}

func (l *Logger) Info(v ...any) {
	l.infoL.Println(v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.infoL.Printf(format, v...)
}

func (l *Logger) Error(v ...any) {
	l.errorL.Println(v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.errorL.Printf(format, v...)
}

func (l *Logger) Debug(v ...any) {
	l.debugL.Println(v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.debugL.Printf(format, v...)
}
