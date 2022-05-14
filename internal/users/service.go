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
	Update(ctx context.Context, userID int, changeset User) error
	Delete(ctx context.Context, userID int) error

	Login(ctx context.Context, name, password string) (User, error)
}

type Repository interface {
	Create(ctx context.Context, user User) (id int, err error)
	Read(ctx context.Context, userID int) (User, error)
	ReadByName(ctx context.Context, name string) (User, error)
	ReadAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, userID int, changeset User) error
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
	u, err := NewUser(name, password)
	if err != nil {
		return User{}, err
	}
	u.IsAdmin = isAdmin
	u.ID, err = s.repo.Create(ctx, u)
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

// Even though update accepts a struct User, it expects that password in it
// will be either empty or not hashed.
// If password is empty it will not update the password
// If password is given it will hash it first and then update
func (s *service) Update(ctx context.Context, userID int, changeset User) error {
	defer s.log.Sync()
	s.log.Info("UserService: Update()")
	if changeset.Password != "" {
		psw, err := HashePassword(changeset.Password)
		if err != nil {
			s.log.Errorw("UserService: Update() - hashing password", zap.String("error", err.Error()))
			return err
		}
		changeset.Password = string(psw)
	}
	if err := s.repo.Update(ctx, userID, changeset); err != nil {
		s.log.Errorw("UserService: Update() - repo call", zap.String("error", err.Error()))
		return err
	}
	return nil
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

func (s *service) Login(ctx context.Context, name, password string) (User, error) {
	defer s.log.Sync()
	u, err := s.repo.ReadByName(ctx, name)
	if err != nil {
		s.log.Errorw("UserService: Login() - repo call", zap.Error(err))
		return User{}, err
	}
	if err := u.ComparePasswords(password); err != nil {
		return User{}, err
	}
	return u, nil
}