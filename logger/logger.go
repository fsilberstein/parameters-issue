package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LogStdOut *zap.SugaredLogger
	LogStdErr *zap.SugaredLogger
)

func init() {
	newLogger()
}

//new : instantiate a new logger
func newLogger() {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.OutputPaths = []string{"stdout"}
	stdOutLog, _ := conf.Build()
	LogStdOut = stdOutLog.Sugar()

	conf.OutputPaths = []string{"stderr"}
	stdErrLog, _ := conf.Build()
	LogStdErr = stdErrLog.Sugar()
}
