package zap

import "go.uber.org/zap/zapcore"

var encoderCfg = zapcore.EncoderConfig{
	MessageKey:  "msg",
	NameKey:     "name",
	LevelKey:    "level",
	EncodeLevel: zapcore.LowercaseLevelEncoder,
	CallerKey:   "caller",
	TimeKey:     "time",
	EncodeTime:  zapcore.ISO8601TimeEncoder,
}
