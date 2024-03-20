package sqly

import (
	"log"
	"time"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
)

var Lg Log

type Log interface {
	Debug(cost time.Duration, sql string, args ...any)
}

type defaultLog struct {
	Level Level
}

func SetLog(lg Log) {
	Lg = lg
}

func NewDefaultLog(l Level) *defaultLog {
	return &defaultLog{Level: l}
}

func (d *defaultLog) Debug(cost time.Duration, sql string, args ...any) {
	if d.Level == DebugLevel {
		log.Println(sql, args, "cost:"+cost.String())
	}
}
