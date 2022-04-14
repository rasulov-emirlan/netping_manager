package logger

import (
	"go.uber.org/zap"
)

func NewZap(filename string, isDev bool, level zap.AtomicLevel) (*zap.SugaredLogger, error) {
	cfg := zap.Config{
		Level:       level,
		Development: isDev,
		Encoding:    "console",
		OutputPaths: []string{filename},
	}
	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return z.Sugar(), nil
}
