package compositors

import (
	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/internal/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
	managerH "github.com/rasulov-emirlan/netping-manager/internal/manager/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/manager/sentry"
	"github.com/rasulov-emirlan/netping-manager/internal/manager/storage/mysql"
	"github.com/rasulov-emirlan/netping-manager/pkg/db"
	"github.com/rasulov-emirlan/netping-manager/pkg/logger"
	"go.uber.org/zap"
)

type CloserFunc func() error

func NewManager(cfg config.Config) (rest.Registrator, CloserFunc, error) {
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
	// 	return nil, err
	// }
	stry := sentry.Sentry{}
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	z, err := logger.NewZap("logs.log", false, level)
	if err != nil {
		return nil, nil, err
	}
	dbConn, err := db.NewMySQL(cfg.Database)
	if err != nil {
		return nil, nil, err
	}
	repo, err := mysql.NewRepository(dbConn)
	if err != nil {
		return nil, nil, err
	}
	s, err := manager.NewService(&stry, z, repo)
	if err != nil {
		return nil, nil, err
	}
	h, err := managerH.NewHandler(s)
	if err != nil {
		return nil, nil, err
	}

	close := func() error {
		if err := dbConn.Close(); err != nil {
			return err
		}
		return nil
	}
	return h, close, nil
}