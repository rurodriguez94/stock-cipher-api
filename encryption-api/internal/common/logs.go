package common

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Log() *zap.Logger {
	_ = log.Sync()
	return log
}

func InitDefault() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Panic("logger cannot be initiate")
	}

	log = l
}
