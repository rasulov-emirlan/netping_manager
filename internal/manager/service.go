package manager

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	AddLocation(ctx context.Context, l Location) ([]*Location, error)
	RemoveLocation(ctx context.Context, locationID int) ([]*Location, error)

	AddSocket(ctx context.Context, s Socket, locationID int) ([]*Location, error)
	RemoveSocket(ctx context.Context, socketID, locationID int) ([]*Location, error)

	CheckAll(ctx context.Context) ([]*Location, error)
	ToggleSocket(ctx context.Context, socketID, locationId int, onOrOff int) ([]*Location, error)
}

const (
	SocketOn  = 1
	SocketOff = 2
)

type Watcher interface {
	AddLocation(ctx context.Context, l Location) ([]*Location, error)
	RemoveLocation(ctx context.Context, locationID int) ([]*Location, error)

	AddSocket(ctx context.Context, s Socket, locationID int) ([]*Location, error)
	RemoveSocket(ctx context.Context, socketID, locationID int) ([]*Location, error)

	CheckAll(ctx context.Context) ([]*Location, error)
	ToggleSocket(ctx context.Context, socketID, locationId int, onOrOff int) ([]*Location, error)
}

type service struct {
	w   Watcher
	log *zap.SugaredLogger
}

func NewService(w Watcher, l *zap.SugaredLogger) (Service, error) {
	return &service{
		w:   w,
		log: l,
	}, nil
}

func (s *service) CheckAll(ctx context.Context) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: CheckAll()")
	v, err := s.w.CheckAll(ctx)
	if err != nil {
		s.log.Errorf("Service: CheckAll() - error: %v", err)
		return nil, err
	}
	return v, nil
}

func (s *service) AddLocation(ctx context.Context, l Location) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: AddLocation()")
	v, err := s.w.AddLocation(ctx, l)
	if err != nil {
		s.log.Errorf("Service: AddLocation() - error: %v", err)
		return nil, err
	}
	return v, nil
}

func (s *service) RemoveLocation(ctx context.Context, locationID int) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: RemoveLocation()")
	v, err := s.w.RemoveLocation(ctx, locationID)
	if err != nil {
		s.log.Errorf("Service: RemoveLocation() - error: %v", err)
		return nil, err
	}
	return v, nil
}

func (s *service) AddSocket(ctx context.Context, socket Socket, locationID int) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: AddSocket()")
	v, err := s.w.AddSocket(ctx, socket, locationID)
	if err != nil {
		s.log.Error("Service: AddSocket() - error: %v", err)
		return nil, err
	}
	return v, err
}

func (s *service) RemoveSocket(ctx context.Context, socketID, locationID int) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: RemoveSocket()")
	v, err := s.w.RemoveSocket(ctx, socketID, locationID)
	if err != nil {
		s.log.Error("Service: RemoveSocket() - error: %v", err)
		return nil, err
	}
	return v, err
}

func (s *service) ToggleSocket(ctx context.Context, socketID, locationId, onOrOff int) ([]*Location, error) {
	defer s.log.Sync()
	s.log.Info("Service: ToggleSocket()")
	v, err := s.w.ToggleSocket(ctx, socketID, locationId, onOrOff)
	if err != nil {
		s.log.Error("Service: ToggleSocket() - error: %v", err)
		return nil, err
	}
	return v, err
}
