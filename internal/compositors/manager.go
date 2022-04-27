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
	// l := []*manager.Location{{
	// 	ID:            1,
	// 	Name:          "Ошская станция",
	// 	RealLocation:  "Город Ош ул.Бакаева",
	// 	SNMPaddress:   "192.168.0.100",
	// 	SNMPcommunity: "SWITCH",
	// 	SNMPport:      161,
	// 	Sockets: []*manager.Socket{{
	// 		ID:         1,
	// 		Name:       "Кондиционер",
	// 		SNMPmib:    ".1.3.6.1.4.1.25728.8900.1.1.3.4",
	// 		IsON:       false,
	// 		ObjectType: manager.TypeAC,
	// 	}},
	// }}
	// w, err := watcher.NewWatcher(l)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// _, err = w.ToggleSocket(context.TODO(), 1, 1, 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
