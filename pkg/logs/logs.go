package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""
	var err error
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(message interface{}, fields ...zapcore.Field) {
	log.Info(fmt.Sprintf("%+v", message), fields...)
}

func Debug(message interface{}, fields ...zapcore.Field) {
	log.Debug(fmt.Sprintf("%+v", message), fields...)
}

func Error(message interface{}, fields ...zapcore.Field) {
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), fields...)
	case string:
		log.Error(v, fields...)
	default:
		log.Error(fmt.Sprintf("%v", v), fields...)
	}
}
