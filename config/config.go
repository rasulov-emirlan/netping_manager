package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type (
	server struct {
		Port         string
		TimeoutRead  time.Duration
		TimeoutWrite time.Duration
		JWTkey       []byte
	}
	Config struct {
		Server      server
		Database    string
		LogFilename string
	}
)

const (
	serverPort         = "SERVER_PORT"
	serverTimeoutRead  = "SERVER_TIMEOUT_READ"
	serverTimeoutWrite = "SERVER_TIMEOUT_WRITE"
	serverJWTkey       = "SERVER_JWT_KEY"

	databaseURL = "DATABASE_URL"

	logFileName = "LOG_NAME"
)

var (
	ErrNoServerData   = errors.New("config: did not find configs for server")
	ErrNoDatabaseData = errors.New("config: did not find configs for database")
)

func NewConfig(filenames ...string) (*Config, error) {
	if len(filenames) > 0 {
		if err := godotenv.Load(filenames...); err != nil {
			return nil, err
		}
	}
	tR, err := time.ParseDuration(os.Getenv(serverTimeoutRead))
	if err != nil {
		return nil, ErrNoServerData
	}
	tW, err := time.ParseDuration(os.Getenv(serverTimeoutWrite))
	if err != nil {
		return nil, ErrNoServerData
	}
	cfg := Config{
		Server: server{
			Port:         os.Getenv(serverPort),
			TimeoutRead:  tR,
			TimeoutWrite: tW,
			JWTkey:       []byte(os.Getenv(serverJWTkey)),
		},
		Database:    os.Getenv(databaseURL),
		LogFilename: os.Getenv(logFileName),
	}
	if cfg.Server.Port == "" || len(cfg.Server.JWTkey) == 0 {
		return nil, ErrNoServerData
	}
	cfg.Server.Port = ":" + cfg.Server.Port
	if cfg.Database == "" {
		return nil, ErrNoDatabaseData
	}
	if cfg.LogFilename == "" {
		cfg.LogFilename = "logs.log"
	}
	return &cfg, nil
}
