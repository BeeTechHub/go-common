package logger

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger *zap.SugaredLogger

// -------------- Start LOGGING ----------------------------
func GetEnv(envName string) string {
	err := godotenv.Load()
	if err != nil {
		Fatal("Error loading .env name:" + envName)
	}

	return os.Getenv(envName)
}

func ConfigWriteLogToFile() bool {
	result, error := strconv.ParseBool(GetEnv("LOG_WRITE_TO_FILE"))
	if error != nil {
		Warnf("ConfigWriteLogToFile error: %s", error.Error())
		return false
	}
	return result
}

// -------------- End LOGGING ----------------------------

func InitLogger() {
	writerSyncer := getLogWriter()
	encoder := getEncoder()

	var core zapcore.Core
	if ConfigWriteLogToFile() {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
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
