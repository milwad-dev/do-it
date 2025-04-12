package logger

import "go.uber.org/zap"

var Log *zap.Logger

func InitLogger(isProduction bool) {
	var err error

	if isProduction {
		Log, err = zap.NewProduction()
	} else {
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
