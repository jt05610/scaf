package zap

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Logger struct {
	z *zap.Logger
}

type LogLevel uint8

func newLogger(ctx context.Context, name string, l zapcore.Level) *Logger {
	z := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(
				&lumberjack.Logger{
					Filename:   name,
					MaxSize:    100,
					MaxAge:     7,
					MaxBackups: 5,
					LocalTime:  true,
					Compress:   true,
				}),
			l,
		),
	)

	z = z.WithOptions(
		zap.WrapCore(
			func(core zapcore.Core) zapcore.Core {
				ucEncoder := encoderCfg
				ucEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
				return zapcore.NewTee(
					core,
					zapcore.NewCore(
						zapcore.NewConsoleEncoder(ucEncoder),
						zapcore.Lock(os.Stdout),
						zapcore.DebugLevel,
					),
				)
			}))
	go func() {
		<-ctx.Done()
		err := z.Sync()
		if err != nil {
			panic(err)
		}
	}()
	return &Logger{
		z: z,
	}

}

func newSampler(ctx context.Context, name string, l zapcore.Level, freq time.Duration, initial, skip int) *Logger {
	z := zap.New(
		zapcore.NewSamplerWithOptions(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.AddSync(
					&lumberjack.Logger{
						Filename:   name,
						MaxSize:    100,
						MaxAge:     7,
						MaxBackups: 5,
						LocalTime:  true,
						Compress:   true,
					}),
				l,
			), freq, initial, skip,
		),
	)

	z = z.WithOptions(
		zap.WrapCore(
			func(core zapcore.Core) zapcore.Core {
				ucEncoder := encoderCfg
				ucEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
				return zapcore.NewTee(
					core,
					zapcore.NewCore(
						zapcore.NewConsoleEncoder(ucEncoder),
						zapcore.Lock(os.Stdout),
						zapcore.DebugLevel,
					),
				)
			}))
	go func() {
		<-ctx.Done()
		err := z.Sync()
		if err != nil {
			panic(err)
		}
	}()
	return &Logger{
		z: z,
	}

}

func NewProd(ctx context.Context, name string) *Logger {
	l := newLogger(ctx, name, zapcore.InfoLevel)
	return l
}

func NewDev(ctx context.Context, name string) *Logger {
	l := newLogger(ctx, name, zapcore.DebugLevel)
	return l
}

func NewProdSampling(ctx context.Context, name string, freq time.Duration, initial, skip int) *Logger {
	l := newSampler(ctx, name, zapcore.InfoLevel, freq, initial, skip)
	return l
}

func NewDevSampling(ctx context.Context, name string, freq time.Duration, initial, skip int) *Logger {
	l := newSampler(ctx, name, zapcore.DebugLevel, freq, initial, skip)
	return l
}
