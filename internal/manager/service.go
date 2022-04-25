package manager

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Service interface {
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
		return nil, errors.New("manager: arguments for NewService cannot be nil")
	}
	return &service{
		sentry: sentry,
		log:    l,
		repo:   repo,
	}, nil
}

func (s *service) CheckAll(ctx context.Context, locationID int) ([]*Socket, error) {
	defer s.log.Sync()
	s.log.Info("Service: CheckAll()")
	sock, err := s.repo.FindSocketsByLocationID(ctx, locationID)
	if err != nil {
		s.log.Errorw("Service: CheckAll() - repo call", zap.String("error", err.Error()))
		return nil, err
	}
	oids := []string{}
	for _, v := range sock {
		oids = append(oids, v.SNMPmib)
	}
	checks, err := s.sentry.CheckSocket(ctx, oids, sock[0].SNMPaddress, "SWITCH", 161)
	if err != nil {
		s.log.Errorw("Service: CheckAll() - sentry call", zap.String("error", err.Error()))
		return nil, err
	}
	for i, v := range sock {
		v.IsON = checks[i]
	}
	return sock, nil
}

func (s *service) ListAllSockets(ctx context.Context) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: ListAllSockets()")
	l, err := s.repo.ListAllSockets(ctx)
	if err != nil {
		s.log.Errorw("Service: ListAllSockets()", zap.String("error", err.Error()))
		return nil, err
	}
	return l, nil
}

func (s *service) AddSocket(ctx context.Context, socket Socket, locationID int) (*Socket, error) {
	defer s.log.Sync()
	s.log.Info("Service: AddSocket()")
	socks, err := s.repo.FindSocketsByLocationID(ctx, locationID)
	if err != nil {
		s.log.Errorw("Service: AddSocket() - repo call", zap.String("error", err.Error()))
		return nil, err
	}
	for _, v := range socks {
		if v.SNMPmib == socket.SNMPmib {
			s.log.Infow("Service: AddSocket() - tried to use mib adress that is already in use by this location")
			return nil, errors.New("manager: dublicate mib address")
		}
	}
	sock, err := s.repo.CreateSocket(ctx, socket, locationID)
	if err != nil {
		s.log.Errorw("Service: AddSocket() - repo call", zap.String("error", err.Error()))
		return nil, err
	}

	return sock, nil
}

func (s *service) RemoveSocket(ctx context.Context, socketID int) error {
	defer s.log.Sync()
	s.log.Info("Service: RemoveSocket()")
	if err := s.repo.DeleteSocket(ctx, socketID); err != nil {
		s.log.Errorw("Service: RemoveSocket() - repo call", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *service) ToggleSocket(ctx context.Context, socketID int, onOrOff bool) error {
	defer s.log.Sync()
	s.log.Info("Service: ToggleSocket()")
	socket, err := s.repo.FindSocketByID(ctx, socketID)
	if err != nil {
		s.log.Errorw("Service, ToggleSocket() - repo call", zap.String("error", err.Error()))
		return err
	}
	socket.SNMPcommunity = "SWITCH"
	socket.SNMPport = 161
	_, err = s.sentry.ToggleSocket(ctx, onOrOff, socket.SNMPmib, socket.SNMPaddress, socket.SNMPcommunity, socket.SNMPport)
	if err != nil {
		s.log.Errorw("Service, ToggleSocket() - sentry call", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *service) UpdateSocket(ctx context.Context, socket Socket, socketID int) (*Socket, error) {
	defer s.log.Sync()
	s.log.Info("Service: UpdateSocket()")
	socket.ID = socketID
	// TODO: add some updates in repositories
	sock, err := s.repo.UpdateSocket(ctx, socket)
	if err != nil {
		s.log.Errorw("Service: UpdateSocket()", zap.String("error", err.Error()))
		return nil, err
	}
	return sock, err
}
