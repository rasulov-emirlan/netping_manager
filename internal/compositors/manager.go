package compositors

import (
	"database/sql"

	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/internal/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
	managerH "github.com/rasulov-emirlan/netping-manager/internal/manager/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/manager/sentry"
	"github.com/rasulov-emirlan/netping-manager/internal/manager/storage/mysql"
	"go.uber.org/zap"
)

func NewManager(cfg config.Config, logger *zap.SugaredLogger, dbConn *sql.DB) (rest.Registrator, error) {
	// Sentry is a netping manager
	stry := sentry.Sentry{}
	repo, err := mysql.NewRepository(dbConn)
	if err != nil {
		return nil, err
	}
	s, err := manager.NewService(&stry, logger, repo)
	if err != nil {
		return nil, err
	}
	h, err := managerH.NewHandler(s, cfg.Server.JWTkey, logger)
	if err != nil {
		return nil, err
	}
	return h, nil
}
