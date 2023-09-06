package logger

import (
	"mceasy/service-demo/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(cfg config.Config) (*zap.SugaredLogger, error) {
	logLevel, exists := loggerLevelMap[cfg.Logger.Level]
	if !exists {
		logLevel = zapcore.DebugLevel
	}

	var config zap.Config
	if cfg.Logger.Mode == "development" {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	} else if cfg.Logger.Mode == "production" {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}

	config.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	sugarLogger := logger.Sugar()

	return sugarLogger, nil
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}
