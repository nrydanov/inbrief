package log

import "go.uber.org/zap"

var L *zap.Logger

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	L = logger
	return err
}
