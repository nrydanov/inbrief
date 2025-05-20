package log

import "go.uber.org/zap"

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	return err
}
