package users

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, name, password string, isAdmin bool) (User, error)
	Read(ctx context.Context, userID int) (User, error)
	ReadByName(ctx context.Context, name string) (User, error)
	ReadAll(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, userID int) error
}

type Repository interface {
	Create(ctx context.Context, name, password string, isAdmin bool) (User, error)
	Read(ctx context.Context, userID int) (User, error)
	ReadByName(ctx context.Context, name string) (User, error)
	ReadAll(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, userID int) error
}

type service struct {
	repo Repository
	log  *zap.SugaredLogger
}

func NewService(repo Repository, log *zap.SugaredLogger) (Service, error) {
	return &service{
		repo: repo,
		log:  log,
	}, nil
}

func (s *service) Create(ctx context.Context, name, password string, isAdmin bool) (User, error) {
	defer s.log.Sync()
	s.log.Info("UserService: Create()")
	u, err := s.repo.Create(ctx, name, password, isAdmin)
	if err != nil {
		s.log.Errorw("UserService: Create() - repo call", zap.String("error", err.Error()))
		return User{}, err
	}
	return u, err
}

func (s *service) Read(ctx context.Context, userID int) (User, error) {
	defer s.log.Sync()
	s.log.Info("UserService: Read()")
	u, err := s.repo.Read(ctx, userID)
	if err != nil {
		s.log.Errorw("UserService: Read() - repo call", zap.String("error", err.Error()))
		return User{}, err
	}
	return u, err
}

func (s *service) ReadByName(ctx context.Context, name string) (User, error) {
	defer s.log.Sync()
	s.log.Info("UserService: ReadByName()")
	u, err := s.repo.ReadByName(ctx, name)
	if err != nil {
		s.log.Errorw("UserService: ReadByName() - repo call", zap.String("error", err.Error()))
		return User{}, err
	}
	return u, err
}

func (s *service) ReadAll(ctx context.Context) ([]User, error) {
	defer s.log.Sync()
	s.log.Info("UserService: ReadAll()")
	u, err := s.repo.ReadAll(ctx)
	if err != nil {
		s.log.Errorw("UserService: ReadAll() - repo call", zap.String("error", err.Error()))
		return nil, err
	}
	return u, err
}

func (s *service) Delete(ctx context.Context, userID int) error {
	defer s.log.Sync()
	s.log.Info("UserService: Delete()")
	if err := s.repo.Delete(ctx, userID); err != nil {
		s.log.Errorw("UserService: Delete() - repo call", zap.String("error", err.Error()))
		return err
	}
	return nil
}
