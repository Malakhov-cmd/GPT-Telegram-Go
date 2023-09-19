package logger

import (
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	log  *zap.Logger
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		log = createLogger()
	})

	return log
}

func createLogger() *zap.Logger {
	// Создание конфигурации для логгера
	config := zap.NewProductionConfig()

	// Устанавливаем пути для файлов логов
	debugLogPath := "../logs/debug/debug.log"
	infoLogPath := "../logs/info/info.log"
	errorLogPath := "../logs/error/error.log"

	// Настройка ротации и архивации файла debug.log
	debugLog := &lumberjack.Logger{
		Filename:   debugLogPath,
		MaxSize:    500, // максимальный размер файла в мегабайтах
		MaxBackups: 5,   // максимальное количество ротированных файлов
		MaxAge:     7,   // максимальный возраст файла в днях
		Compress:   true,
	}
	config.OutputPaths = append(config.OutputPaths, debugLogPath)

	// Настройка ротации и архивации файла info.log
	infoLog := &lumberjack.Logger{
		Filename:   infoLogPath,
		MaxSize:    300,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
	}
	config.OutputPaths = append(config.OutputPaths, infoLogPath)

	// Настройка ротации и архивации файла error.log
	errorLog := &lumberjack.Logger{
		Filename:   errorLogPath,
		MaxSize:    200,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
	}
	config.ErrorOutputPaths = append(config.ErrorOutputPaths, errorLogPath)

	// Устанавливаем уровни для файлов логов
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Устанавливаем уровень вывода логов в консоль
	consoleConfig := zap.NewProductionEncoderConfig()

	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Создаем ядро логгера
	core := zapcore.NewTee(
		// Ядро для записи в файлы
		zapcore.NewCore(zapcore.NewConsoleEncoder(config.EncoderConfig), zapcore.Lock(zapcore.AddSync(debugLog)), zap.NewAtomicLevelAt(config.Level.Level())),
		zapcore.NewCore(zapcore.NewConsoleEncoder(config.EncoderConfig), zapcore.Lock(zapcore.AddSync(infoLog)), zap.NewAtomicLevelAt(config.Level.Level())),
		zapcore.NewCore(zapcore.NewConsoleEncoder(config.EncoderConfig), zapcore.Lock(zapcore.AddSync(errorLog)), zap.NewAtomicLevelAt(zap.ErrorLevel)),
		// Ядро для записи в консоль
		zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel),
	)

	// Создаем логгер
	logger := zap.New(core)

	// Пример использования
	logger.Info("logger инициализирован")

	// Чистим ресурсы
	logger.Sync()

	return logger
}
