package logger

import (
	"go.uber.org/zap"
	"os"
)

var Log *zap.Logger

func InitLogger(isProduction bool) {
	var err error

	if isProduction {
		Log, err = zap.NewProduction()
	} else {
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			err := os.MkdirAll("logs", os.ModePerm)
			if err != nil {
				panic("failed to create logs directory: " + err.Error())
			}
		}

		cfg := zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			OutputPaths:      []string{"stdout", "logs/app.log"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig:    zap.NewProductionEncoderConfig(),
		}
		Log, err = cfg.Build()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}
