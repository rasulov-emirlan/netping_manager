package compositors

import (
	"database/sql"

	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/internal/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/users"
	usersH "github.com/rasulov-emirlan/netping-manager/internal/users/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/users/storage/mysql"
	"go.uber.org/zap"
)

func NewUsers(cfg *config.Config, logger *zap.SugaredLogger, dbConn *sql.DB) (rest.Registrator, error) {
	repo, err := mysql.NewRepository(dbConn)
	if err != nil {
		return nil, err
	}
	s, err := users.NewService(repo, logger)
	if err != nil {
		return nil, err
	}
	h, err := usersH.NewHandler(s, cfg.Server.JWTkey, logger)
	return h, err
}
