package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init инициализирует логгер
func Init(level string, format string, file string) error {
	// Парсинг уровня логирования
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Создание конфигурации
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Настройка формата вывода
	if format == "console" {
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config.Encoding = "json"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Настройка файла логов
	if file != "" {
		// Создание директории для логов
		if err := os.MkdirAll(file, 0755); err != nil {
			return err
		}
		config.OutputPaths = append(config.OutputPaths, file)
	}

	// Создание логгера
	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger возвращает экземпляр логгера
func GetLogger() *zap.Logger {
	if logger == nil {
		// Создание логгера по умолчанию
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}
	return logger
}

// Debug логирует debug сообщение
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info логирует info сообщение
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn логирует warning сообщение
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error логирует error сообщение
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal логирует fatal сообщение и завершает программу
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Sync синхронизирует буферы логгера
func Sync() error {
	return GetLogger().Sync()
}
