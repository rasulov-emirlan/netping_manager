package manager

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Service interface {
	AddSocket(ctx context.Context, socket Socket, locationAddress string) (*Socket, error)
	RemoveSocket(ctx context.Context, socketID int) error

	CheckAll(ctx context.Context, locationAddress string) ([]*Socket, error)
	ToggleSocket(ctx context.Context, socketID int, onOrOff bool) error
}

type Repository interface {
	CreateSocket(ctx context.Context, s Socket) (*Socket, error)
	UpdateSocket(ctx context.Context, s Socket) (*Socket, error)
	DeleteSocket(ctx context.Context, socketID int) error

	FindSocketsByLocation(ctx context.Context, locationAddress string) ([]*Socket, error)
	FindSocketByID(ctx context.Context, socketID int) (*Socket, error)
}

type Sentry interface{
	CheckSocket(ctx context.Context, mib, address, community string, port int) (*Socket, error)
	ToggleSocket(ctx context.Context, turnOn bool, socketMIB, address, community string, port int) (*Socket, error)
}

type service struct {
	sentry   Sentry
	repo Repository
	log *zap.SugaredLogger
}

func NewService(sentry Sentry, l *zap.SugaredLogger, repo Repository) (Service, error) {
	if sentry == nil || l == nil {
		return nil, errors.New("manager: arguments for NewService cannot be nil")
	}
	return &service{
		sentry:   sentry,
		log: l,
		repo: repo,
	}, nil
}

func (s *service) CheckAll(ctx context.Context, locationAddress string) ([]*Socket, error) {
	defer s.log.Sync()
	s.log.Info("Service: CheckAll()")
	sock, err := s.repo.FindSocketsByLocation(ctx, locationAddress)
	if err != nil {
		s.log.Errorw("Service: CheckAll() - repo call", zap.String("error", err.Error()))
		return nil, err
	}
	return sock, nil
}

func (s *service) AddSocket(ctx context.Context, socket Socket, locationAddress string) (*Socket, error) {
	defer s.log.Sync()
	s.log.Info("Service: AddSocket()")
	socket.SNMPaddress=locationAddress
	// TODO: add repo call
	sock, err := s.repo.CreateSocket(ctx, socket)
	if err != nil {
		return nil, err
	}
	
	return sock, nil
}

func (s *service) RemoveSocket(ctx context.Context, socketID int) error {
	defer s.log.Sync()
	s.log.Info("Service: RemoveSocket()")
	if err := s.repo.DeleteSocket(ctx, socketID); err != nil {
		s.log.Errorw("Service: RemoveSocket() - repo call", zap.String("error", err.Error()))
		return  err
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
	_, err = s.sentry.ToggleSocket(ctx,onOrOff, socket.SNMPmib, socket.SNMPaddress, socket.SNMPcommunity, socket.SNMPport)
	if err != nil {
		s.log.Errorw("Service, ToggleSocket() - sentry call", zap.String("error", err.Error()))
		return err
	}
	return nil
}
