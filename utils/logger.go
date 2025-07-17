package utils

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log *zap.Logger
)

// InitLogger 初始化日志
func InitLogger() error {
	// 创建logs目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// 配置日志轮转
	infoWriter := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "info.log"),
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     30, // 天
		Compress:   true,
	}

	errorWriter := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "error.log"),
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     30, // 天
		Compress:   true,
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建Core
	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(infoWriter),
		zapcore.InfoLevel,
	)

	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(errorWriter),
		zapcore.ErrorLevel,
	)

	// 创建Logger
	core := zapcore.NewTee(infoCore, errorCore)
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// timeEncoder 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Info 记录信息日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Error 记录错误日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Warn 记录警告日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Debug 记录调试日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Fatal 记录致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
