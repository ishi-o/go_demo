package log

import (
	"fmt"
	"os"

	"github.com/ishi-o/go_demo/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func Init(cfg *config.LogConfig) error {
	level := getLogLevel(cfg.Level)

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	consoleWriter := zapcore.AddSync(os.Stdout)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writer, level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), consoleWriter, level),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Printf(format string, args ...interface{}) {
	Logger.Info(fmt.Sprintf(format, args...))
}
