package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/tnnmigga/corev2/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.SugaredLogger

func init() {
	// 先按默认值临时创建一个logger
	Init()
}

func Init() {
	var logLevel zap.AtomicLevel
	err := logLevel.UnmarshalText([]byte(conf.String("log.level", "debug")))
	if err != nil {
		panic(fmt.Errorf("log Init level error: %v", err))
	}
	conf := zap.Config{
		Level:             logLevel,
		Development:       false,
		Encoding:          conf.String("log.encoding", "console"),
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{conf.String("log.stdout", "stdout")},
		ErrorOutputPaths:  []string{conf.String("log.stderr", "stderr")},
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
		panic(fmt.Errorf("log Init conf build error: %v", err))
	}
	lg = l.Sugar()
}

func Logger() *zap.SugaredLogger {
	return lg
}

func Debug(args ...any) {
	lg.Debug(args...)
}

func Debugf(format string, args ...any) {
	lg.Debugf(format, args...)
}

func Info(args ...any) {
	lg.Info(args...)
}

func Infof(format string, args ...any) {
	lg.Infof(format, args...)
}

func Warn(args ...any) {
	lg.Warn(args...)
}

func Warnf(format string, args ...any) {
	lg.Warnf(format, args...)
}

func Error(args ...any) {
	lg.Error(args...)
}

func Errorf(format string, args ...any) {
	lg.Errorf(format, args...)
}

func Panic(args ...any) {
	lg.Panic(args...)
}

func Panicf(format string, args ...any) {
	lg.Panicf(format, args...)
}
