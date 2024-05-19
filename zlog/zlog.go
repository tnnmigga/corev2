package zlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/tnnmigga/corev2/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	// 先按默认值临时创建一个logger
	Init()
}

func Init() {
	var logLevel zap.AtomicLevel
	err := logLevel.UnmarshalText([]byte(conf.String("zlog.level", "debug")))
	if err != nil {
		panic(fmt.Errorf("log Init level error: %v", err))
	}
	conf := zap.Config{
		Level:             logLevel,
		Development:       false,
		Encoding:          conf.String("log.encoding", "console"),
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{conf.String("zlog.stdout", "stdout")},
		ErrorOutputPaths:  []string{conf.String("zlog.stderr", "stderr")},
		DisableCaller:     false,
		DisableStacktrace: true,
	}
	conf.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05.000000"))
	}
	conf.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		index := strings.LastIndex(caller.Function, "/")
		encoder.AppendString(fmt.Sprintf("%s:%d", caller.Function[index+1:], caller.Line))
	}
	l, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Errorf("zlog Init conf build error: %v", err))
	}
	logger = l.Sugar()
}

func Logger() *zap.SugaredLogger {
	return logger
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warn(args ...any) {
	logger.Warn(args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Panic(args ...any) {
	logger.Panic(args...)
}

func Panicf(format string, args ...any) {
	logger.Panicf(format, args...)
}

// func Fatal(args ...any) {
// 	logger.Error(args...)
// 	proc.Exit()
// }

// func Fatalf(format string, args ...any) {
// 	logger.Errorf(format, args...)
// 	proc.Exit()
// }
