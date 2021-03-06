package manager

import (
	"context"
	"errors"
	"sort"

	"go.uber.org/zap"
)

type Service interface {
	// Probably this code
	// could be faster without pointers XP
	AddSocket(ctx context.Context, socket Socket, locationID int) (*Socket, error)
	UpdateSocket(ctx context.Context, socket Socket, socketID int) (*Socket, error)
	ListAllSockets(ctx context.Context) ([]*Location, error)
	RemoveSocket(ctx context.Context, socketID int) error

	CheckAll(ctx context.Context, locationID int) ([]*Socket, error)
	ToggleSocket(ctx context.Context, socketID int, onOrOff bool) error
}

type Repository interface {
	CreateSocket(ctx context.Context, s Socket, locationdID int) (*Socket, error)
	UpdateSocket(ctx context.Context, s Socket) (*Socket, error)
	DeleteSocket(ctx context.Context, socketID int) error

	FindSocketsByLocationID(ctx context.Context, locationID int) ([]*Socket, error)
	FindSocketByID(ctx context.Context, socketID int) (*Socket, error)
	ListAllSockets(ctx context.Context) ([]*Location, error)
}

type Sentry interface {
	CheckSocket(ctx context.Context, mib []string, address, community string, port int) ([]bool, error)
	ToggleSocket(ctx context.Context, turnOn bool, socketMIB, address, community string, port int) (*Socket, error)
}

type service struct {
	sentry Sentry
	repo   Repository
	log    *zap.SugaredLogger
}

func NewService(sentry Sentry, l *zap.SugaredLogger, repo Repository) (Service, error) {
	if sentry == nil || l == nil {
		return nil, errors.New("manager: arguments for ManagerNewService cannot be nil")
	}
	return &service{
		sentry: sentry,
		log:    l,
		repo:   repo,
	}, nil
}

func (s *service) CheckAll(ctx context.Context, locationID int) ([]*Socket, error) {
	defer s.log.Sync()
	s.log.Info("ManagerService: CheckAll()")
	sock, err := s.repo.FindSocketsByLocationID(ctx, locationID)
	if err != nil {
		s.log.Errorw("ManagerService: CheckAll() - repo call", zap.String("error", err.Error()))
		return nil, err
	}
	oids := []string{}
	for _, v := range sock {
		oids = append(oids, v.SNMPmib)
	}
	if len(sock) == 0 {
		return nil, errors.New("manager: no sockets")
	}
	checks, err := s.sentry.CheckSocket(ctx, oids, sock[0].SNMPaddress, "SWITCH", 161)
	if err != nil {
		s.log.Errorw("ManagerService: CheckAll() - sentry call", zap.String("error", err.Error()))
		return nil, err
	}
	for i, v := range sock {
		v.IsON = checks[i]
	}
	return sock, nil
}

func (s *service) ListAllSockets(ctx context.Context) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("ManagerService: ListAllSockets()")
	l, err := s.repo.ListAllSockets(ctx)
	if err != nil {
		s.log.Errorw("ManagerService: ListAllSockets()", zap.String("error", err.Error()))
		return nil, err
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i].ID < l[j].ID
	})
	return l, nil
}

func (s *service) AddSocket(ctx context.Context, socket Socket, locationID int) (*Socket, error) {
	defer s.log.Sync()
	s.log.Info("ManagerService: AddSocket()")
	socks, err := s.repo.FindSocketsByLocationID(ctx, locationID)
	if err != nil {
		s.log.Errorw("ManagerService: AddSocket() - repo call", zap.String("error", err.Error()))
		return nil, err
	}
	for _, v := range socks {
		if v.SNMPmib == socket.SNMPmib {
			s.log.Infow("ManagerService: AddSocket() - tryied to use mib address that is already in use by this location")
			return nil, errors.New("manager: dublicate mib address")
		}
	}
	sock, err := s.repo.CreateSocket(ctx, socket, locationID)
	if err != nil {
		s.log.Errorw("ManagerService: AddSocket() - repo call", zap.String("error", err.Error()))
		return nil, err
	}

	return sock, nil
}

func (s *service) RemoveSocket(ctx context.Context, socketID int) error {
	defer s.log.Sync()
	s.log.Info("ManagerService: RemoveSocket()")
	if err := s.repo.DeleteSocket(ctx, socketID); err != nil {
		s.log.Errorw("ManagerService: RemoveSocket() - repo call", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *service) ToggleSocket(ctx context.Context, socketID int, onOrOff bool) error {
	defer s.log.Sync()
	s.log.Info("ManagerService: ToggleSocket()")
	socket, err := s.repo.FindSocketByID(ctx, socketID)
	if err != nil {
		s.log.Errorw("ManagerService, ToggleSocket() - repo call", zap.String("error", err.Error()))
		return err
	}
	socket.SNMPcommunity = "SWITCH"
	socket.SNMPport = 161
	_, err = s.sentry.ToggleSocket(ctx, onOrOff, socket.SNMPmib, socket.SNMPaddress, socket.SNMPcommunity, socket.SNMPport)
	if err != nil {
		s.log.Errorw("ManagerService, ToggleSocket() - sentry call", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *service) UpdateSocket(ctx context.Context, socket Socket, socketID int) (*Socket, error) {
	defer s.log.Sync()
	s.log.Info("ManagerService: UpdateSocket()")
	socket.ID = socketID
	// TODO: add some updates in repositories
	sock, err := s.repo.UpdateSocket(ctx, socket)
	if err != nil {
		s.log.Errorw("ManagerService: UpdateSocket()", zap.String("error", err.Error()))
		return nil, err
	}
	return sock, err
}
