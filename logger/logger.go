package logger

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger *zap.SugaredLogger

func InitLogger(logLevel string, writeLogToFile bool) {
	fmt.Printf("logger init with loglevel:%s", logLevel)
	level := zapcore.DebugLevel
	if len(logLevel) <= 0 {
		level = zapcore.DebugLevel
	} else {
		tmp := strings.TrimSpace(logLevel)
		if strings.EqualFold(tmp, "error") {
			level = zapcore.ErrorLevel
		} else if strings.EqualFold(tmp, "warn") {
			level = zapcore.WarnLevel
		} else if strings.EqualFold(tmp, "info") {
			level = zapcore.InfoLevel
		} else {
			level = zapcore.DebugLevel
		}
	}
	writerSyncer := getLogWriter()
	encoder := getEncoder()

	var core zapcore.Core
	if writeLogToFile {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writerSyncer, level),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	} else {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	appLogger = logger.Sugar()
	defer appLogger.Sync() //Xả hết buffer ra
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Write log to file using gopkg.in/natefinch/lumberjack.v2
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/news-api.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	appLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	appLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	appLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	appLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	appLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	appLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	appLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	appLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	appLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	appLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	appLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	appLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	appLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	appLogger.Fatalf(template, args...)
}
