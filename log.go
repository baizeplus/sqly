package sqly

import "log"

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
)

var Lg Log

type Log interface {
	Debug(args ...interface{})
}

type defaultLog struct {
	Level Level
}

func NewDefaultLog(l Level) *defaultLog {
	return &defaultLog{Level: l}
}

func (d *defaultLog) Debug(args ...interface{}) {
	if d.Level == DebugLevel {
		log.Println(args)
	}
}
