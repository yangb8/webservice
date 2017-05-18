package logger

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLocation uint8

const (
	StdErr LogLocation = iota
	File
)

func SetLogger(l LogLocation) {
	switch l {
	case File:
		log.SetOutput(&lumberjack.Logger{
			Filename:   "./foo.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
	case StdErr:
		// do nothing, os.Stderr is the default Writer for log.Output in package log
		// var std = New(os.Stderr, "", LstdFlags)
	}
}
