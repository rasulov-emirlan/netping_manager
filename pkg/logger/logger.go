package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CloseLogger func() error

func NewZap(filename string, isDev bool, level zap.AtomicLevel) (*zap.SugaredLogger, CloseLogger, error) {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	fileencoder := zapcore.NewJSONEncoder(conf)
	logFile, err := os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	writer := zapcore.AddSync(logFile)
	core := zapcore.NewTee(
		zapcore.NewCore(fileencoder, writer, level),
	)
	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	close := func() error {
		if err := logFile.Close(); err != nil {
			return err
		}
		return nil
	}
	return l.Sugar(), close, nil
}
